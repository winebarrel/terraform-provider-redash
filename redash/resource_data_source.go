package redash

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redashgo "github.com/winebarrel/redash-go/v2"
)

// Redash masks secret option values (e.g. "password") with this fixed
// placeholder in API responses. The real values are never returned.
// cf. https://github.com/getredash/redash/blob/master/redash/utils/configuration.py
const secretPlaceholder = "--------"

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
				Description:   "Data Source options (JSON string). Use `jsonencode()`. Stored in Terraform state. Conflicts with `options_wo` and `options_wo_version`.",
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"options_wo", "options_wo_version"},
			},
			"options_wo": {
				Description:   "Data Source options (JSON string) that are not stored in Terraform state. Use `jsonencode()`. Use this instead of `options` when the JSON contains credentials. Requires Terraform 1.11 or later. Conflicts with `options`. Pair with `options_wo_version`. Changing this value alone does not update the data source; increment `options_wo_version` to apply it. Redash does not return these values, so they are not refreshed into state.",
				Type:          schema.TypeString,
				Optional:      true,
				WriteOnly:     true,
				Sensitive:     true,
				RequiredWith:  []string{"options_wo_version"},
				ConflictsWith: []string{"options"},
			},
			"options_wo_version": {
				Description:   "Version of `options_wo`, stored in Terraform state. Increment this when `options_wo` should be applied. Required with `options_wo`. Conflicts with `options`.",
				Type:          schema.TypeInt,
				Optional:      true,
				RequiredWith:  []string{"options_wo"},
				ConflictsWith: []string{"options"},
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

	options, ok, diags := dataSourceOptions(d)
	if diags.HasError() {
		return diags
	}
	if ok {
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
	// Refresh does not include write-only values. When options is unset, leave
	// it unset so options_wo is not copied from the API into state.
	_, storeOptions := d.GetOk("options")

	err := readDataSource0(ctx, d, meta, storeOptions)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func readDataSource0(ctx context.Context, d *schema.ResourceData, meta any, storeOptions bool) error {
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

	if !storeOptions {
		return nil
	}

	// The API masks secret values with a placeholder. Writing the placeholder
	// to the state would cause permanent drift against the configuration, so
	// keep the values already in the state for masked keys.
	// cf. https://github.com/winebarrel/terraform-provider-redash/issues/177
	if v, ok := d.GetOk("options"); ok {
		stateOptions := map[string]any{}

		if err := json.Unmarshal([]byte(v.(string)), &stateOptions); err == nil {
			for key, value := range ds.Options {
				if value == secretPlaceholder {
					if stateValue, ok := stateOptions[key]; ok {
						ds.Options[key] = stateValue
					}
				}
			}
		}
	}

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

	// Send options_wo whenever it is configured, including when only name or
	// type changed. The client serializes a nil options map as null, which
	// would clear credentials. options_wo itself never shows a plan diff;
	// options_wo_version is the stored trigger for applying a new value.
	options, ok, diags := dataSourceOptions(d)
	if diags.HasError() {
		return diags
	}
	if ok {
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
	// Import records the API options, including masked secrets.
	err := readDataSource0(ctx, d, meta, true)
	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

// dataSourceOptions returns options from options_wo when it is set, otherwise
// from options. options_wo is read from the raw config because d.Get does not
// return write-only values.
func dataSourceOptions(d *schema.ResourceData) (map[string]any, bool, diag.Diagnostics) {
	options, ok, diags := dataSourceOptionsWriteOnly(d)
	if diags.HasError() || ok {
		return options, ok, diags
	}

	v, ok := d.GetOk("options")
	if !ok {
		return nil, false, nil
	}

	options, err := unmarshalOptions(v.(string))
	if err != nil {
		return nil, false, diag.FromErr(err)
	}

	return options, true, nil
}

func dataSourceOptionsWriteOnly(d *schema.ResourceData) (map[string]any, bool, diag.Diagnostics) {
	v, diags := d.GetRawConfigAt(cty.GetAttrPath("options_wo"))
	if diags.HasError() {
		return nil, false, diags
	}
	if v.IsNull() {
		return nil, false, nil
	}
	if !v.IsKnown() {
		return nil, false, diag.Errorf("options_wo is unknown")
	}

	options, err := unmarshalOptions(v.AsString())
	if err != nil {
		return nil, false, diag.FromErr(err)
	}

	return options, true, nil
}

func unmarshalOptions(s string) (map[string]any, error) {
	options := map[string]any{}
	if err := json.Unmarshal([]byte(s), &options); err != nil {
		return nil, err
	}

	return options, nil
}
