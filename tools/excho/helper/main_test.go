package helper

import (
	"io"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
)

var (
	mockCtx echo.Context
)

func TestMain(m *testing.M) {
	mockCtx = newMockEchoContext(nil)

	os.Exit(m.Run())
}

type (
	mockEchoContext struct {
		echo.Context

		mockBind func(i interface{}) error
	}
)

func newMockEchoContext(m *mockEchoContext) *mockEchoContext {
	if m == nil {
		m = &mockEchoContext{}
	}

	m.Context = echo.New().NewContext(
		httptest.NewRequest(echo.GET, "/", nil),
		httptest.NewRecorder(),
	)

	return m
}

func (m *mockEchoContext) getResponseBody() []byte {
	httpResp := m.Response().Writer.(*httptest.ResponseRecorder)
	body, _ := io.ReadAll(httpResp.Body)
	return body
}
