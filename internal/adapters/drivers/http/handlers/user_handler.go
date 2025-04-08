package handlers

import (
	"net/http"

	"github.com/DamiaoCanndido/na-mosca-server/internal/ports"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *ports.UserService
}

func NewUserHandler(service *ports.UserService) *UserHandler {
	return &UserHandler{service: service}
}

type RegisterUserRequest struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=Password"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.RegisterUser(req.Name, req.Email, req.Password)
	if err != nil {
		if err.Error() == "email já está em uso" {
			c.JSON(http.StatusConflict, gin.H{"error": "Este email já está cadastrado"})
			return
		}
		if err.Error() == "a senha deve ter no mínimo 6 caracteres" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "A senha deve ter no mínimo 6 caracteres"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar usuário"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Authenticate(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
} 