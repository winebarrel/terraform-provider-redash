package redash

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
	redash_go "github.com/winebarrel/redash-go/v2"
)

func TestProviderConfigure_httpHeaders(t *testing.T) {
	tests := []struct {
		name        string
		httpHeaders map[string]any
		want        map[string]string
	}{
		{
			name:        "unset",
			httpHeaders: nil,
			want: map[string]string{
				"Authorization": "Key test-key",
				"X-Request-Id":  "",
			},
		},
		{
			name: "single header",
			httpHeaders: map[string]any{
				"X-Request-Id": "abc123",
			},
			want: map[string]string{
				"Authorization": "Key test-key",
				"X-Request-Id":  "abc123",
			},
		},
		{
			name: "multiple headers",
			httpHeaders: map[string]any{
				"X-Request-Id":        "abc123",
				"Proxy-Authorization": "Bearer jwt",
			},
			want: map[string]string{
				"Authorization":       "Key test-key",
				"X-Request-Id":        "abc123",
				"Proxy-Authorization": "Bearer jwt",
			},
		},
		{
			name: "override authorization",
			httpHeaders: map[string]any{
				"Authorization": "Bearer custom",
			},
			want: map[string]string{
				"Authorization": "Bearer custom",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got http.Header
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				got = r.Header.Clone()
				w.WriteHeader(http.StatusOK)
				_, _ = io.WriteString(w, "{}")
			}))
			t.Cleanup(srv.Close)

			raw := map[string]any{
				"url":     srv.URL,
				"api_key": "test-key",
			}
			if tt.httpHeaders != nil {
				raw["http_headers"] = tt.httpHeaders
			}
			d := schema.TestResourceDataRaw(t, Provider().Schema, raw)

			meta, diags := providerConfigure(context.Background(), d)
			require.False(t, diags.HasError())

			client := meta.(*redash_go.Client)
			res, closeBody, err := client.Get(context.Background(), "ping", nil)
			require.NoError(t, err)
			defer closeBody()
			defer res.Body.Close()

			for k, v := range tt.want {
				require.Equal(t, v, got.Get(k), k)
			}
		})
	}
}

func TestHeaderTransport_doesNotModifyRequest(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)

	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Key test-key")

	res, err := headerTransport{
		"Authorization": "Bearer custom",
		"X-Request-Id":  "abc123",
	}.RoundTrip(req)
	require.NoError(t, err)
	defer res.Body.Close()

	require.Equal(t, http.Header{"Authorization": {"Key test-key"}}, req.Header)
}
