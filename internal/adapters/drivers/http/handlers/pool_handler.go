package handlers

import (
	"github.com/DamiaoCanndido/na-mosca-server/internal/adapters/drivers/http/dtos"
	"github.com/DamiaoCanndido/na-mosca-server/internal/ports"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PoolHandler struct {
	service *ports.PoolService
}

func NewPoolHandler(service *ports.PoolService) *PoolHandler {
	return &PoolHandler{service: service}
}

func (h *PoolHandler) CreatePool(c *gin.Context) {
	var req dtos.CreatePoolRequest

	c.BindJSON(&req)

	// Validação personalizada
	if errors := req.Validate(); len(errors) > 0 {
		c.JSON(400, gin.H{"error": "Erro de validação", "details": errors})
		return
	}

	userID, _ := c.Get("user_id")
	userUUID, _ := uuid.Parse(userID.(string))

	poolName, err := h.service.Create(req.Name, userUUID)

	if err != nil {
		c.JSON(400, gin.H{"error": "Erro ao criar pool"})
		return
	}

	c.JSON(201, gin.H{"pool_name": poolName})
}