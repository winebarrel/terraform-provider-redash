package redash

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redash_go "github.com/winebarrel/redash-go"
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
			"redash_alert_destination": dataSourceAlertDestination(),
			"redash_alert":             dataSourceAlert(),
			"redash_data_source":       dataSourceDataSource(),
			"redash_group":             dataSourceGroup(),
			"redash_query":             dataSourceQuery(),
			"redash_user":              dataSourceUser(),
			"redash_users":             dataSourceUsers(),
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

	client, err := redash_go.NewClient(url, apiKey)

	if err != nil {
		return nil, diag.FromErr(err)
	}

	client.Debug = logging.IsDebugOrHigher()

	return client, nil
}
