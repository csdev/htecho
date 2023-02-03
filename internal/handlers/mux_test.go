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

func TestMuxOptions(t *testing.T) {
	tests := []struct {
		description     string
		options         *Options
		requestHeaders  map[string]string
		responseHeaders map[string]string
	}{
		{
			description: "auth headers are excluded by default",
			options:     DefaultOptions(),
			requestHeaders: map[string]string{
				"authorization":       "bearer asdf",
				"proxy-authorization": "bearer zxcv",
			},
			responseHeaders: map[string]string{
				"X-Echo-Header-Authorization":       "",
				"X-Echo-Header-Proxy-Authorization": "",
			},
		},
		{
			description: "auth headers can be included",
			options:     &Options{IncludeAuth: true},
			requestHeaders: map[string]string{
				"authorization":       "bearer asdf",
				"proxy-authorization": "bearer zxcv",
			},
			responseHeaders: map[string]string{
				"X-Echo-Header-Authorization":       "bearer asdf",
				"X-Echo-Header-Proxy-Authorization": "bearer zxcv",
			},
		},
		{
			description: "ip headers are excluded by default",
			options:     DefaultOptions(),
			requestHeaders: map[string]string{
				"x-forwarded-for": "127.0.0.1",
				"forwarded":       "proto=http",
			},
			responseHeaders: map[string]string{
				"X-Echo-Header-X-Forwarded-For": "",
				"X-Echo-Header-Forwarded":       "",
			},
		},
		{
			description: "ip headers can be included",
			options:     &Options{IncludeIps: true},
			requestHeaders: map[string]string{
				"x-forwarded-for": "127.0.0.1",
				"forwarded":       "proto=http",
			},
			responseHeaders: map[string]string{
				"X-Echo-Header-X-Forwarded-For": "127.0.0.1",
				"X-Echo-Header-Forwarded":       "proto=http",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			mux := NewMux(test.options)
			req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(""))
			require.NoError(t, err)

			for name, value := range test.requestHeaders {
				req.Header.Add(name, value)
			}

			resp := httptest.NewRecorder()
			mux.ServeHTTP(resp, req)

			for name, value := range test.responseHeaders {
				assert.Equal(t, value, resp.Header().Get(name))
			}
		})
	}
}
