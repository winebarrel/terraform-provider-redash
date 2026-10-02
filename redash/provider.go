package redash

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redash_go "github.com/winebarrel/redash-go/v2"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Description: "Redash API endpoint URL. This can also be set from the REDASH_URL environment variable.",
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("REDASH_URL", nil),
			},
			"api_key": {
				Description: "Redash User API Key. This can also be set from the REDASH_API_KEY environment variable.",
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("REDASH_API_KEY", nil),
				Sensitive:   true,
			},
			"http_headers": {
				Description: "Extra HTTP headers sent on every Redash API request.",
				Type:        schema.TypeMap,
				Optional:    true,
				Sensitive:   true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		ConfigureContextFunc: providerConfigure,
		ResourcesMap: map[string]*schema.Resource{
			"redash_alert_destination":  resourceAlertDestination(),
			"redash_alert_subscription": resourceAlertSubscription(),
			"redash_alert":              resourceAlert(),
			"redash_data_source":        resourceDataSource(),
			"redash_group_data_source":  resourceGroupDataSource(),
			"redash_group_member":       resourceGroupMember(),
			"redash_group":              resourceGroup(),
			"redash_query":              resourceQuery(),
			"redash_user":               resourceUser(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"redash_alert_destination":  dataSourceAlertDestination(),
			"redash_alert":              dataSourceAlert(),
			"redash_data_source":        dataSourceDataSource(),
			"redash_group_data_sources": dataSourceGroupDataSources(),
			"redash_group":              dataSourceGroup(),
			"redash_query":              dataSourceQuery(),
			"redash_user":               dataSourceUser(),
			"redash_users":              dataSourceUsers(),
		},
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
	url := d.Get("url").(string)

	if url == "" {
		return nil, diag.Errorf("url is required")
	}

	apiKey := d.Get("api_key").(string)

	if apiKey == "" {
		return nil, diag.Errorf("api_key is required")
	}

	var httpClient *http.Client
	if headers := httpHeaders(d); len(headers) > 0 {
		httpClient = &http.Client{Transport: headerTransport(headers)}
	}

	client, err := redash_go.NewClientWithHTTPClient(url, apiKey, httpClient)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	client.SetDebug(logging.IsDebugOrHigher())

	return client, nil
}

// httpHeaders reads the optional http_headers provider argument.
// The Terraform SDK stores map values as map[string]any.
// An unset argument returns nil, and the Redash client then uses http.DefaultClient.
func httpHeaders(d *schema.ResourceData) map[string]string {
	raw, ok := d.GetOk("http_headers")
	if !ok {
		return nil
	}

	in := raw.(map[string]any)
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v.(string)
	}
	return out
}

// headerTransport sets configured headers on each request. RoundTrip must not
// modify the incoming request, so it clones it first. The Redash client sets
// Authorization before http.Client.Do, which then runs RoundTrip, so that
// header is already on the request and is left in place unless http_headers
// also sets it.
type headerTransport map[string]string

func (t headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	for k, v := range t {
		req.Header.Set(k, v)
	}

	return http.DefaultTransport.RoundTrip(req)
}
