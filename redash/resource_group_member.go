package redash

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redashgo "github.com/winebarrel/redash-go/v2"
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
		Description:   "Adds one user to one group. When that user is already a member, create records the membership and does not add another copy. Destroy removes every stored copy of the group ID, because Redash deletes only one copy per request.",
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

	members, err := client.ListGroupMembers(ctx, groupId)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, member := range members {
		if member.ID == userId {
			d.SetId(fmt.Sprintf("%d/%d", groupId, userId))
			return nil
		}
	}

	if _, err := client.AddGroupMember(ctx, groupId, userId); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d/%d", groupId, userId))

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
	// Redash does not return 404 for a missing group. This endpoint filters users by the stored group id and does not load the group. Deleting a group removes one stored copy, so a user who still has that id is still returned.
	members, err := client.ListGroupMembers(ctx, groupId)
	if err != nil {
		return diag.FromErr(err)
	}

	found := false
	for _, member := range members {
		if member.ID == userId {
			found = true
			break
		}
	}

	if !found {
		d.SetId("")
		return nil
	}

	d.Set("group_id", groupId) //nolint:errcheck
	d.Set("user_id", userId)   //nolint:errcheck

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

	if err := removeGroupMember(ctx, client, memberId, groupId); err != nil {
		return diag.FromErr(err)
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

// Redash removes one stored copy of the group id per delete.
func removeGroupMember(ctx context.Context, client *redashgo.Client, userID, groupID int) error {
	for {
		n, err := groupCopies(ctx, client, userID, groupID)
		if err != nil {
			return err
		}

		if n == 0 {
			return nil
		}

		if err := client.RemoveGroupMember(ctx, groupID, userID); err != nil {
			return err
		}

		left, err := groupCopies(ctx, client, userID, groupID)
		if err != nil {
			return err
		}

		if left >= n {
			return fmt.Errorf("user %d: group %d still has %d %s", userID, groupID, left, copyNoun(left))
		}
	}
}

func copyNoun(n int) string {
	if n == 1 {
		return "copy"
	}

	return "copies"
}

func groupCopies(ctx context.Context, client *redashgo.Client, userID, groupID int) (int, error) {
	user, err := client.GetUser(ctx, userID)
	if err != nil {
		return 0, err
	}

	n := 0
	for _, group := range user.Groups {
		id, ok := groupIDFromJSON(group)
		if !ok {
			return 0, fmt.Errorf("user %d: group id has type %T", userID, group)
		}

		if id == groupID {
			n++
		}
	}

	return n, nil
}

func groupIDFromJSON(v any) (int, bool) {
	switch n := v.(type) {
	case float64:
		return int(n), true
	case int:
		return n, true
	default:
		return 0, false
	}
}
