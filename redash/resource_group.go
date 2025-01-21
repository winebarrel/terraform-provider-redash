package redash

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redashgo "github.com/winebarrel/redash-go/v2"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: createGroup,
		ReadContext:   schema.NoopContext,
		DeleteContext: deleteGroup,
		Importer: &schema.ResourceImporter{
			StateContext: importGroup,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func createGroup(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*redashgo.Client)

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

func deleteGroup(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, _ := strconv.Atoi(d.Id())
	client := meta.(*redashgo.Client)

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
