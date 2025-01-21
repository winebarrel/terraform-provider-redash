package redash

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redashgo "github.com/winebarrel/redash-go/v2"
)

func dataSourceAlertDestination() *schema.Resource {
	return &schema.Resource{
		ReadContext: readAlertDestinationByName,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"options": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func readAlertDestinationByName(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*redashgo.Client)
	name := d.Get("name").(string)
	dests, err := client.ListDestinations(ctx)

	if err != nil {
		return diag.FromErr(err)
	}

	for _, dest := range dests {
		if dest.Name == name {
			d.SetId(strconv.Itoa(dest.ID))
			return readAlertDestination(ctx, d, meta)
		}
	}

	return diag.Errorf("Alert Destination (%s) not found", name)
}
