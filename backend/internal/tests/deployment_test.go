package tests

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestDeploymentIntegration(t *testing.T) {
	// create app prerequisite
	appPayload := `{"name":"dep-app", "git_url":"https://example.com/repo.git"}`
	appResp, err := http.Post(
		testServer.URL+"/api/apps",
		"application/json",
		strings.NewReader(appPayload),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer appResp.Body.Close()
	if appResp.StatusCode != http.StatusCreated {
		t.Fatalf("create app expected 201 got %d", appResp.StatusCode)
	}
	var appCreated map[string]interface{}
	json.NewDecoder(appResp.Body).Decode(&appCreated)
	appID := appCreated["id"].(string)

	//  deploy app to Minikube
	payload := `{"app_id":"` + appID + `", "version":"v1", "image_url":"nginx:stable"}`
	resp, err := http.Post(
		testServer.URL+"/api/deployments/deploy",
		"application/json",
		strings.NewReader(payload),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("deploy expected 202 got %d", resp.StatusCode)
	}

	var created map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&created)
	deployID := created["id"].(string)

	// poll status until RUNNING or FAILED or timeout
	deadline := time.Now().Add(90 * time.Second)
	var lastStatus string
	for time.Now().Before(deadline) {
		resp2, err := http.Get(testServer.URL + "/api/deployments/" + deployID + "/status")
		if err != nil {
			t.Fatalf("status request failed: %v", err)
		}
		if resp2.StatusCode != http.StatusOK {
			resp2.Body.Close()
			time.Sleep(3 * time.Second)
			continue
		}
		var statusResp map[string]interface{}
		json.NewDecoder(resp2.Body).Decode(&statusResp)
		resp2.Body.Close()
		lastStatus, _ = statusResp["status"].(string)
		if lastStatus == "RUNNING" || lastStatus == "FAILED" {
			break
		}
		time.Sleep(5 * time.Second)
	}
	if lastStatus == "" {
		t.Fatalf("status polling timed out without result")
	}

	// list
	resp3, _ := http.Get(testServer.URL + "/api/deployments")
	if resp3.StatusCode != http.StatusOK {
		t.Fatalf("list deployments expected 200 got %d", resp3.StatusCode)
	}
	resp3.Body.Close()
}
