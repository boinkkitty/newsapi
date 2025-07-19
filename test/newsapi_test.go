package test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestNewsAPI_(t *testing.T) {
	serviceAddr := forward(t, "news-service", "news-api-service", 8080)
	url := serviceAddr + "/news"
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, url, nil)
	require.NoError(t, err)

	client := http.DefaultClient
	res, err := client.Do(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
}
