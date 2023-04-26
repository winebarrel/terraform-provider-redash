package redash

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redashgo "github.com/winebarrel/redash-go"
)

func resourceGroupMember() *schema.Resource {
	return &schema.Resource{
		CreateContext: createGroupUser,
		ReadContext:   schema.NoopContext,
		DeleteContext: deleteGroupUser,
		Importer: &schema.ResourceImporter{
			StateContext: importGroupMember,
		},
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"user_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func createGroupUser(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*redashgo.Client)
	groupId := d.Get("group_id").(int)
	userId := d.Get("user_id").(int)
	member, err := client.AddGroupMember(ctx, groupId, userId)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d/%d", groupId, member.ID))

	return nil
}

func deleteGroupUser(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	groupMemberId := strings.SplitN(d.Id(), "/", 2)
	groupId, _ := strconv.Atoi(groupMemberId[0])
	memberId, _ := strconv.Atoi(groupMemberId[1])
	client := meta.(*redashgo.Client)

	err := client.RemoveGroupMember(ctx, groupId, memberId)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func importGroupMember(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	groupMemberId := strings.SplitN(d.Id(), "/", 2)
	groupId, err := strconv.Atoi(groupMemberId[0])

	if err != nil {
		return nil, err
	}

	memberId, err := strconv.Atoi(groupMemberId[1])

	if err != nil {
		return nil, err
	}

	d.Set("group_id", groupId) //nolint:errcheck
	d.Set("user_id", memberId) //nolint:errcheck

	return []*schema.ResourceData{d}, nil
}
