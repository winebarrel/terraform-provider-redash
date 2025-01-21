package redash

import (
	"context"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redashgo "github.com/winebarrel/redash-go/v2"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: readUserByName,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"name", "email"},
			},
			"email": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"name", "email"},
			},
		},
	}
}

func readUserByName(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*redashgo.Client)
	var name, email string
	msgId := []string{}

	if _, ok := d.GetOk("name"); ok {
		name = d.Get("name").(string)
		msgId = append(msgId, name)
	}

	if _, ok := d.GetOk("email"); ok {
		email = d.Get("email").(string)
		msgId = append(msgId, email)
	}

	input := &redashgo.ListUsersInput{
		Page:     1,
		PageSize: 20,
	}

	for {
		rs, err := client.ListUsers(ctx, input)

		if err != nil {
			return diag.FromErr(err)
		}

		count := rs.Count

		for _, user := range rs.Results {
			if user.Name == name || user.Email == email {
				d.SetId(strconv.Itoa(user.ID))
				d.Set("name", user.Name)   //nolint:errcheck
				d.Set("email", user.Email) //nolint:errcheck
				return nil
			}
		}

		if count <= rs.PageSize*rs.Page {
			break
		}

		input.Page++
	}

	return diag.Errorf("User (%s) not found", strings.Join(msgId, "/"))
}
