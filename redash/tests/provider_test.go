package test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/winebarrel/terraform-provider-redash/redash"
)

var (
	testAccProviders map[string]*schema.Provider
	testAccProvider  *schema.Provider
)

const (
	testAccRedashURL    = "http://localhost:5001"
	testAccRedashAPIKey = "G1LARLeRTzoWF7asyy32Qdvken2OZ2LhzoOzwA3r"
)

func init() {
	testAccProvider = redash.Provider()
	testAccProviders = map[string]*schema.Provider{
		"redash": testAccProvider,
	}

}

func testAccPreCheck(t *testing.T) {
	t.Setenv("REDASH_URL", testAccRedashURL)
	t.Setenv("REDASH_API_KEY", testAccRedashAPIKey)
}

func TestProvider(t *testing.T) {
	assert := assert.New(t)

	provider := redash.Provider()
	err := provider.InternalValidate()
	assert.NoError(err)

	raw := map[string]interface{}{
		"url":     "https://example.com",
		"api_key": "api_key",
	}

	diagnostics := provider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	assert.False(diagnostics.HasError())
}

func TestProvider_withoutURL(t *testing.T) {
	assert := assert.New(t)

	provider := redash.Provider()
	err := provider.InternalValidate()
	assert.NoError(err)

	raw := map[string]interface{}{
		"api_key": "api_key",
	}

	diagnostics := provider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	assert.True(diagnostics.HasError())
	assert.Equal("url is required", diagnostics[0].Summary)
}

func TestProvider_withoutAPIKey(t *testing.T) {
	assert := assert.New(t)

	provider := redash.Provider()
	err := provider.InternalValidate()
	assert.NoError(err)

	raw := map[string]interface{}{
		"url": "https://example.com",
	}

	diagnostics := provider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	assert.True(diagnostics.HasError())
	assert.Equal("api_key is required", diagnostics[0].Summary)
}
