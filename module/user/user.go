package user

import (
	"context"

	"github.com/leguminosa/profile-open-portal/entity"
	"github.com/leguminosa/profile-open-portal/pkg/validator"
	"github.com/leguminosa/profile-open-portal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserModule struct {
	userRepository repository.UserRepositoryInterface
}

type NewUserModuleOptions struct {
	UserRepository repository.UserRepositoryInterface
}

func New(opts NewUserModuleOptions) *UserModule {
	return &UserModule{
		userRepository: opts.UserRepository,
	}
}

func (m *UserModule) Register(ctx context.Context, req entity.RegisterModuleRequest) (entity.RegisterModuleResponse, error) {
	var (
		resp = entity.RegisterModuleResponse{
			Valid:    true,
			Messages: []string{},
			User:     req.User,
		}
		err error
	)

	// validate request
	var (
		messages []string
		valid    bool
	)
	if messages, valid = validator.ValidatePhoneNumber(req.User.PhoneNumber); !valid {
		resp.Valid = false
		resp.Messages = append(resp.Messages, messages...)
	}
	if messages, valid = validator.ValidateFullName(req.User.Fullname); !valid {
		resp.Valid = false
		resp.Messages = append(resp.Messages, messages...)
	}
	if messages, valid = validator.ValidatePassword(req.User.PlainPassword); !valid {
		resp.Valid = false
		resp.Messages = append(resp.Messages, messages...)
	}

	if !resp.Valid {
		return resp, nil
	}

	// hashing plain password before inserting to database
	var hashedPassword []byte
	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(req.User.PlainPassword), bcrypt.DefaultCost)
	if err != nil {
		return resp, err
	}
	req.User.HashedPassword = string(hashedPassword)

	resp.User.ID, err = m.userRepository.InsertUser(ctx, req.User)
	if err != nil {
		return resp, err
	}

	resp.Valid = true
	return resp, nil
}
