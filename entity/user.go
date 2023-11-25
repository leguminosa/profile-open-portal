package entity

import "github.com/leguminosa/profile-open-portal/tools"

type (
	User struct {
		ID             int    `json:"id"            db:"id"`
		Fullname       string `json:"fullname"      db:"fullname"`
		PhoneNumber    string `json:"phone_number"  db:"phone_number"`
		HashedPassword string `json:"-"             db:"password"`
		PlainPassword  string `json:"password"      db:"-"`
	}
)

func (u *User) Exist() bool {
	return u.ID != 0
}

func (u *User) HashPassword(hash tools.HashInterface) error {
	hashedPassword, err := hash.HashPassword(u.PlainPassword)
	if err != nil {
		return err
	}
	u.HashedPassword = string(hashedPassword)

	return nil
}

type (
	RegisterAPIRequest struct {
		*User
	}
	RegisterAPIResponse struct {
		UserID  int    `json:"user_id,omitempty"`
		Message string `json:"message,omitempty"`
	}
	RegisterModuleRequest struct {
		User *User
	}
	RegisterModuleResponse struct {
		Valid    bool
		Messages []string
		User     *User
	}
)

type (
	LoginAPIRequest struct {
		PhoneNumber string `json:"phone_number"`
		Password    string `json:"password"`
	}
	LoginAPIResponse struct {
		UserID  int    `json:"user_id,omitempty"`
		JWT     string `json:"jwt,omitempty"`
		Message string `json:"message,omitempty"`
	}
	LoginModuleRequest struct {
		PhoneNumber string
		Password    string
	}
	LoginModuleResponse struct {
		User *User
		JWT  string
	}
)
