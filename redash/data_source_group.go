package redash

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGroup() *schema.Resource {
	return &schema.Resource{
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
	client := meta.(*providerMeta).client
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
