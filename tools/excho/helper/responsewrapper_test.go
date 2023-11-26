package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOK(t *testing.T) {
	c := newMockEchoContext(nil)
	tests := []struct {
		name    string
		i       interface{}
		want    string
		wantErr bool
	}{
		{
			name: "success",
			i: map[string]interface{}{
				"message": "success",
			},
			want:    "{\"message\":\"success\"}\n",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := OK(c, tt.i)
			if !assert.Equal(t, tt.wantErr, err != nil) {
				return
			}

			got := c.getResponseBody()
			assert.Equal(t, tt.want, string(got))
		})
	}
}

func TestBadRequest(t *testing.T) {
	c := newMockEchoContext(nil)
	tests := []struct {
		name    string
		message string
		want    string
		wantErr bool
	}{
		{
			name:    "success",
			message: "bad request",
			want:    "{\"message\":\"bad request\"}\n",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := BadRequest(c, tt.message)
			if !assert.Equal(t, tt.wantErr, err != nil) {
				return
			}

			got := c.getResponseBody()
			assert.Equal(t, tt.want, string(got))
		})
	}
}

func TestForbidden(t *testing.T) {
	c := newMockEchoContext(nil)
	tests := []struct {
		name    string
		message string
		want    string
		wantErr bool
	}{
		{
			name:    "success",
			message: "forbidden",
			want:    "{\"message\":\"forbidden\"}\n",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Forbidden(c, tt.message)
			if !assert.Equal(t, tt.wantErr, err != nil) {
				return
			}

			got := c.getResponseBody()
			assert.Equal(t, tt.want, string(got))
		})
	}
}

func TestConflict(t *testing.T) {
	c := newMockEchoContext(nil)
	tests := []struct {
		name    string
		message string
		want    string
		wantErr bool
	}{
		{
			name:    "success",
			message: "conflicted",
			want:    "{\"message\":\"conflicted\"}\n",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Conflict(c, tt.message)
			if !assert.Equal(t, tt.wantErr, err != nil) {
				return
			}

			got := c.getResponseBody()
			assert.Equal(t, tt.want, string(got))
		})
	}
}

func TestInternalServerError(t *testing.T) {
	c := newMockEchoContext(nil)
	tests := []struct {
		name    string
		message string
		want    string
		wantErr bool
	}{
		{
			name:    "success",
			message: "server error",
			want:    "{\"message\":\"server error\"}\n",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := InternalServerError(c, tt.message)
			if !assert.Equal(t, tt.wantErr, err != nil) {
				return
			}

			got := c.getResponseBody()
			assert.Equal(t, tt.want, string(got))
		})
	}
}

func TestJSON(t *testing.T) {
	c := newMockEchoContext(nil)
	tests := []struct {
		name    string
		code    int
		i       interface{}
		want    string
		wantErr bool
	}{
		{
			name: "success",
			code: 200,
			i: map[string]interface{}{
				"message": "success",
			},
			want:    "{\"message\":\"success\"}\n",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := JSON(c, tt.code, tt.i)
			if !assert.Equal(t, tt.wantErr, err != nil) {
				return
			}

			got := c.getResponseBody()
			assert.Equal(t, tt.want, string(got))
		})
	}
}
