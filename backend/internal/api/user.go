package api

import (
	"mini-paas/backend/internal/models"
	"mini-paas/backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(s services.UserService) *UserHandler {
	return &UserHandler{userService: s}
}

// POST /api/users
func (h *UserHandler) CreateUserHandler(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &models.User{
		Name:  req.Name,
		Email: req.Email,
	}

	newUsr, err := h.userService.CreateUser(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, UserResponse{
		ID:    newUsr.ID.String(),
		Name:  newUsr.Name,
		Email: newUsr.Email,
	})
}

// // GET /api/users
// func (h *UserHandler)  {

// }

// GET /api/users/user/:id
func (h *UserHandler) GetUserByIDHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return
	}

	user, err := h.userService.GetUserByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
	})
}

// GET /api/users/user/:email
func (h *UserHandler) GetUserByEmailHandler(c *gin.Context) {
	email := c.Query("email")

	user, err := h.userService.GetUserByEmail(c, email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
	})
}

// DELETE /api/users/user/:id
