package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"mini-paas/backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type LogWSHandler struct {
	logService services.K8sLogService
	writeWait  time.Duration // limit time send a message to Websocket conn, exceeds writeTime -> close conn
	pongWait   time.Duration //
	pingPeriod time.Duration // 1 ping/ 54s
	maxQueue   int           //
}

func NewLogWSHandler(ls services.K8sLogService) *LogWSHandler {
	return &LogWSHandler{
		logService: ls,
		writeWait:  10 * time.Second,
		pongWait:   60 * time.Second,
		pingPeriod: (60 * time.Second * 9) / 10,
		maxQueue:   1000,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *LogWSHandler) StreamDeploymentLogs(c *gin.Context) {
	deployID := c.Param("id")
	namespace := c.DefaultQuery("namespace", "default")
	followStr := c.DefaultQuery("follow", "true")
	follow := followStr == "true"
	tailStr := c.DefaultQuery("tailLines", "")
	var tail *int64

	if tailStr != "" {
		if v, err := strconv.ParseInt(tailStr, 10, 64); err == nil {
			tail = &v
		}
	}

	deploymentName := deployID

	// Upgrade to WS

	wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	defer wsConn.Close()

	// set deadlines and handlers
	wsConn.SetReadLimit(512)
	_ = wsConn.SetReadDeadline(time.Now().Add(h.pongWait))
	wsConn.SetPongHandler(func(string) error {
		_ = wsConn.SetReadDeadline(time.Now().Add(h.pongWait))
		return nil
	})

	// find pods
	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	pods, err := h.logService.FindPodsForDeployment(ctx, deploymentName, namespace)
	if err != nil {
		wsConn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("error: %v", err)))
		return
	}

	pod := pods[0]

	// open log stream
	stream, err := h.logService.StreamPodLogs(ctx, namespace, pod.Name, follow, tail)
	if err != nil {
		wsConn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("error opening logs: %v", err)))
		return
	}

	// line channel and writer goroutine
	lineCh := make(chan string, h.maxQueue)
	go services.StreamToLines(ctx, stream, lineCh)

	// writer: send lines to ws with write deadline
	writeErrCh := make(chan error, 1)
	go func() {
		ticker := time.NewTicker(h.pingPeriod)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				writeErrCh <- nil
				return

			case <-ticker.C:
				if err := wsConn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(h.writeWait)); err != nil {
					writeErrCh <- nil
					return
				}
			case line, ok := <-lineCh:
				if !ok {
					writeErrCh <- nil
					return
				}
				_ = wsConn.SetWriteDeadline(time.Now().Add(h.writeWait))
				if err := wsConn.WriteMessage(websocket.TextMessage, []byte(line)); err != nil {
					writeErrCh <- err
					return
				}
			}
		}
	}()

	// reader: read control messages to detect client disconnect
	go func() {
		for {
			_, _, err := wsConn.ReadMessage()
			if err != nil {
				cancel()
				return
			}
		}
	}()

	select {
	case err := <-writeErrCh:
		if err != nil {
			// log err
		}
	case <-ctx.Done():
		// cancelled
	}
}
