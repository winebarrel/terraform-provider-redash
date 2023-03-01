package redash

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redash_go "github.com/winebarrel/redash-go"
)

func resourceDataSource() *schema.Resource {
	return &schema.Resource{
		CreateContext: createDataSource,
		ReadContext:   readDataSource,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"options": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func createDataSource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*redash_go.Client)

	input := &redash_go.CreateDataSourceInput{
		Name: d.Get("name").(string),
		Type: d.Get("type").(string),
	}

	if v, ok := d.GetOk("options"); ok {
		input.Options = v.(map[string]any)
	}

	ds, err := client.CreateDataSource(ctx, input)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(ds.ID))

	return readDataSource(ctx, d, meta)
}

func readDataSource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}
