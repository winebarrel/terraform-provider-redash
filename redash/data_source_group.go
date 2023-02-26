package redash

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redashgo "github.com/winebarrel/redash-go"
)

func dataSourceGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Group data source.",
		ReadContext: readGroupByName,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func readGroupByName(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*redashgo.Client)
	name := d.Get("name").(string)
	groups, err := client.ListGroups(ctx)

	if err != nil {
		return diag.FromErr(err)
	}

	for _, group := range groups {
		if group.Name == name {
			d.SetId(strconv.Itoa(group.ID))
			return nil
		}
	}

	return diag.Errorf("Group (%s) not found", name)
}
