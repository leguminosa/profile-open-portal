package user

import (
	"context"
	"errors"

	"github.com/leguminosa/profile-open-portal/entity"
	"github.com/leguminosa/profile-open-portal/repository"
	"github.com/leguminosa/profile-open-portal/tools"
	"github.com/leguminosa/profile-open-portal/tools/validator"
)

type UserModule struct {
	userRepository repository.UserRepositoryInterface
	hash           tools.HashInterface
	jwt            tools.JWTInterface
}

type NewUserModuleOptions struct {
	UserRepository repository.UserRepositoryInterface
	Hash           tools.HashInterface
	JWT            tools.JWTInterface
}

// New creates new user module.
func New(opts NewUserModuleOptions) *UserModule {
	return &UserModule{
		userRepository: opts.UserRepository,
		hash:           opts.Hash,
		jwt:            opts.JWT,
	}
}

// Register creates new user after validating the request.
func (m *UserModule) Register(ctx context.Context, user *entity.User) (entity.RegisterModuleResponse, error) {
	var (
		resp = entity.RegisterModuleResponse{
			User:     user,
			Valid:    true,
			Messages: []string{},
		}
		err error
	)

	// validate request
	var (
		messages []string
		valid    bool
	)
	if messages, valid = validator.ValidatePhoneNumber(user.PhoneNumber); !valid {
		resp.Valid = false
		resp.Messages = append(resp.Messages, messages...)
	}
	if messages, valid = validator.ValidateFullName(user.Fullname); !valid {
		resp.Valid = false
		resp.Messages = append(resp.Messages, messages...)
	}
	if messages, valid = validator.ValidatePassword(user.PlainPassword); !valid {
		resp.Valid = false
		resp.Messages = append(resp.Messages, messages...)
	}

	if !resp.Valid {
		return resp, nil
	}

	// hashing plain password before inserting to database
	err = user.HashPassword(m.hash)
	if err != nil {
		return resp, err
	}

	resp.User.ID, err = m.userRepository.InsertUser(ctx, user)
	if err != nil {
		return resp, err
	}

	resp.Valid = true
	return resp, nil
}

var (
	// ErrLoginFailed obscures the error message to prevent brute force attack
	ErrLoginFailed = errors.New("phone number or password is not correct")
)

// Login generate jwt and increment success login count on successful attempt.
func (m *UserModule) Login(ctx context.Context, user *entity.User) (entity.LoginModuleResponse, error) {
	var (
		resp = entity.LoginModuleResponse{
			User: user,
		}
		err error
	)

	// get user from database
	resp.User, err = m.userRepository.GetUserByPhoneNumber(ctx, user.PhoneNumber)
	if err != nil {
		return resp, ErrLoginFailed
	}

	// check whether user with requested phone number exist in database
	if !resp.User.Exist() {
		return resp, ErrLoginFailed
	}

	// compare hashed password stored in database with user input
	err = m.hash.ComparePassword([]byte(resp.User.HashedPassword), user.PlainPassword)
	if err != nil {
		return resp, ErrLoginFailed
	}

	resp.JWT, err = m.jwt.Generate(resp.User)
	if err != nil {
		return resp, ErrLoginFailed
	}

	err = m.userRepository.IncrementLoginCount(ctx, resp.User.ID)
	if err != nil {
		return resp, ErrLoginFailed
	}

	return resp, nil
}

// GetProfile returns user profile.
func (m *UserModule) GetProfile(ctx context.Context, userID int) (*entity.User, error) {
	return m.userRepository.GetUserByID(ctx, userID)
}

// UpdateProfile only updates fullname and/or phone number if user input is not empty.
func (m *UserModule) UpdateProfile(ctx context.Context, user *entity.User) (entity.UpdateProfileModuleResponse, error) {
	var resp entity.UpdateProfileModuleResponse

	// get user to db first to check whether user currentValue
	currentValue, err := m.userRepository.GetUserByID(ctx, user.ID)
	if err != nil {
		return resp, err
	}

	if !currentValue.Exist() {
		return resp, errors.New("user not found")
	}

	// only update if user input is not empty
	if user.Fullname != "" {
		currentValue.Fullname = user.Fullname
	}
	if user.PhoneNumber != "" {
		resp.Conflict = m.isPhoneNumberExist(ctx, user.PhoneNumber)
		currentValue.PhoneNumber = user.PhoneNumber
	}

	if resp.Conflict {
		// don't update if phone number already exist
		resp.Message = "phone number already exist"
		return resp, nil
	}

	return resp, m.userRepository.UpdateUser(ctx, currentValue)
}

func (m *UserModule) isPhoneNumberExist(ctx context.Context, phoneNumber string) bool {
	user, err := m.userRepository.GetUserByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return false
	}
	return user.Exist()
}
