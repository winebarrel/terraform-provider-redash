package redash

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redashgo "github.com/winebarrel/redash-go"
)

func dataSourceAlert() *schema.Resource {
	return &schema.Resource{
		ReadContext: readAlertByName,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"query_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"options": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"column": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"op": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"custom_subject": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"custom_body": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template": {
							Type:       schema.TypeString,
							Computed:   true,
							Deprecated: "This attribute is for backward compatibility.",
						},
					},
				},
			},
			"rearm": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"muted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func readAlertByName(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*redashgo.Client)
	name := d.Get("name").(string)
	alerts, err := client.ListAlerts(ctx)

	if err != nil {
		return diag.FromErr(err)
	}

	for _, alert := range alerts {
		if alert.Name == name {
			d.SetId(strconv.Itoa(alert.ID))
			return readAlert(ctx, d, meta)
		}
	}

	return diag.Errorf("Alert (%s) not found", name)
}
