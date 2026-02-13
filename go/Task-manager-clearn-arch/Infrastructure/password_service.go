package infrastructure

import (
	"golang.org/x/crypto/bcrypt"
)

type PasswordService interface {
	Hash(password string) (string, error)
	Compare(password, hash string) bool
}

type passwordService struct{}

func NewPasswordService() PasswordService {
	return &passwordService{}
}

func (ps *passwordService) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (ps *passwordService) Compare(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
