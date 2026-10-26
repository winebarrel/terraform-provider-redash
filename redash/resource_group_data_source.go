package redash

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redashgo "github.com/winebarrel/redash-go/v2"
)

func resourceGroupDataSource() *schema.Resource {
	return &schema.Resource{
		Description:   "!> **Warning:** Do not use this resource together with `redash_group_data_sources` for the same group. Both manage the same grants and will conflict.",
		CreateContext: createGroupDataSource,
		ReadContext:   readGroupDataSource,
		UpdateContext: updateGroupDataSource,
		DeleteContext: deleteGroupDataSource,
		Importer: &schema.ResourceImporter{
			StateContext: importGroupDataSource,
		},
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"data_source_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"view_only": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func createGroupDataSource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*providerMeta).client
	groupId := d.Get("group_id").(int)
	dsId := d.Get("data_source_id").(int)
	gds, err := client.AddGroupDataSource(ctx, groupId, dsId)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.Get("view_only").(bool) {
		_, err = client.UpdateGroupDataSource(ctx, groupId, gds.ID, &redashgo.UpdateGroupDataSourceInput{
			ViewOnly: true,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(fmt.Sprintf("%d/%d", groupId, gds.ID))

	return nil
}

func readGroupDataSource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	groupIdStr, gdsIdStr, _ := strings.Cut(d.Id(), "/")
	groupId, _ := strconv.Atoi(groupIdStr)
	gdsId, _ := strconv.Atoi(gdsIdStr)
	client := meta.(*providerMeta).client

	dsList, err := client.ListGroupDataSources(ctx, groupId)
	if err != nil {
		// The group was deleted outside Terraform.
		if strings.Contains(err.Error(), "HTTP status code not OK: 404") {
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	for _, ds := range dsList {
		if ds.ID == gdsId {
			d.Set("group_id", groupId)      //nolint:errcheck
			d.Set("data_source_id", ds.ID)  //nolint:errcheck
			d.Set("view_only", ds.ViewOnly) //nolint:errcheck
			return nil
		}
	}

	d.SetId("")

	return nil
}

func updateGroupDataSource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	groupIdStr, gdsIdStr, _ := strings.Cut(d.Id(), "/")
	groupId, _ := strconv.Atoi(groupIdStr)
	gdsId, _ := strconv.Atoi(gdsIdStr)
	client := meta.(*providerMeta).client

	if d.HasChange("view_only") {
		_, err := client.UpdateGroupDataSource(ctx, groupId, gdsId, &redashgo.UpdateGroupDataSourceInput{
			ViewOnly: d.Get("view_only").(bool),
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func deleteGroupDataSource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	groupIdStr, gdsIdStr, _ := strings.Cut(d.Id(), "/")
	groupId, _ := strconv.Atoi(groupIdStr)
	gdsId, _ := strconv.Atoi(gdsIdStr)
	client := meta.(*providerMeta).client

	if err := removeGroupDataSource(ctx, client, groupId, gdsId); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

// removeGroupDataSource revokes a grant, and treats one that is already gone
// as revoked. Redash drops the grants of a deleted data source, so when the
// data source is destroyed in the same apply, the removal can fail after the
// grant has disappeared.
func removeGroupDataSource(ctx context.Context, client *redashgo.Client, groupId int, dsId int) error {
	err := client.RemoveGroupDataSource(ctx, groupId, dsId)
	if err == nil {
		return nil
	}

	dsList, listErr := client.ListGroupDataSources(ctx, groupId)
	if listErr != nil {
		return err
	}

	for _, ds := range dsList {
		if ds.ID == dsId {
			return err
		}
	}

	return nil
}

func importGroupDataSource(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	groupIdStr, gdsIdStr, ok := strings.Cut(d.Id(), "/")
	if !ok {
		return nil, fmt.Errorf("invalid import ID: %q", d.Id())
	}

	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		return nil, err
	}

	gdsId, err := strconv.Atoi(gdsIdStr)
	if err != nil {
		return nil, err
	}

	d.Set("group_id", groupId)     //nolint:errcheck
	d.Set("data_source_id", gdsId) //nolint:errcheck

	return []*schema.ResourceData{d}, nil
}
