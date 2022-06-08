package delivery

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	testUserId      = "1"
	testMessage     = "Hello"
	testBearerToken = "Bearer token"
	testOffset      = "offset"
	testLimit       = 10
)

func TestHandler_createPost(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		message            string
		bearerToken        string
		expectedStatusCode int
	}{
		{
			message:            testMessage,
			bearerToken:        testBearerToken,
			expectedStatusCode: http.StatusCreated,
		},
		{
			message:            testMessage,
			bearerToken:        "",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			message:            "",
			bearerToken:        testBearerToken,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			message:            "",
			bearerToken:        "",
			expectedStatusCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		tc := tc

		body, err := json.Marshal(struct {
			Message string `json:"message" binding:"required"`
		}{
			Message: tc.message,
		})
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/api/v1/posts", bytes.NewReader(body))
		assert.NoError(t, err)
		req.Header.Set("Authorization", tc.bearerToken)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		router.ServeHTTP(w, req)
		assert.Equal(t, tc.expectedStatusCode, w.Result().StatusCode)
	}
}

func TestHandler_getPostsByUserId(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		userId             string
		offset             string
		limit              int64
		expectedStatusCode int
	}{
		{
			userId:             testUserId,
			offset:             "",
			limit:              0,
			expectedStatusCode: http.StatusOK,
		},
		{
			userId:             testUserId,
			offset:             testOffset,
			limit:              testLimit,
			expectedStatusCode: http.StatusOK,
		},
		{
			userId:             "",
			offset:             "",
			limit:              0,
			expectedStatusCode: http.StatusNotFound,
		},
		{
			userId:             "",
			offset:             testOffset,
			limit:              testLimit,
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		tc := tc

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/api/v1/posts/"+tc.userId, nil)
		assert.NoError(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, tc.expectedStatusCode, w.Result().StatusCode)
	}
}
