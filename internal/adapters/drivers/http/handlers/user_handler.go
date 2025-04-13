package handlers

import (
	"net/http"

	"github.com/DamiaoCanndido/na-mosca-server/internal/adapters/drivers/http/dtos"
	"github.com/DamiaoCanndido/na-mosca-server/internal/ports"

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

	user, err := h.service.RegisterUser(req.Name, req.Email, req.Password)
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