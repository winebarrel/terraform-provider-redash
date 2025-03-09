package redash

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redashgo "github.com/winebarrel/redash-go/v2"
)

func dataSourceQuery() *schema.Resource {
	return &schema.Resource{
		ReadContext: readQueryByName,
		Schema: map[string]*schema.Schema{
			"data_source_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"query": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"published": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"schedule": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func readQueryByName(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*redashgo.Client)
	name := d.Get("name").(string)

	input := &redashgo.ListQueriesInput{
		Q:        name,
		Page:     1,
		PageSize: 20,
	}

	for {
		rs, err := client.ListQueries(ctx, input)

		if err != nil {
			return diag.FromErr(err)
		}

		count := rs.Count

		for _, query := range rs.Results {
			if query.Name == name {
				d.SetId(strconv.Itoa(query.ID))
				return readQuery(ctx, d, meta)
			}
		}

		if count <= rs.PageSize*rs.Page {
			break
		}

		input.PageSize++
	}

	return diag.Errorf("Query (%s) not found", name)
}
