package userRepo

import (
	"bolao/internal/domain"
	"crypto/rand"
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/argon2"
	"gorm.io/gorm"
)

const (
	memory      = 64 * 1024
	iterations  = 3
	parallelism = 2
	saltLength  = 16
	keyLength   = 32
)

type UserRepository struct {
	db *gorm.DB
}

// Create implements domain.UserRepository.
func (r *UserRepository) Create(user *domain.User) error {
	panic("unimplemented")
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func generateSalt() ([]byte, error) {
	salt := make([]byte, saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func hashPassword(password string) (string, error) {
	salt, err := generateSalt()
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		iterations,
		memory,
		parallelism,
		keyLength,
	)

	// Combine salt and hash
	combined := make([]byte, saltLength+keyLength)
	copy(combined[:saltLength], salt)
	copy(combined[saltLength:], hash)

	// Encode to base64
	return base64.StdEncoding.EncodeToString(combined), nil
}

func verifyPassword(password, encodedHash string) bool {
	// Decode the combined salt and hash
	combined, err := base64.StdEncoding.DecodeString(encodedHash)
	if err != nil {
		return false
	}

	// Extract salt and hash
	salt := combined[:saltLength]
	hash := combined[saltLength:]

	// Compute hash of the provided password
	newHash := argon2.IDKey(
		[]byte(password),
		salt,
		iterations,
		memory,
		parallelism,
		keyLength,
	)

	// Compare hashes
	for i := range hash {
		if hash[i] != newHash[i] {
			return false
		}
	}
	return true
}

func (r *UserRepository) RegisterUser(user *domain.User) error {
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) VerifyPassword(user *domain.User, password string) error {
	if !verifyPassword(password, user.Password) {
		return errors.New("invalid password")
	}
	return nil
}
