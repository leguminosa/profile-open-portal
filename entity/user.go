package entity

import (
	"time"

	"github.com/leguminosa/profile-open-portal/tools"
)

type (
	// User represents both users table and return value exposed as api object.
	User struct {
		ID             int       `json:"id"             db:"id"`
		Fullname       string    `json:"fullname"       db:"fullname"`
		PhoneNumber    string    `json:"phone_number"   db:"phone_number"`
		HashedPassword string    `json:"-"              db:"password"`
		LoginCount     int       `json:"-"              db:"login_count"`
		CreatedAt      time.Time `json:"-"              db:"created_at"`
		UpdatedAt      time.Time `json:"-"              db:"updated_at"`

		PlainPassword string `json:"password,omitempty" db:"-"`
	}
	RegisterModuleResponse struct {
		User     *User
		Valid    bool
		Messages []string
	}
	LoginModuleResponse struct {
		User *User
		JWT  string
	}
)

// Exist returns true if user has been saved to database.
func (u *User) Exist() bool {
	return u.ID != 0
}

// HashPassword fills HashedPassword field using PlainPassword field.
func (u *User) HashPassword(hash tools.HashInterface) error {
	hashedPassword, err := hash.HashPassword(u.PlainPassword)
	if err != nil {
		return err
	}
	u.HashedPassword = string(hashedPassword)

	return nil
}
