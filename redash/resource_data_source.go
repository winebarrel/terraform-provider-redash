package redash

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redashgo "github.com/winebarrel/redash-go/v2"
)

func resourceDataSource() *schema.Resource {
	return &schema.Resource{
		CreateContext: createDataSource,
		ReadContext:   readDataSource,
		UpdateContext: updateDataSource,
		DeleteContext: deleteDataSource,
		Importer: &schema.ResourceImporter{
			StateContext: importDataSource,
		},
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
				Description: "Data Source options (JSON string). Use `jsonencode()`.",
				Type:        schema.TypeString,
				Optional:    true,
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
		options := map[string]any{}
		err := json.Unmarshal([]byte(v.(string)), &options)

		if err != nil {
			return diag.FromErr(err)
		}

		input.Options = options
	}

	ds, err := client.CreateDataSource(ctx, input)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(ds.ID))

	return readDataSource(ctx, d, meta)
}

func readDataSource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	err := readDataSource0(ctx, d, meta)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func readDataSource0(ctx context.Context, d *schema.ResourceData, meta any) error {
	id, err := strconv.Atoi(d.Id())

	if err != nil {
		return err
	}

	client := meta.(*redashgo.Client)
	ds, err := client.GetDataSource(ctx, id)

	if err != nil {
		return err
	}

	d.Set("name", ds.Name) //nolint:errcheck
	d.Set("type", ds.Type) //nolint:errcheck

	options, err := json.Marshal(ds.Options)

	if err != nil {
		return err
	}

	d.Set("options", string(options)) //nolint:errcheck

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
		options := map[string]any{}
		err := json.Unmarshal([]byte(v.(string)), &options)

		if err != nil {
			return diag.FromErr(err)
		}

		input.Options = options
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

func importDataSource(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	err := readDataSource0(ctx, d, meta)

	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
