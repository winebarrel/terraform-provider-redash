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
	groupGdsId := strings.SplitN(d.Id(), "/", 2)
	groupId, _ := strconv.Atoi(groupGdsId[0])
	gdsId, _ := strconv.Atoi(groupGdsId[1])
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
	groupGdsId := strings.SplitN(d.Id(), "/", 2)
	groupId, _ := strconv.Atoi(groupGdsId[0])
	gdsId, _ := strconv.Atoi(groupGdsId[1])
	client := meta.(*redashgo.Client)

	err := client.RemoveGroupDataSource(ctx, groupId, gdsId)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func importGroupDataSource(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	groupGdsId := strings.SplitN(d.Id(), "/", 2)
	groupId, err := strconv.Atoi(groupGdsId[0])

	if err != nil {
		return nil, err
	}

	gdsId, err := strconv.Atoi(groupGdsId[1])

	if err != nil {
		return nil, err
	}

	d.Set("group_id", groupId)     //nolint:errcheck
	d.Set("data_source_id", gdsId) //nolint:errcheck

	return []*schema.ResourceData{d}, nil
}
