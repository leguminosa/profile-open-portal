package entity

type (
	User struct {
		ID             int    `json:"id" db:"id"`
		Fullname       string `json:"fullname" db:"fullname"`
		PhoneNumber    string `json:"phone_number" db:"phone_number"`
		HashedPassword string `json:"-" db:"password"`
		PlainPassword  string `json:"password" db:"-"`
	}

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
