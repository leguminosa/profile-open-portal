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
	got, err := b.HashPassword("Abcd9!")
	assert.Error(t, err)
	assert.Nil(t, got)

	// password too long (more than 72 characters)
	b.cost = bcrypt.DefaultCost
	got, err = b.HashPassword(strings.Repeat("a", 73))
	assert.Error(t, err)
	assert.Nil(t, got)

	// success (can't assert the value exactly due to its non-deterministic nature)
	got, err = b.HashPassword("Abcd9!")
	assert.NoError(t, err)
	assert.NotEmpty(t, got)
}

func TestBcrypt_ComparePassword(t *testing.T) {
	b := &Bcrypt{}

	// hashed password is too short
	err := b.ComparePassword([]byte("a"), "")
	assert.Error(t, err)

	// success
	err = b.ComparePassword([]byte("$2a$10$q5ZnlsdYtSjIdkgesnGFquAeKYZ55YU5f5on/s4KthnjD2pDBldYy"), "Abcd9!")
	assert.NoError(t, err)
}
