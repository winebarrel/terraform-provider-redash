package redash

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			client.PublishQuery(ctx, query.ID)
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
