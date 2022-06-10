package delivery

import (
	"bytes"
	"encoding/json"
	"github.com/KirillMironov/rapu/gateway/test/mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	testEmail    = "lisa@gmail.com"
	testPassword = "qwerty"
	testUsername = "Lisa"
)

var (
	handler = NewHandler(mock.UsersClient{}, mock.PostsClient{}, mock.Logger{})
	router  = handler.InitRoutes()
)

func TestHandler_signUp(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		username           string
		email              string
		password           string
		expectedStatusCode int
	}{
		{
			username:           testUsername,
			email:              testEmail,
			password:           testPassword,
			expectedStatusCode: http.StatusCreated,
		},
		{
			username:           "",
			email:              testEmail,
			password:           testPassword,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			username:           testUsername,
			email:              "",
			password:           testPassword,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			username:           testUsername,
			email:              testEmail,
			password:           "",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			username:           "",
			email:              "",
			password:           testPassword,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			username:           "",
			email:              testEmail,
			password:           "",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			username:           testUsername,
			email:              "",
			password:           "",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			username:           "",
			email:              "",
			password:           "",
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		tc := tc

		body, err := json.Marshal(struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}{
			Username: tc.username,
			Email:    tc.email,
			Password: tc.password,
		})
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/api/v1/users/sign-up", bytes.NewReader(body))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		router.ServeHTTP(w, req)
		assert.Equal(t, tc.expectedStatusCode, w.Result().StatusCode)
	}
}

func TestHandler_signIn(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		email              string
		password           string
		expectedStatusCode int
	}{
		{
			email:              testEmail,
			password:           testPassword,
			expectedStatusCode: http.StatusOK,
		},
		{
			email:              "",
			password:           testPassword,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			email:              testEmail,
			password:           "",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			email:              "",
			password:           "",
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		tc := tc

		body, err := json.Marshal(struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}{
			Email:    tc.email,
			Password: tc.password,
		})
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/api/v1/users/sign-in", bytes.NewReader(body))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		router.ServeHTTP(w, req)
		assert.Equal(t, tc.expectedStatusCode, w.Result().StatusCode)
	}
}
