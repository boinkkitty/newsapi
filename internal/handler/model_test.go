package handler_test

import (
	"github.com/boinkkitty/newsapi/internal/handler"
	"testing"
)

func TestNewsPostReqBody_Validate(t *testing.T) {
	testCases := []struct {
		name            string
		req             handler.NewsPostReqBody
		isExpectedError bool
	}{
		{
			name:            "author empty",
			req:             handler.NewsPostReqBody{},
			isExpectedError: true,
		},
		{
			name:            "title empty",
			req:             handler.NewsPostReqBody{Author: "Alice"},
			isExpectedError: true,
		},
		{
			name: "summary empty",
			req: handler.NewsPostReqBody{
				Author: "Alice",
				Title:  "Some Title",
			},
			isExpectedError: true,
		},
		{
			name: "time invalid",
			req: handler.NewsPostReqBody{
				Author:  "Alice",
				Title:   "Some Title",
				Summary: "Some summary",
			},
			isExpectedError: true,
		},
		{
			name: "source invalid",
			req: handler.NewsPostReqBody{
				Author:    "Alice",
				Title:     "Some Title",
				Summary:   "Some summary",
				CreatedAt: "invalid time",
			},
			isExpectedError: true,
		},
		{
			name: "tags empty",
			req: handler.NewsPostReqBody{
				Author:    "Alice",
				Title:     "Some Title",
				Summary:   "Some summary",
				CreatedAt: "2023-01-01T00:00:00Z",
				Source:    "https://example.com",
			},
			isExpectedError: true,
		},
		{
			name: "validate",
			req: handler.NewsPostReqBody{
				Author:    "Alice",
				Title:     "Some Title",
				Summary:   "Some summary",
				CreatedAt: "2023-01-01T00:00:00Z",
				Source:    "https://example.com",
				Tags:      []string{"go", "news"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.req.Validate()

			if tc.isExpectedError && err == nil {
				t.Fatal("expected error but got nil")
			}

			if !tc.isExpectedError && err != nil {
				t.Fatalf("expected nil but got error: %v", err)
			}
		})
	}
}
