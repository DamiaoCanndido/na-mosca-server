package dtos

import (
	"regexp"
	"strings"
)

type RegisterUserRequest struct {
	Name            string `json:"name" binding:"required" validate:"required,min=3,max=100"`
	Email           string `json:"email" binding:"required" validate:"required,email"`
	Password        string `json:"password" binding:"required" validate:"required,min=6"`
	ConfirmPassword string `json:"confirmPassword" binding:"required" validate:"required,eqfield=Password"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required" validate:"required,email"`
	Password string `json:"password" binding:"required" validate:"required"`
}

func (r *RegisterUserRequest) Validate() map[string]string {
	errors := make(map[string]string)

	// Validação do nome
	if strings.TrimSpace(r.Name) == "" {
		errors["name"] = "O nome é obrigatório"
	} else if len(r.Name) < 3 {
		errors["name"] = "O nome deve ter pelo menos 3 caracteres"
	} else if len(r.Name) > 100 {
		errors["name"] = "O nome não pode ter mais de 100 caracteres"
	}

	// Validação do email
	if strings.TrimSpace(r.Email) == "" {
		errors["email"] = "O email é obrigatório"
	} else {
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(r.Email) {
			errors["email"] = "Por favor, insira um email válido"
		}
	}

	// Validação da senha
	if strings.TrimSpace(r.Password) == "" {
		errors["password"] = "A senha é obrigatória"
	} else if len(r.Password) < 6 {
		errors["password"] = "A senha deve ter pelo menos 6 caracteres"
	} else {
		// Verifica se a senha contém pelo menos um número
		hasNumber := regexp.MustCompile(`[0-9]`).MatchString(r.Password)
		// Verifica se a senha contém pelo menos uma letra maiúscula
		hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(r.Password)
		// Verifica se a senha contém pelo menos uma letra minúscula
		hasLower := regexp.MustCompile(`[a-z]`).MatchString(r.Password)

		if !hasNumber {
			errors["password"] = "A senha deve conter pelo menos um número"
		} else if !hasUpper {
			errors["password"] = "A senha deve conter pelo menos uma letra maiúscula"
		} else if !hasLower {
			errors["password"] = "A senha deve conter pelo menos uma letra minúscula"
		}
	}

	// Validação da confirmação de senha
	if strings.TrimSpace(r.ConfirmPassword) == "" {
		errors["confirmPassword"] = "A confirmação de senha é obrigatória"
	} else if r.Password != r.ConfirmPassword {
		errors["confirmPassword"] = "As senhas não coincidem"
	}

	return errors
}

func (l *LoginRequest) Validate() map[string]string {
	errors := make(map[string]string)

	// Validação do email
	if strings.TrimSpace(l.Email) == "" {
		errors["email"] = "O email é obrigatório"
	} else {
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(l.Email) {
			errors["email"] = "Por favor, insira um email válido"
		}
	}

	// Validação da senha
	if strings.TrimSpace(l.Password) == "" {
		errors["password"] = "A senha é obrigatória"
	}

	return errors
} 