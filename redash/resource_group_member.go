package redash

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Redash stores membership as users.groups, an array rewritten in full by
// POST/DELETE /api/groups/<id>/members. Parallel calls for one user lose
// updates. Serialize membership changes per user inside this process.
var groupMemberUserLocks sync.Map

func lockGroupMemberUser(userID int) func() {
	mu, _ := groupMemberUserLocks.LoadOrStore(userID, &sync.Mutex{})
	m := mu.(*sync.Mutex)
	m.Lock()

	return m.Unlock
}

func resourceGroupMember() *schema.Resource {
	return &schema.Resource{
		CreateContext: createGroupUser,
		ReadContext:   readGroupMember,
		DeleteContext: deleteGroupUser,
		Importer: &schema.ResourceImporter{
			StateContext: importGroupMember,
		},
		Schema: map[string]*schema.Schema{
			"group_id": {
				Description: "ID of the group.",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"user_id": {
				Description: "ID of the user to add to the group.",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func createGroupUser(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*providerMeta).client
	groupId := d.Get("group_id").(int)
	userId := d.Get("user_id").(int)

	unlock := lockGroupMemberUser(userId)
	defer unlock()

	member, err := client.AddGroupMember(ctx, groupId, userId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d/%d", groupId, member.ID))

	return nil
}

func readGroupMember(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	groupIdStr, userIdStr, ok := strings.Cut(d.Id(), "/")
	if !ok {
		return diag.Errorf("invalid ID: %q", d.Id())
	}

	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		return diag.FromErr(err)
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		return diag.FromErr(err)
	}

	client := meta.(*providerMeta).client
	// Redash does not return 404 for a deleted group. It returns an empty list,
	// because it only filters users by group ID without looking up the group.
	members, err := client.ListGroupMembers(ctx, groupId)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, m := range members {
		if m.ID == userId {
			d.Set("group_id", groupId) //nolint:errcheck
			d.Set("user_id", userId)   //nolint:errcheck
			return nil
		}
	}

	d.SetId("")

	return nil
}

func deleteGroupUser(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	groupIdStr, memberIdStr, ok := strings.Cut(d.Id(), "/")
	if !ok {
		return diag.Errorf("invalid ID: %q", d.Id())
	}

	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		return diag.FromErr(err)
	}

	memberId, err := strconv.Atoi(memberIdStr)
	if err != nil {
		return diag.FromErr(err)
	}

	client := meta.(*providerMeta).client

	unlock := lockGroupMemberUser(memberId)
	defer unlock()

	err = client.RemoveGroupMember(ctx, groupId, memberId)
	if err != nil {
		// Redash does user.group_ids.remove(group_id). A missing id raises
		// ValueError and the API returns 500. If the user is not in the
		// group, the membership is already gone.
		members, listErr := client.ListGroupMembers(ctx, groupId)
		if listErr != nil {
			return diag.FromErr(err)
		}

		for _, member := range members {
			if member.ID == memberId {
				return diag.FromErr(err)
			}
		}
	}

	d.SetId("")

	return nil
}

func importGroupMember(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	groupIdStr, memberIdStr, ok := strings.Cut(d.Id(), "/")
	if !ok {
		return nil, fmt.Errorf("invalid import ID: %q", d.Id())
	}

	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		return nil, err
	}

	memberId, err := strconv.Atoi(memberIdStr)
	if err != nil {
		return nil, err
	}

	d.Set("group_id", groupId) //nolint:errcheck
	d.Set("user_id", memberId) //nolint:errcheck

	return []*schema.ResourceData{d}, nil
}
