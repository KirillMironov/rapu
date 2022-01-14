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
	testEmail    = "lisa@gmail.com"
	testPassword = "qwerty"
	testUsername = "Lisa"
)

var (
	handler = NewHandler(mocks.UsersClientMock{}, mocks.PostsClientMock{}, mocks.LoggerMock{})
	router  = handler.InitRoutes()
)

func TestHandler_signUp(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		form               signUpForm
		expectedStatusCode int
	}{
		{signUpForm{
			Username: testUsername,
			Email:    testEmail,
			Password: testPassword,
		}, http.StatusCreated},
		{signUpForm{
			Username: "",
			Email:    testEmail,
			Password: testPassword,
		}, http.StatusBadRequest},
		{signUpForm{
			Username: testUsername,
			Email:    "",
			Password: testPassword,
		}, http.StatusBadRequest},
		{signUpForm{
			Username: testUsername,
			Email:    testEmail,
			Password: "",
		}, http.StatusBadRequest},
		{signUpForm{
			Username: "",
			Email:    "",
			Password: testPassword,
		}, http.StatusBadRequest},
		{signUpForm{
			Username: "",
			Email:    testEmail,
			Password: "",
		}, http.StatusBadRequest},
		{signUpForm{
			Username: testUsername,
			Email:    "",
			Password: "",
		}, http.StatusBadRequest},
		{signUpForm{
			Username: "",
			Email:    "",
			Password: "",
		}, http.StatusBadRequest},
	}

	for _, tc := range testCases {
		tc := tc
		body, err := json.Marshal(tc.form)
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
		form               signInForm
		expectedStatusCode int
	}{
		{signInForm{Email: testEmail, Password: testPassword}, http.StatusOK},
		{signInForm{Email: "", Password: testPassword}, http.StatusBadRequest},
		{signInForm{Email: testEmail, Password: ""}, http.StatusBadRequest},
		{signInForm{Email: "", Password: ""}, http.StatusBadRequest},
	}

	for _, tc := range testCases {
		tc := tc
		body, err := json.Marshal(tc.form)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/api/v1/users/sign-in", bytes.NewReader(body))
		assert.NoError(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, tc.expectedStatusCode, w.Result().StatusCode)
	}
}
