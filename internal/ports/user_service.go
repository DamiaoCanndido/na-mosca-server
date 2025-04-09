package ports

import (
	"errors"
	"os"
	"time"

	"github.com/DamiaoCanndido/na-mosca-server/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterUser(name, email, password string) (*domain.User, error) {
	// Verificar se o email já existe
	existingUser, err := s.repo.FindByEmail(email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email já está em uso")
	}

	// Verificar tamanho mínimo da senha
	if len(password) < 6 {
		return nil, errors.New("a senha deve ter no mínimo 6 caracteres")
	}

	user := &domain.User{
		ID:        uuid.New(),
		Name:      name,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
	}

	err = s.repo.RegisterUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Authenticate(email, password string) (string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	err = s.repo.VerifyPassword(user, password)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID.String(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
} 