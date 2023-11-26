package helper

import (
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestUserIDFromContext(t *testing.T) {
	tests := []struct {
		name string
		c    echo.Context
		want int
	}{
		{
			name: "success",
			c:    newMockEchoContext(nil),
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UserIDFromContext(tt.c)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSetUserIDToContext(t *testing.T) {
	c := newMockEchoContext(nil)

	userID := 1
	SetUserIDToContext(c, userID)

	assert.Equal(t, 1, UserIDFromContext(c))
}
