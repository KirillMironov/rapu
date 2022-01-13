package delivery

import (
	"bytes"
	"encoding/json"
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
		form               createPostForm
		bearerToken        string
		expectedStatusCode int
	}{
		{createPostForm{UserId: testUserId, Message: testMessage}, testBearerToken, http.StatusCreated},
		{createPostForm{UserId: testUserId, Message: testMessage}, "", http.StatusUnauthorized},
		{createPostForm{UserId: "", Message: testMessage}, testBearerToken, http.StatusBadRequest},
		{createPostForm{UserId: "", Message: testMessage}, "", http.StatusUnauthorized},
		{createPostForm{UserId: testUserId, Message: ""}, testBearerToken, http.StatusBadRequest},
		{createPostForm{UserId: testUserId, Message: ""}, "", http.StatusUnauthorized},
		{createPostForm{UserId: "", Message: ""}, testBearerToken, http.StatusBadRequest},
		{createPostForm{UserId: "", Message: ""}, "", http.StatusUnauthorized},
	}

	for _, tc := range testCases {
		tc := tc
		body, err := json.Marshal(tc.form)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/api/v1/posts", bytes.NewReader(body))
		assert.NoError(t, err)
		req.Header.Set("Authorization", tc.bearerToken)

		router.ServeHTTP(w, req)
		assert.Equal(t, tc.expectedStatusCode, w.Result().StatusCode)
	}
}

func TestHandler_getPostsByUserId(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		userId             string
		form               getPostsByUserIdForm
		expectedStatusCode int
	}{
		{testUserId, getPostsByUserIdForm{}, http.StatusOK},
		{testUserId, getPostsByUserIdForm{Offset: testOffset, Limit: testLimit}, http.StatusOK},
		{"", getPostsByUserIdForm{}, http.StatusNotFound},
		{"", getPostsByUserIdForm{Offset: testOffset, Limit: testLimit}, http.StatusNotFound},
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
