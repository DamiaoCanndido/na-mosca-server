package dtos

import "github.com/DamiaoCanndido/na-mosca-server/internal/domain"

type CreatePoolRequest struct {
	Name string `json:"name" binding:"required" validate:"required,min=3,max=100"`
}

type PoolResponse struct {
	ID           string              `json:"id"`
	Name         string              `json:"name"`
	Participants []domain.User       `json:"participants"`
	Games        []domain.Game       `json:"games"`
	OwnerID      string              `json:"owner_id"`
}

func (r *CreatePoolRequest) Validate() map[string]string {
	errors := make(map[string]string)

	if r.Name == "" {
		errors["name"] = "Name is required"
	} else if len(r.Name) < 3 || len(r.Name) > 100 {
		errors["name"] = "Name must be between 3 and 100 characters"
	}

	return errors
}