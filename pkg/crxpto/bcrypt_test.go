package crxpto

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestNewBcrypt(t *testing.T) {
	assert.NotEmpty(t, NewBcrypt())
}

func TestBcrypt_HashPassword(t *testing.T) {
	b := &Bcrypt{}

	// cost too big (more than 31)
	b.cost = bcrypt.MaxCost + 1
	got, err := b.HashPassword("a")
	assert.Error(t, err)
	assert.Nil(t, got)

	// password too long (more than 72 characters)
	b.cost = bcrypt.DefaultCost
	got, err = b.HashPassword(strings.Repeat("a", 73))
	assert.Error(t, err)
	assert.Nil(t, got)

	// success (can't assert the value exactly due to its non-deterministic nature)
	got, err = b.HashPassword("abcde")
	assert.NoError(t, err)
	assert.NotEmpty(t, got)
}
