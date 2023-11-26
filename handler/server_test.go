package handler

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/leguminosa/profile-open-portal/module"
	"github.com/leguminosa/profile-open-portal/tools"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserModule := module.NewMockUserModuleInterface(ctrl)
	mockAuth := tools.NewMockAuthInterface(ctrl)

	assert.NotEmpty(t, NewServer(NewServerOptions{
		UserModule: mockUserModule,
		Auth:       mockAuth,
	}))
}
