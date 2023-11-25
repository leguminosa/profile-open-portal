package crxpto

import (
	"golang.org/x/crypto/bcrypt"
)

// Bcrypt wraps golang bcrypt primitive implementation.
type Bcrypt struct {
	cost int
}

// NewBcrypt returns a new Bcrypt instance.
func NewBcrypt() *Bcrypt {
	return &Bcrypt{
		cost: bcrypt.DefaultCost,
	}
}

// HashPassword hashes a password using bcrypt.
func (b *Bcrypt) HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), b.cost)
}

// ComparePassword compares a hashed password with a plain password.
func (b *Bcrypt) ComparePassword(hashedPassword []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
}
