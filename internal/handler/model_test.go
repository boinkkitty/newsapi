package handler_test

import (
	"github.com/boinkkitty/newsapi/internal/news"
	"net/url"
	"testing"
	"time"

	"github.com/boinkkitty/newsapi/internal/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewsPostReqBody_Validate(t *testing.T) {
	type expectations struct {
		err  string
		news *news.Record
	}
	testCases := []struct {
		name         string
		req          handler.NewsPostReqBody
		expectations expectations
	}{
		{
			name: "author empty",
			req:  handler.NewsPostReqBody{},
			expectations: expectations{
				err: "author is empty",
			},
		},
		{
			name: "title empty",
			req:  handler.NewsPostReqBody{Author: "Alice"},
			expectations: expectations{
				err: "title is empty",
			},
		},
		{
			name: "content empty",
			req: handler.NewsPostReqBody{
				Author: "Alice",
				Title:  "Some Title",
			},
			expectations: expectations{
				err: "content is empty",
			},
		},
		{
			name: "summary empty",
			req: handler.NewsPostReqBody{
				Author:  "Alice",
				Title:   "Some Title",
				Content: "Some Content",
			},
			expectations: expectations{
				err: "summary is empty",
			},
		},
		{
			name: "time invalid",
			req: handler.NewsPostReqBody{
				Author:    "Alice",
				Title:     "Some Title",
				Content:   "Some Content",
				Summary:   "Some summary",
				CreatedAt: "invalid time",
			},
			expectations: expectations{
				err: `parsing time "invalid time"`,
			},
		},
		{
			name: "source invalid",
			req: handler.NewsPostReqBody{
				Author:    "Alice",
				Title:     "Some Title",
				Content:   "Some Content",
				Summary:   "Some summary",
				CreatedAt: "2023-01-01T00:00:00Z",
			},
			expectations: expectations{
				err: "source is empty",
			},
		},
		{
			name: "tags empty",
			req: handler.NewsPostReqBody{
				Author:    "Alice",
				Title:     "Some Title",
				Content:   "Some Content",
				Summary:   "Some summary",
				CreatedAt: "2023-01-01T00:00:00Z",
				Source:    "https://example.com",
			},
			expectations: expectations{
				err: "tags cannot be empty",
			},
		},
		{
			name: "validate",
			req: handler.NewsPostReqBody{
				Author:    "Alice",
				Title:     "Some Title",
				Content:   "Some Content",
				Summary:   "Some summary",
				CreatedAt: "2023-01-01T00:00:00Z",
				Source:    "https://example.com",
				Tags:      []string{"go", "news"},
			},
			expectations: expectations{
				news: &news.Record{
					Author:  "Alice",
					Title:   "Some Title",
					Content: "Some Content",
					Summary: "Some summary",
					Tags:    []string{"go", "news"},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			news, err := tc.req.Validate()

			if tc.expectations.err != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectations.err)
			} else {
				assert.NoError(t, err)

				// Check time
				parseTime, parseErr := time.Parse(time.RFC3339, tc.req.CreatedAt)
				require.NoError(t, parseErr)
				tc.expectations.news.CreatedAt = parseTime

				// Check url
				parsedSource, parseSourceErr := url.Parse(tc.req.Source)
				require.NoError(t, parseSourceErr)
				tc.expectations.news.Source = parsedSource.String()
				assert.Equal(t, tc.expectations.news, news)
			}
		})
	}
}
