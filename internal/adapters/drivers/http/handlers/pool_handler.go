package handlers

import (
	"net/http"

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
		c.JSON(400, gin.H{"error": "Validation error", "details": errors})
		return
	}

	userID, _ := c.Get("user_id")
	userUUID, _ := uuid.Parse(userID.(string))

	poolName, err := h.service.Create(req.Name, userUUID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error creating pool"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"name": poolName})
}

func (h *PoolHandler) GetPoolByID(c *gin.Context) {
	poolIDStr := c.Param("pool_id")
	poolID, err := uuid.Parse(poolIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pool ID"})
		return
	}

	pool, err := h.service.GetByID(poolID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pool not found"})
		return
	}

	c.JSON(http.StatusOK, dtos.PoolResponse{
		ID:           pool.ID.String(),
		Name:         pool.Name,
		Participants: pool.Participants,
		Games:        pool.Games,
		OwnerID:      pool.OwnerID.String(),
	})
}