package test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/terraform-provider-redash/redash"
)

var (
	testAccProviderFactories map[string]func() (*schema.Provider, error)
	testAccProvider          *schema.Provider
)

const (
	testAccRedashURL    = "http://localhost:5001"
	testAccRedashAPIKey = "6nh64ZsT66WeVJvNZ6WB5D2JKZULeC2VBdSD68wt"
)

func init() {
	testAccProvider = redash.Provider()
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"redash": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}

func testAccPreCheck(t *testing.T) {
	t.Setenv("REDASH_URL", testAccRedashURL)
	t.Setenv("REDASH_API_KEY", testAccRedashAPIKey)
}

// newTestRedash answers GET /api/session like Redash: 200 for apiKey, and for
// any other key the 404 that Redash returns to a request it cannot log in.
func newTestRedash(t *testing.T, apiKey string) string {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Key "+apiKey {
			w.WriteHeader(http.StatusNotFound)
			_, _ = io.WriteString(w, `{"message": "Couldn't find resource. Please login and try again."}`)
			return
		}

		if r.Method != http.MethodGet || r.URL.Path != "/api/session" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		_, _ = io.WriteString(w, `{"org_slug": "default"}`)
	}))
	t.Cleanup(srv.Close)

	return srv.URL
}

func TestProvider(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	provider := redash.Provider()
	err := provider.InternalValidate()
	require.NoError(err)

	raw := map[string]any{
		"url":     newTestRedash(t, "api_key"),
		"api_key": "api_key",
	}

	diagnostics := provider.Configure(t.Context(), terraform.NewResourceConfigRaw(raw))
	assert.False(diagnostics.HasError())
}

func TestProvider_wrongAPIKey(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	provider := redash.Provider()
	err := provider.InternalValidate()
	require.NoError(err)

	raw := map[string]any{
		"url":     newTestRedash(t, "api_key"),
		"api_key": "wrong_key",
	}

	diagnostics := provider.Configure(t.Context(), terraform.NewResourceConfigRaw(raw))
	require.True(diagnostics.HasError())
	assert.Contains(diagnostics[0].Summary, "Redash rejected the API key")
}

func TestProvider_withoutURL(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	provider := redash.Provider()
	err := provider.InternalValidate()
	require.NoError(err)

	raw := map[string]any{
		"api_key": "api_key",
	}

	diagnostics := provider.Configure(t.Context(), terraform.NewResourceConfigRaw(raw))
	assert.True(diagnostics.HasError())
	assert.Equal("url is required", diagnostics[0].Summary)
}

func TestProvider_withoutAPIKey(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	provider := redash.Provider()
	err := provider.InternalValidate()
	require.NoError(err)

	raw := map[string]any{
		"url": "https://example.com",
	}

	diagnostics := provider.Configure(t.Context(), terraform.NewResourceConfigRaw(raw))
	assert.True(diagnostics.HasError())
	assert.Equal("api_key is required", diagnostics[0].Summary)
}

// With a wrong API key every read answers 404, the same status as a deleted
// object, so the plan must fail instead of dropping the data source from the
// state and planning to create it again.
func TestAccProvider_wrongAPIKeyKeepsState(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceConfigBasicPg,
			},
			{
				Config:      `provider "redash" { api_key = "wrong-key" }` + testAccDataSourceConfigBasicPg,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Redash rejected the API key`),
			},
			{
				Config:   testAccDataSourceConfigBasicPg,
				PlanOnly: true,
			},
		},
	})
}
