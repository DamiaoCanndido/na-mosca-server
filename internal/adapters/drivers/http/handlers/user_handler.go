package handlers

import (
	"net/http"

	"github.com/DamiaoCanndido/na-mosca-server/internal/adapters/drivers/http/dtos"
	"github.com/DamiaoCanndido/na-mosca-server/internal/ports"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *ports.UserService
}

func NewUserHandler(service *ports.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req dtos.RegisterUserRequest

	c.BindJSON(&req)

	// Validação personalizada
	if errors := req.Validate(); len(errors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro de validação", "details": errors})
		return
	}

	avatarUrl := ""
    if req.AvatarUrl != nil {
        avatarUrl = *req.AvatarUrl
    }

	user, err := h.service.RegisterUser(req.Name, avatarUrl, req.Email, req.Password)
	if err != nil {
		if err.Error() == "email já está em uso" {
			c.JSON(http.StatusConflict, gin.H{"error": "Este email já está cadastrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar usuário"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dtos.LoginRequest

	c.BindJSON(&req)

	// Validação personalizada
	if errors := req.Validate(); len(errors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro de validação", "details": errors})
		return
	}

	token, err := h.service.Authenticate(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
} 

func (h *UserHandler) GetMe(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token not provided"})
		return
	}

	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	// Convert userID to uuid.UUID
	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	user, err := h.service.GetMe(token, userUUID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	c.JSON(http.StatusOK, user)
}