package dtos

type CreatePoolRequest struct {
	Name    string    `json:"name" binding:"required" validate:"required,min=3,max=100"`
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