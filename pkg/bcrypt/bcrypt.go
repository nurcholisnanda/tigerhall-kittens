package bcrypt

import "golang.org/x/crypto/bcrypt"

type Bcrypt struct{}

//go:generate mockgen -source=bcrypt.go -destination=mock/bcrypt.go -package=mock
type BcryptInterface interface {
	HashPassword(password string) ([]byte, error)
	ComparePassword(hashedPassword, password string) error
}

func NewBcrypt() BcryptInterface {
	return &Bcrypt{}
}

func (h *Bcrypt) HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (h *Bcrypt) ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
