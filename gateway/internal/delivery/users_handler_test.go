package delivery

import (
	"bytes"
	"encoding/json"
	"github.com/KirillMironov/rapu/gateway/test/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	email    = "lisa@gmail.com"
	password = "qwerty"
	username = "Lisa"
)

var (
	handler = NewHandler(mocks.UsersClientMock{}, mocks.PostsClientMock{}, mocks.LoggerMock{})
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
		{username, email, password, http.StatusCreated},
		{"", email, password, http.StatusBadRequest},
		{username, "", password, http.StatusBadRequest},
		{username, email, "", http.StatusBadRequest},
		{"", "", password, http.StatusBadRequest},
		{"", email, "", http.StatusBadRequest},
		{username, "", "", http.StatusBadRequest},
		{"", "", "", http.StatusBadRequest},
	}

	for _, tc := range testCases {
		tc := tc

		body, err := json.Marshal(struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}{tc.username, tc.email, tc.password})
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/api/v1/users/sign-up", bytes.NewReader(body))
		assert.NoError(t, err)

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
		{email, password, http.StatusOK},
		{"", password, http.StatusBadRequest},
		{email, "", http.StatusBadRequest},
		{"", "", http.StatusBadRequest},
	}

	for _, tc := range testCases {
		tc := tc

		body, err := json.Marshal(struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}{tc.email, tc.password})
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/api/v1/users/sign-in", bytes.NewReader(body))
		assert.NoError(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, tc.expectedStatusCode, w.Result().StatusCode)
	}
}

func TestHandler_auth(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		token              string
		expectedStatusCode int
	}{
		{"token", http.StatusOK},
		{"", http.StatusBadRequest},
	}

	for _, tc := range testCases {
		tc := tc

		body, err := json.Marshal(map[string]string{accessToken: tc.token})
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/api/v1/users/auth", bytes.NewReader(body))
		assert.NoError(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, tc.expectedStatusCode, w.Result().StatusCode)
	}
}
