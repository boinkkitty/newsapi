package handler_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/boinkkitty/newsapi/internal/handler"
	"github.com/boinkkitty/newsapi/internal/store"
	"github.com/google/uuid"
)

func Test_PostNews(t *testing.T) {
	testcases := []struct {
		name           string
		body           io.Reader
		store          handler.NewsStorer
		expectedStatus int
	}{
		{
			name:           "invalid request body json",
			body:           strings.NewReader(`{`),
			store:          mockNewsStore{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid request body",
			body: strings.NewReader(`{
				"id": "550e8400-e29b-41d4-a716-446655440000",
				"author": "Jane Doe",
				"title": "Go Makes Testing Easy",
				"content": "some content",
				"summary": "A short summary about Go testing",
				"created_at": "2024-07-16T15:04:05Z",
				"source": "https://example.com/go-testing"
			}`),
			store:          mockNewsStore{isExpectedError: true},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "db error",
			body: strings.NewReader(`{
				"id": "550e8400-e29b-41d4-a716-446655440000",
				"author": "Jane Doe",
				"title": "Go Makes Testing Easy",
				"content": "some content",
				"summary": "A short summary about Go testing",
				"created_at": "2024-07-16T15:04:05Z",
				"source": "https://example.com",
				"tags": ["go", "testing", "development"]
			}`),
			store:          mockNewsStore{isExpectedError: true},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "success",
			body: strings.NewReader(`{
				"id": "550e8400-e29b-41d4-a716-446655440000",
				"author": "Jane Doe",
				"title": "Go Makes Testing Easy",
				"content": "some content",
				"summary": "A short summary about Go testing",
				"created_at": "2024-07-16T15:04:05Z",
				"source": "https://example.com",
				"tags": ["go", "testing", "development"]
			}`),
			store:          mockNewsStore{},
			expectedStatus: http.StatusCreated,
		},
	}

	// Iterate through test cases
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", tc.body)
			response := httptest.NewRecorder()

			// Post Method
			handler.PostNews(tc.store)(response, request)

			// Assert
			assertStatusCode(t, response.Result().StatusCode, tc.expectedStatus)
		})
	}
}

func Test_GetAllNews(t *testing.T) {
	testcases := []struct {
		name           string
		store          handler.NewsStorer
		expectedStatus int
	}{
		{
			name:           "db error",
			store:          mockNewsStore{isExpectedError: true},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "success",
			store:          mockNewsStore{},
			expectedStatus: http.StatusOK,
		},
	}

	// Iterate through test cases
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
			response := httptest.NewRecorder()

			// Post Method
			handler.GetAllNews(tc.store)(response, request)

			// Assert
			assertStatusCode(t, response.Result().StatusCode, tc.expectedStatus)
		})
	}
}

func Test_GetNewsByID(t *testing.T) {
	testcases := []struct {
		name           string
		store          handler.NewsStorer
		newsID         string
		expectedStatus int
	}{
		{
			name:           "invalid news id",
			store:          mockNewsStore{},
			newsID:         "invalid-uuid",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "not found",
			store:          mockNewsStore{isExpectedError: true},
			newsID:         uuid.NewString(),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "success",
			store:          mockNewsStore{},
			newsID:         uuid.NewString(),
			expectedStatus: http.StatusOK,
		},
	}

	// Iterate through test cases
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
			response := httptest.NewRecorder()
			request.SetPathValue("news_id", tc.newsID)

			// Post Method
			handler.GetNewsByID(tc.store)(response, request)

			// Assert
			assertStatusCode(t, response.Result().StatusCode, tc.expectedStatus)
		})
	}
}

func Test_UpdateNewsByID(t *testing.T) {
	testcases := []struct {
		name           string
		body           io.Reader
		store          handler.NewsStorer
		newsID         string
		expectedStatus int
	}{
		{
			name:           "invalid request body json",
			body:           strings.NewReader(`{`),
			store:          mockNewsStore{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid request body",
			body: strings.NewReader(`{
				"id": "550e8400-e29b-41d4-a716-446655440000",
				"author": "Jane Doe",
				"title": "Go Makes Testing Easy",
				"content": "some content",
				"summary": "A short summary about Go testing",
				"created_at": "2024-07-16T15:04:05Z",
				"source": "https://example.com/go-testing"
			}`),
			store:          mockNewsStore{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "db error",
			body: strings.NewReader(`{
				"id": "550e8400-e29b-41d4-a716-446655440000",
				"author": "Jane Doe",
				"title": "Go Makes Testing Easy",
				"content": "some content",
				"summary": "A short summary about Go testing",
				"created_at": "2024-07-16T15:04:05Z",
				"source": "https://example.com",
				"tags": ["go", "testing", "development"]
			}`),
			store:          mockNewsStore{isExpectedError: true},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "success",
			body: strings.NewReader(`{
				"id": "550e8400-e29b-41d4-a716-446655440000",
				"author": "Jane Doe",
				"title": "Go Makes Testing Easy",
				"content": "some content",
				"summary": "A short summary about Go testing",
				"created_at": "2024-07-16T15:04:05Z",
				"source": "https://example.com",
				"tags": ["go", "testing", "development"]
			}`),
			store:          mockNewsStore{},
			expectedStatus: http.StatusOK,
		},
	}

	// Iterate through test cases
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPut, "/", tc.body)
			response := httptest.NewRecorder()

			// Post Method
			handler.UpdateNewsByID(tc.store)(response, request)

			// Assert
			assertStatusCode(t, response.Result().StatusCode, tc.expectedStatus)
		})
	}
}

func Test_DeleteNewsByID(t *testing.T) {
	testcases := []struct {
		name           string
		store          handler.NewsStorer
		newsID         string
		expectedStatus int
	}{
		{
			name:           "invalid news id",
			store:          mockNewsStore{},
			newsID:         "invalid-uuid",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "db error",
			store:          mockNewsStore{isExpectedError: true},
			newsID:         uuid.NewString(),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "success",
			store:          mockNewsStore{},
			newsID:         uuid.NewString(),
			expectedStatus: http.StatusNoContent,
		},
	}

	// Iterate through test cases
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodDelete, "/", http.NoBody)
			response := httptest.NewRecorder()
			// Set Id path value
			request.SetPathValue("news_id", tc.newsID)

			// Post Method
			handler.DeleteNewsByID(tc.store)(response, request)

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

type mockNewsStore struct {
	isExpectedError bool
}

func (m mockNewsStore) Create(_ store.News) (news store.News, err error) {
	if m.isExpectedError {
		return news, errors.New("some error")
	}
	return news, nil
}

func (m mockNewsStore) FindByID(_ uuid.UUID) (news store.News, err error) {
	if m.isExpectedError {
		return news, errors.New("some error")
	}
	return news, nil
}

func (m mockNewsStore) FindAll() (news []store.News, err error) {
	if m.isExpectedError {
		return news, errors.New("some error")
	}
	return news, nil
}

func (m mockNewsStore) DeleteByID(_ uuid.UUID) (err error) {
	if m.isExpectedError {
		return errors.New("some error")
	}
	return nil
}

func (m mockNewsStore) UpdateByID(_ store.News) (err error) {
	if m.isExpectedError {
		return errors.New("some error")
	}
	return nil
}
