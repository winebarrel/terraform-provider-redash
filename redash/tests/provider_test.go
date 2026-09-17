package test

import (
	"testing"

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

func TestProvider(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	provider := redash.Provider()
	err := provider.InternalValidate()
	require.NoError(err)

	raw := map[string]any{
		"url":     "https://example.com",
		"api_key": "api_key",
	}

	diagnostics := provider.Configure(t.Context(), terraform.NewResourceConfigRaw(raw))
	assert.False(diagnostics.HasError())
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
