package dtos

import (
	"regexp"
	"strings"
)

type RegisterUserRequest struct {
	Name            string  `json:"name" binding:"required" validate:"required,min=3,max=100"`
	AvatarUrl      	*string `json:"avatar_url,omitempty"`
	Email           string  `json:"email" binding:"required" validate:"required,email"`
	Password        string  `json:"password" binding:"required" validate:"required,min=6"`
	ConfirmPassword string  `json:"confirm_password" binding:"required" validate:"required,eqfield=Password"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required" validate:"required,email"`
	Password string `json:"password" binding:"required" validate:"required,min=6"`
}

func (r *RegisterUserRequest) Validate() map[string]string {
	errors := make(map[string]string)

	if strings.TrimSpace(r.Name) == "" {
		errors["name"] = "Name is required"
	} else if len(r.Name) < 3 {
		errors["name"] = "Name must have at least 3 characters"
	} else if len(r.Name) > 100 {
		errors["name"] = "Name cannot have more than 100 characters"
	}

	if r.AvatarUrl != nil && len(*r.AvatarUrl) > 255 {
		errors["avatar_url"] = "Avatar URL cannot have more than 255 characters"
	}
	urlRegex := regexp.MustCompile(`^(https):\/\/[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,}(\/\S*)?$`)
    if r.AvatarUrl != nil && !urlRegex.MatchString(*r.AvatarUrl) {
        errors["avatar_url"] = "Please enter a valid URL for the avatar"
    }

	if strings.TrimSpace(r.Email) == "" {
		errors["email"] = "Email is required"
	} else {
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(r.Email) {
			errors["email"] = "Please enter a valid email"
		}
	}

	if strings.TrimSpace(r.Password) == "" {
		errors["password"] = "Password is required"
	} else if len(r.Password) < 6 {
		errors["password"] = "Password must have at least 6 characters"
	} else {
		hasNumber := regexp.MustCompile(`[0-9]`).MatchString(r.Password)
		hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(r.Password)
		hasLower := regexp.MustCompile(`[a-z]`).MatchString(r.Password)

		if !hasNumber {
			errors["password"] = "Password must contain at least one number"
		} else if !hasUpper {
			errors["password"] = "Password must contain at least one uppercase letter"
		} else if !hasLower {
			errors["password"] = "Password must contain at least one lowercase letter"
		}
	}

	if strings.TrimSpace(r.ConfirmPassword) == "" {
		errors["confirm_password"] = "Password confirmation is required"
	} else if r.Password != r.ConfirmPassword {
		errors["confirm_password"] = "Passwords do not match"
	}

	return errors
}

func (l *LoginRequest) Validate() map[string]string {
	errors := make(map[string]string)

	// Email validation
	if strings.TrimSpace(l.Email) == "" {
		errors["email"] = "Email is required"
	} else {
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(l.Email) {
			errors["email"] = "Please enter a valid email"
		}
	}

	// Password validation
	if strings.TrimSpace(l.Password) == "" {
		errors["password"] = "Password is required"
	} else if len(l.Password) < 6 {
		errors["password"] = "Password must have at least 6 characters"
	}

	return errors
}