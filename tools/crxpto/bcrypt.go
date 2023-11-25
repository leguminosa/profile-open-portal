package crxpto

import (
	"golang.org/x/crypto/bcrypt"
)

type Bcrypt struct {
	cost int
}

func NewBcrypt() *Bcrypt {
	return &Bcrypt{
		cost: bcrypt.DefaultCost,
	}
}

func (b *Bcrypt) HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), b.cost)
}

func (b *Bcrypt) ComparePassword(hashedPassword []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
}
