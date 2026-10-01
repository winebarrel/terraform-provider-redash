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
	var gotHeader, gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeader = r.Header.Get("X-Request-Id")
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "{}")
	}))
	t.Cleanup(srv.Close)

	d := schema.TestResourceDataRaw(t, Provider().Schema, map[string]any{
		"url":     srv.URL,
		"api_key": "test-key",
		"http_headers": map[string]any{
			"X-Request-Id": "abc123",
		},
	})

	meta, diags := providerConfigure(context.Background(), d)
	require.False(t, diags.HasError())

	client := meta.(*redash_go.Client)
	res, closeBody, err := client.Get(context.Background(), "ping", nil)
	require.NoError(t, err)
	defer closeBody()
	defer res.Body.Close()

	require.Equal(t, "abc123", gotHeader)
	require.Equal(t, "Key test-key", gotAuth)
}

func TestProviderConfigure_withoutHTTPHeaders(t *testing.T) {
	var gotHeader, gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeader = r.Header.Get("X-Request-Id")
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "{}")
	}))
	t.Cleanup(srv.Close)

	d := schema.TestResourceDataRaw(t, Provider().Schema, map[string]any{
		"url":     srv.URL,
		"api_key": "test-key",
	})

	meta, diags := providerConfigure(context.Background(), d)
	require.False(t, diags.HasError())

	client := meta.(*redash_go.Client)
	res, closeBody, err := client.Get(context.Background(), "ping", nil)
	require.NoError(t, err)
	defer closeBody()
	defer res.Body.Close()

	require.Empty(t, gotHeader)
	require.Equal(t, "Key test-key", gotAuth)
}
