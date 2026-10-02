package redash

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redashgo "github.com/winebarrel/redash-go/v2"
)

func resourceGroupDataSources() *schema.Resource {
	return &schema.Resource{
		Description: "Authoritative for a given group. Updates the data sources granted to that group to match this list. Data sources not listed here are removed.\n\n" +
			"!> **Warning:** Do not use this resource together with `redash_group_data_source` for the same group. Both manage the same grants and will conflict.",
		CreateContext: createGroupDataSources,
		ReadContext:   readGroupDataSourcesSet,
		UpdateContext: updateGroupDataSources,
		DeleteContext: deleteGroupDataSources,
		Importer: &schema.ResourceImporter{
			StateContext: importGroupDataSources,
		},
		CustomizeDiff: rejectDuplicateGroupDataSources,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Description: "ID of the group.",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"data_source": {
				Description: "Data sources granted to the group. Omit this to grant none.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"data_source_id": {
							Description: "ID of the data source.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"name": {
							Description: "Name of the data source. Set this to show the name in the plan. After apply, the name returned by Redash is stored.",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
						"view_only": {
							Description: "When true, the group has view-only access. Defaults to false (full access).",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
					},
				},
			},
		},
	}
}

func createGroupDataSources(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	groupId := d.Get("group_id").(int)
	if diags := syncGroupDataSources(ctx, d, meta, groupId); diags != nil {
		return diags
	}

	d.SetId(strconv.Itoa(groupId))

	return readGroupDataSourcesSet(ctx, d, meta)
}

func readGroupDataSourcesSet(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	groupId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := meta.(*redashgo.Client)
	dsList, err := client.ListGroupDataSources(ctx, groupId)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSources := make([]map[string]any, 0, len(dsList))
	for _, ds := range dsList {
		dataSources = append(dataSources, map[string]any{
			"data_source_id": ds.ID,
			"name":           ds.Name,
			"view_only":      ds.ViewOnly,
		})
	}

	d.Set("group_id", groupId)        //nolint:errcheck
	d.Set("data_source", dataSources) //nolint:errcheck

	return nil
}

func updateGroupDataSources(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	groupId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if diags := syncGroupDataSources(ctx, d, meta, groupId); diags != nil {
		return diags
	}

	return readGroupDataSourcesSet(ctx, d, meta)
}

func deleteGroupDataSources(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	groupId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := meta.(*redashgo.Client)
	dsList, err := client.ListGroupDataSources(ctx, groupId)
	if err != nil {
		return diag.FromErr(err)
	}

	seen := make(map[int]struct{}, len(dsList))
	for _, ds := range dsList {
		if _, ok := seen[ds.ID]; ok {
			continue
		}
		seen[ds.ID] = struct{}{}

		if err := client.RemoveGroupDataSource(ctx, groupId, ds.ID); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId("")

	return nil
}

func importGroupDataSources(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	groupId, err := strconv.Atoi(d.Id())
	if err != nil {
		return nil, fmt.Errorf("invalid import ID: %q", d.Id())
	}

	d.Set("group_id", groupId) //nolint:errcheck

	return []*schema.ResourceData{d}, nil
}

func syncGroupDataSources(ctx context.Context, d *schema.ResourceData, meta any, groupId int) diag.Diagnostics {
	desired, err := configuredGroupDataSources(d)
	if err != nil {
		return diag.FromErr(err)
	}

	client := meta.(*redashgo.Client)
	dsList, err := client.ListGroupDataSources(ctx, groupId)
	if err != nil {
		return diag.FromErr(err)
	}

	current := make(map[int]bool, len(dsList))
	for _, ds := range dsList {
		current[ds.ID] = ds.ViewOnly
	}

	for id, viewOnly := range desired {
		have, ok := current[id]
		if !ok {
			if _, err := client.AddGroupDataSource(ctx, groupId, id); err != nil {
				return diag.FromErr(err)
			}
			have = false
		}

		if have != viewOnly {
			_, err := client.UpdateGroupDataSource(ctx, groupId, id, &redashgo.UpdateGroupDataSourceInput{
				ViewOnly: viewOnly,
			})
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	for id := range current {
		if _, ok := desired[id]; ok {
			continue
		}

		if err := client.RemoveGroupDataSource(ctx, groupId, id); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func configuredGroupDataSources(d interface{ Get(string) any }) (map[int]bool, error) {
	set := d.Get("data_source").(*schema.Set)
	grants := make(map[int]bool, set.Len())

	for _, raw := range set.List() {
		item := raw.(map[string]any)
		id := item["data_source_id"].(int)
		if _, ok := grants[id]; ok {
			return nil, fmt.Errorf("duplicate data_source_id %d", id)
		}
		grants[id] = item["view_only"].(bool)
	}

	return grants, nil
}

func rejectDuplicateGroupDataSources(_ context.Context, d *schema.ResourceDiff, _ any) error {
	_, err := configuredGroupDataSources(d)
	return err
}
