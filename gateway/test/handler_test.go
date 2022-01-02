package test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
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
		{credentials{"lisa@gmail.com", "qwerty", "Lisa"}, http.StatusCreated},
		{credentials{"", "qwerty", "Lisa"}, http.StatusBadRequest},
		{credentials{"lisa@gmail.com", "", "Lisa"}, http.StatusBadRequest},
		{credentials{"lisa@gmail.com", "qwerty", ""}, http.StatusBadRequest},
		{credentials{"lisa@gmail.com", "qwerty", ""}, http.StatusBadRequest},
		{credentials{"lisa@gmail.com", "", "Lisa"}, http.StatusBadRequest},
		{credentials{"", "qwerty", "Lisa"}, http.StatusBadRequest},
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
		{credentials{"lisa@gmail.com", "qwerty"}, http.StatusOK},
		{credentials{"lisa@gmail.com", ""}, http.StatusBadRequest},
		{credentials{"", "qwerty"}, http.StatusBadRequest},
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
