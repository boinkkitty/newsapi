package handler_test

import (
	"github.com/boinkkitty/newsapi/internal/handler"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_PostNews(t *testing.T) {
	testcases := []struct {
		name           string
		expectedStatus int
	}{
		{
			name:           "not implemented",
			expectedStatus: http.StatusNotImplemented,
		},
	}

	// Iterate through test cases
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", nil)
			response := httptest.NewRecorder()

			// Post Method
			handler.PostNews()(response, request)

			// Assert
			assertStatusCode(t, response.Result().StatusCode, tc.expectedStatus)
		})
	}
}

func Test_GetAllNews(t *testing.T) {
	testcases := []struct {
		name           string
		expectedStatus int
	}{
		{
			name:           "not implemented",
			expectedStatus: http.StatusNotImplemented,
		},
	}

	// Iterate through test cases
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			response := httptest.NewRecorder()

			// Post Method
			handler.GetAllNews()(response, request)

			// Assert
			assertStatusCode(t, response.Result().StatusCode, tc.expectedStatus)
		})
	}
}

func Test_GetNewsByID(t *testing.T) {
	testcases := []struct {
		name           string
		expectedStatus int
	}{
		{
			name:           "not implemented",
			expectedStatus: http.StatusNotImplemented,
		},
	}

	// Iterate through test cases
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			response := httptest.NewRecorder()

			// Post Method
			handler.GetNewsByID()(response, request)

			// Assert
			assertStatusCode(t, response.Result().StatusCode, tc.expectedStatus)
		})
	}
}

func Test_UpdateNewsByID(t *testing.T) {
	testcases := []struct {
		name           string
		expectedStatus int
	}{
		{
			name:           "not implemented",
			expectedStatus: http.StatusNotImplemented,
		},
	}

	// Iterate through test cases
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPut, "/", nil)
			response := httptest.NewRecorder()

			// Post Method
			handler.UpdateNewsByID()(response, request)

			// Assert
			assertStatusCode(t, response.Result().StatusCode, tc.expectedStatus)
		})
	}
}

func Test_DeleteNewsByID(t *testing.T) {
	testcases := []struct {
		name           string
		expectedStatus int
	}{
		{
			name:           "not implemented",
			expectedStatus: http.StatusNotImplemented,
		},
	}

	// Iterate through test cases
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodDelete, "/", nil)
			response := httptest.NewRecorder()

			// Post Method
			handler.DeleteNewsByID()(response, request)

			// Assert
			assertStatusCode(t, response.Result().StatusCode, tc.expectedStatus)
		})
	}
}

func assertStatusCode(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("got: %d, want: %d", got, want)
	}
}
