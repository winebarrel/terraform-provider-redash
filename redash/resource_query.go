package redash

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	redashgo "github.com/winebarrel/redash-go/v2"
)

func resourceQuery() *schema.Resource {
	return &schema.Resource{
		CreateContext: createQuery,
		ReadContext:   readQuery,
		UpdateContext: updateQuery,
		DeleteContext: deleteQuery,
		Importer: &schema.ResourceImporter{
			StateContext: importQuery,
		},
		Schema: map[string]*schema.Schema{
			"data_source_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"query": {
				Type:     schema.TypeString,
				Required: true,
			},
			"published": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"schedule": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"interval": {
							Description: "Interval in seconds.",
							Type:        schema.TypeInt,
							Required:    true,
							ValidateFunc: func(val any, key string) (warns []string, errs []error) {
								v := val.(int)

								if v < 1 {
									errs = append(errs, fmt.Errorf("%q must be >= 1, got: %d", key, v))
								}

								return
							},
						},
					},
				},
			},
			"options": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parameter": {
							Type:     schema.TypeSet,
							MinItems: 1,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"title": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"type": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"text",
											"text-pattern",
											"number",
											"enum",
											"query",
											"date",
											"datetime-local",
											"datetime-with-seconds",
											"date-range",
											"datetime-range",
											"datetime-range-with-seconds",
										}, false),
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"regex": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"enum": {
										Type:     schema.TypeList,
										MaxItems: 1,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"options": {
													Type:     schema.TypeList,
													MinItems: 1,
													Required: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"multi_values": {
													Type:     schema.TypeList,
													MaxItems: 1,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"quotation": {
																Type:         schema.TypeString,
																Optional:     true,
																ValidateFunc: validation.StringInSlice([]string{`"`, `'`}, false),
															},
															"separator": {
																Type:     schema.TypeString,
																Optional: true,
																Default:  ",",
															},
														},
													},
												},
											},
										},
									},
									"query_id": {
										Type:     schema.TypeInt,
										Optional: true,
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

func getQueryOptions(d *schema.ResourceData) []redashgo.QueryOptionsParameter {
	v, ok := d.GetOk("options")

	if !ok {
		return nil
	}

	options := v.([]any)

	if len(options) == 0 {
		return nil
	}

	optionsBlk := options[0].(map[string]any)
	parameterList, ok := optionsBlk["parameter"]

	if !ok {
		return nil
	}

	parameters := parameterList.(*schema.Set).List()

	if len(parameters) == 0 {
		return nil
	}

	queryParams := []redashgo.QueryOptionsParameter{}

	for _, p := range parameters {
		param := p.(map[string]any)
		paramName := param["name"].(string)

		qp := redashgo.QueryOptionsParameter{
			Title: paramName,
			Name:  paramName,
			Type:  param["type"].(string),
		}

		if title, ok := param["title"]; ok {
			qp.Title = title.(string)
		}

		if value, ok := param["value"]; ok {
			qp.Value = value.(string)
		}

		if regex, ok := param["regex"]; ok {
			qp.Regex = regex.(string)
		}

		if enumBlk, ok := param["enum"]; ok {
			enumList := enumBlk.([]any)

			if len(enumList) == 1 {
				enum := enumList[0].(map[string]any)
				enumOptions := []string{}

				for _, o := range enum["options"].([]any) {
					enumOptions = append(enumOptions, o.(string))
				}

				qp.EnumOptions = strings.Join(enumOptions, "\n")

				if multiValuesBlk, ok := enum["multi_values"]; ok {
					multiValuesList := multiValuesBlk.([]any)

					if len(multiValuesList) >= 1 {
						multiValues := multiValuesList[0].(map[string]any)
						qbmvo := &redashgo.QueryOptionsParameterMultiValuesOptions{}

						if quotation, ok := multiValues["quotation"]; ok {
							qs := quotation.(string)
							qbmvo.Prefix = qs
							qbmvo.Suffix = qs
						}

						if separator, ok := multiValues["separator"]; ok {
							qbmvo.Separator = separator.(string)
						}

						qp.MultiValuesOptions = qbmvo
					}
				}
			}
		}

		if queryID, ok := param["query_id"]; ok {
			if i := queryID.(int); i >= 1 {
				qp.QueryID = i
			}
		}

		queryParams = append(queryParams, qp)
	}

	return queryParams
}

func createQuery(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*redashgo.Client)

	input := &redashgo.CreateQueryInput{
		DataSourceID: d.Get("data_source_id").(int),
		Name:         d.Get("name").(string),
		Description:  d.Get("description").(string),
		Query:        d.Get("query").(string),
		Schedule:     &redashgo.CreateQueryInputSchedule{},
	}

	if v, ok := d.GetOk("schedule"); ok {
		schedules := v.([]any)

		if len(schedules) == 1 {
			schedule := schedules[0].(map[string]any)
			input.Schedule.Interval = schedule["interval"].(int)
		}
	}

	queryParams := getQueryOptions(d)

	if len(queryParams) >= 1 {
		input.Options = &redashgo.CreateQueryInputOptions{
			Parameters: queryParams,
		}
	}

	if v, ok := d.GetOk("tags"); ok {
		tags := []string{}

		for _, t := range v.([]any) {
			tags = append(tags, t.(string))
		}

		if len(tags) > 0 {
			input.Tags = tags
		}
	}

	query, err := client.CreateQuery(ctx, input)

	if err != nil {
		return diag.FromErr(err)
	}

	if v, ok := d.GetOk("published"); ok {
		published := v.(bool)

		if published {
			err = client.PublishQuery(ctx, query.ID)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	d.SetId(strconv.Itoa(query.ID))

	return readQuery(ctx, d, meta)
}

func readQuery(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	err := readQuery0(ctx, d, meta)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func readQuery0(ctx context.Context, d *schema.ResourceData, meta any) error {
	id, err := strconv.Atoi(d.Id())

	if err != nil {
		return err
	}

	client := meta.(*redashgo.Client)
	query, err := client.GetQuery(ctx, id)

	if err != nil {
		return err
	}

	d.Set("data_source_id", query.DataSourceID) //nolint:errcheck
	d.Set("name", query.Name)                   //nolint:errcheck
	d.Set("description", query.Description)     //nolint:errcheck
	d.Set("query", query.Query)                 //nolint:errcheck
	d.Set("published", !query.IsDraft)          //nolint:errcheck

	if query.Schedule != nil && query.Schedule.Interval != 0 {
		schedule := map[string]any{"interval": query.Schedule.Interval}
		d.Set("schedule", []any{schedule}) //nolint:errcheck
	}

	if len(query.Tags) > 0 {
		tags := []any{}

		for _, t := range query.Tags {
			tags = append(tags, t)
		}

		d.Set("tags", tags) //nolint:errcheck
	}

	if len(query.Options.Parameters) >= 1 {
		parameters := []any{}

		for _, param := range query.Options.Parameters {
			m := map[string]any{}
			m["title"] = param.Title
			m["name"] = param.Name
			m["type"] = param.Type

			if param.Value != nil {
				m["value"] = param.Value
			}

			if param.Regex != "" {
				m["regex"] = param.Regex
			}

			if param.QueryID > 0 {
				m["query_id"] = param.QueryID
			}

			if param.EnumOptions != "" {
				enum := map[string]any{}
				m["enum"] = []any{enum}
				enum["options"] = strings.Split(param.EnumOptions, "\n")

				if param.MultiValuesOptions != nil {
					mv := map[string]any{}
					enum["multi_values"] = []any{mv}
					// NOTE: Suffix is ​​the same value as Prefix
					mv["quotation"] = param.MultiValuesOptions.Prefix
					mv["separator"] = param.MultiValuesOptions.Separator
				}
			}

			parameters = append(parameters, m)
		}

		d.Set("options", []any{map[string]any{"parameter": parameters}}) //nolint:errcheck
	}

	return nil
}

func updateQuery(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, _ := strconv.Atoi(d.Id())
	client := meta.(*redashgo.Client)
	isDraft := !d.Get("published").(bool)

	input := &redashgo.UpdateQueryInput{
		DataSourceID: d.Get("data_source_id").(int),
		Name:         d.Get("name").(string),
		Description:  d.Get("description").(string),
		Query:        d.Get("query").(string),
		Schedule:     &redashgo.UpdateQueryInputSchedule{},
		IsDraft:      &isDraft,
	}

	schedules := d.Get("schedule").([]any)

	if len(schedules) == 1 {
		schedule := schedules[0].(map[string]any)
		input.Schedule.Interval = schedule["interval"].(int)
	}

	queryParams := getQueryOptions(d)

	if len(queryParams) >= 1 {
		input.Options = &redashgo.UpdateQueryInputOptions{
			Parameters: queryParams,
		}
	}

	if d.HasChange("tags") {
		tags := []string{}

		if v, ok := d.GetOk("tags"); ok {
			for _, t := range v.([]any) {
				tags = append(tags, t.(string))
			}
		}

		input.Tags = &tags
	}

	_, err := client.UpdateQuery(ctx, id, input)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func deleteQuery(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, _ := strconv.Atoi(d.Id())
	client := meta.(*redashgo.Client)

	err := client.ArchiveQuery(ctx, id)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func importQuery(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	err := readQuery0(ctx, d, meta)

	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
