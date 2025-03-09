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
			"query_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				AtLeastOneOf: []string{"query_id", "name"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"query_id", "name"},
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
			"options": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parameter": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"title": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"regex": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"enum": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"options": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"multi_values": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"quotation": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"separator": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
											},
										},
									},
									"query_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
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
	var queryId *int

	if rawId, ok := d.GetOk("query_id"); ok {
		n := rawId.(int)
		queryId = &n
	}

	name := d.Get("name").(string)

	if queryId != nil {
		query, err := client.GetQuery(ctx, *queryId)

		if err != nil {
			return diag.Errorf("Query not found: %s", err)
		}

		if name != "" && query.Name != name {
			return diag.Errorf("Query (%s) not found", name)
		}

		d.SetId(strconv.Itoa(query.ID))
		return readQuery(ctx, d, meta)
	}

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
