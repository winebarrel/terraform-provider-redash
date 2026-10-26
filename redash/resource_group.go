package redash

import (
	"context"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redashgo "github.com/winebarrel/redash-go/v2"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: createGroup,
		ReadContext:   readGroup,
		UpdateContext: updateGroup,
		DeleteContext: deleteGroup,
		Importer: &schema.ResourceImporter{
			StateContext: importGroup,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func createGroup(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*providerMeta).client

	input := &redashgo.CreateGroupInput{
		Name: d.Get("name").(string),
	}

	query, err := client.CreateGroup(ctx, input)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(query.ID))

	return nil
}

func readGroup(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := meta.(*providerMeta).client
	group, err := client.GetGroup(ctx, id)
	if err != nil {
		// Redash answers GET /api/groups/:id with 500, not 404, for a deleted group.
		if strings.Contains(err.Error(), "HTTP status code not OK: 500") {
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	d.Set("name", group.Name) //nolint:errcheck

	return nil
}

func updateGroup(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, _ := strconv.Atoi(d.Id())
	client := meta.(*providerMeta).client

	input := &redashgo.UpdateGroupInput{
		Name: d.Get("name").(string),
	}

	_, err := client.UpdateGroup(ctx, id, input)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func deleteGroup(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, _ := strconv.Atoi(d.Id())
	client := meta.(*providerMeta).client

	err := client.DeleteGroup(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func importGroup(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	_, err := strconv.Atoi(d.Id())
	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
