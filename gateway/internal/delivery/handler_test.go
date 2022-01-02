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
	handler = NewHandler(mocks.UsersClientMock{}, mocks.LoggerMock{})
	router  = handler.InitRoutes()
)

func TestHandler_signUp(t *testing.T) {
	t.Parallel()

	type credentials struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	testCases := []struct {
		body               credentials
		expectedStatusCode int
	}{
		{credentials{email, password, username}, http.StatusCreated},
		{credentials{"", password, username}, http.StatusBadRequest},
		{credentials{email, "", username}, http.StatusBadRequest},
		{credentials{email, password, ""}, http.StatusBadRequest},
		{credentials{email, password, ""}, http.StatusBadRequest},
		{credentials{email, "", username}, http.StatusBadRequest},
		{credentials{"", password, username}, http.StatusBadRequest},
		{credentials{"", "", ""}, http.StatusBadRequest},
	}
	for _, tc := range testCases {
		tc := tc
		body, err := json.Marshal(tc.body)
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/users/sign-up", bytes.NewReader(body))
		router.ServeHTTP(w, req)
		assert.Equal(t, tc.expectedStatusCode, w.Result().StatusCode)
	}
}

func TestHandler_signIn(t *testing.T) {
	t.Parallel()

	type credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	testCases := []struct {
		body               credentials
		expectedStatusCode int
	}{
		{credentials{email, password}, http.StatusOK},
		{credentials{email, ""}, http.StatusBadRequest},
		{credentials{"", password}, http.StatusBadRequest},
		{credentials{"", ""}, http.StatusBadRequest},
	}
	for _, tc := range testCases {
		tc := tc
		body, err := json.Marshal(tc.body)
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/users/sign-in", bytes.NewReader(body))
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

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/auth?token="+tc.token, nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, tc.expectedStatusCode, w.Result().StatusCode)
	}
}
