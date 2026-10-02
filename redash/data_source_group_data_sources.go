package redash

import (
	"context"
	"slices"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redashgo "github.com/winebarrel/redash-go/v2"
)

func dataSourceGroupDataSources() *schema.Resource {
	return &schema.Resource{
		ReadContext: readGroupDataSources,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Description: "ID of the group.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"data_sources": {
				Description: "Data sources granted to the group.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"view_only": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func readGroupDataSources(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*redashgo.Client)
	groupId := d.Get("group_id").(int)

	dsList, err := client.ListGroupDataSources(ctx, groupId)
	if err != nil {
		return diag.FromErr(err)
	}

	slices.SortFunc(dsList, func(a, b redashgo.DataSource) int {
		return a.ID - b.ID
	})

	dataSources := make([]map[string]any, 0, len(dsList))
	for _, ds := range dsList {
		dataSources = append(dataSources, map[string]any{
			"id":        ds.ID,
			"name":      ds.Name,
			"view_only": ds.ViewOnly,
		})
	}

	d.SetId(strconv.Itoa(groupId))
	d.Set("data_sources", dataSources) //nolint:errcheck

	return nil
}
