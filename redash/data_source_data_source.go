package redash

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redashgo "github.com/winebarrel/redash-go/v2"
)

func dataSourceDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: readDataSourceByName,
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

func readDataSourceByName(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*redashgo.Client)
	name := d.Get("name").(string)
	dsList, err := client.ListDataSources(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, ds := range dsList {
		if ds.Name == name {
			d.SetId(strconv.Itoa(ds.ID))
			// Always store options here. The resource read skips that when
			// options is unset, so options_wo credentials stay out of state.
			err := readDataSource0(ctx, d, meta, true)
			if err != nil {
				return diag.FromErr(err)
			}

			return nil
		}
	}

	return diag.Errorf("Data Source (%s) not found", name)
}
