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
		CreateContext: createGroupDataSource,
		ReadContext:   schema.NoopContext,
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
	client := meta.(*redashgo.Client)
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

func updateGroupDataSource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	groupIdStr, gdsIdStr, _ := strings.Cut(d.Id(), "/")
	groupId, _ := strconv.Atoi(groupIdStr)
	gdsId, _ := strconv.Atoi(gdsIdStr)
	client := meta.(*redashgo.Client)

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
	client := meta.(*redashgo.Client)

	err := client.RemoveGroupDataSource(ctx, groupId, gdsId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

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
