package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMux(t *testing.T) {
	mux := NewMux(DefaultOptions())

	t.Run("it handles a request", func(t *testing.T) {
		const body = `{"key": "value"}`
		req, err := http.NewRequest(http.MethodPost, "/path/?query=value", strings.NewReader(body))
		require.NoError(t, err)

		req.Header.Add("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		mux.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, http.MethodPost, resp.Header().Get("X-Echo-Method"))
		assert.Equal(t, "query=value", resp.Header().Get("X-Echo-Query"))
		assert.Equal(t, "application/json", resp.Header().Get("X-Echo-Header-Content-Type"))
		assert.Equal(t, "application/json", resp.Header().Get("Content-Type"))

		data, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, body, string(data))
	})

	t.Run("it omits the content-type if the request does not have one", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(""))
		require.NoError(t, err)

		resp := httptest.NewRecorder()
		mux.ServeHTTP(resp, req)

		assert.Equal(t, "", resp.Header().Get("Content-Type"))
	})

	t.Run("it excludes sensitive headers", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(""))
		require.NoError(t, err)

		req.Header.Add("authorization", "bearer asdf")
		req.Header.Add("proxy-authorization", "bearer zxcv")

		resp := httptest.NewRecorder()
		mux.ServeHTTP(resp, req)

		assert.Equal(t, "", resp.Header().Get("Authorization"))
		assert.Equal(t, "", resp.Header().Get("Proxy-Authorization"))
	})

	t.Run("it echoes duplicate headers", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(""))
		require.NoError(t, err)

		req.Header.Add("X-Foo", "bar")
		req.Header.Add("X-Foo", "baz")

		resp := httptest.NewRecorder()
		mux.ServeHTTP(resp, req)

		assert.Equal(t, []string{"bar", "baz"}, resp.Header().Values("X-Echo-Header-X-Foo"))
	})
}
