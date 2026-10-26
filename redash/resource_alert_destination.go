package redash

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redashgo "github.com/winebarrel/redash-go/v2"
)

func resourceAlertDestination() *schema.Resource {
	return &schema.Resource{
		CreateContext: createAlertDestination,
		ReadContext:   readAlertDestination,
		DeleteContext: deleteAlertDestination,
		Importer: &schema.ResourceImporter{
			StateContext: importAlertDestination,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"options": {
				Description: "Alert Destination options (JSON string). Use `jsonencode()`.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
		},
	}
}

func createAlertDestination(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*providerMeta).client

	input := &redashgo.CreateDestinationInput{
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

	dest, err := client.CreateDestination(ctx, input)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(dest.ID))

	return readAlertDestination(ctx, d, meta)
}

func readAlertDestination(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	err := readAlertDestination0(ctx, d, meta)
	if err != nil {
		// Redash answers GET /api/destinations/:id with 500, not 404, for a
		// deleted destination, so check the list before dropping it.
		if gone, listErr := alertDestinationGone(ctx, d, meta); listErr == nil && gone {
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	return nil
}

func readAlertDestination0(ctx context.Context, d *schema.ResourceData, meta any) error {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	client := meta.(*providerMeta).client
	dest, err := client.GetDestination(ctx, id)
	if err != nil {
		return err
	}

	d.Set("name", dest.Name) //nolint:errcheck
	d.Set("type", dest.Type) //nolint:errcheck

	options, err := json.Marshal(dest.Options)
	if err != nil {
		return err
	}

	d.Set("options", string(options)) //nolint:errcheck

	return nil
}

func deleteAlertDestination(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, _ := strconv.Atoi(d.Id())
	client := meta.(*providerMeta).client

	err := client.DeleteDestination(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func importAlertDestination(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	err := readAlertDestination0(ctx, d, meta)
	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func alertDestinationGone(ctx context.Context, d *schema.ResourceData, meta any) (bool, error) {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, err
	}

	client := meta.(*providerMeta).client
	dests, err := client.ListDestinations(ctx)
	if err != nil {
		return false, err
	}

	for _, dest := range dests {
		if dest.ID == id {
			return false, nil
		}
	}

	return true, nil
}
