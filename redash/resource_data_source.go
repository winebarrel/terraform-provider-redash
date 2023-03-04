package redash

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redashgo "github.com/winebarrel/redash-go"
)

func resourceDataSource() *schema.Resource {
	return &schema.Resource{
		CreateContext: createDataSource,
		ReadContext:   readDataSource,
		UpdateContext: updateDataSource,
		DeleteContext: deleteDataSource,
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
	client := meta.(*redashgo.Client)

	input := &redashgo.CreateDataSourceInput{
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
	id, _ := strconv.Atoi(d.Id())
	client := meta.(*redashgo.Client)
	ds, err := client.GetDataSource(ctx, id)

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", ds.Name)
	d.Set("type", ds.Type)
	d.Set("options", ds.Options)

	return nil
}

func updateDataSource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, _ := strconv.Atoi(d.Id())
	client := meta.(*redashgo.Client)

	input := &redashgo.UpdateDataSourceInput{
		Name: d.Get("name").(string),
		Type: d.Get("type").(string),
	}

	if v, ok := d.GetOk("options"); ok {
		input.Options = v.(map[string]any)
	}

	_, err := client.UpdateDataSource(ctx, id, input)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func deleteDataSource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, _ := strconv.Atoi(d.Id())
	client := meta.(*redashgo.Client)

	err := client.DeleteDataSource(ctx, id)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
