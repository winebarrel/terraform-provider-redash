package redash

import (
	"context"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redashgo "github.com/winebarrel/redash-go"
)

func dataSourceUsers() *schema.Resource {
	return &schema.Resource{
		Description: "Data Source to get user IDs.",
		ReadContext: readUsersByFilter,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Name filter. (e.g. `J* Smith`)",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"email": {
				Description: "Email filter. (e.g. `*@example.com`)",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"ids": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Computed: true,
			},
		},
	}
}

func globToRegexp(s string) string {
	s = regexp.QuoteMeta(s)
	s = strings.ReplaceAll(s, `\*`, ".*")
	s = strings.ReplaceAll(s, `\?`, ".")
	return "^" + s + "$"
}

func readUsersByFilter(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*redashgo.Client)
	rName := regexp.MustCompile("")
	rEmail := regexp.MustCompile("")

	if _, ok := d.GetOk("name"); ok {
		name := d.Get("name").(string)
		rName = regexp.MustCompile(globToRegexp(name))
	}

	if _, ok := d.GetOk("email"); ok {
		email := d.Get("email").(string)
		rEmail = regexp.MustCompile(globToRegexp(email))
	}

	input := &redashgo.ListUsersInput{
		Page:     1,
		PageSize: 20,
	}

	ids := []int{}

	for {
		rs, err := client.ListUsers(ctx, input)

		if err != nil {
			return diag.FromErr(err)
		}

		count := rs.Count

		for _, user := range rs.Results {
			if rName.MatchString(user.Name) && rEmail.MatchString(user.Email) {
				ids = append(ids, user.ID)
			}
		}

		if count <= rs.PageSize*rs.Page {
			break
		}

		input.Page++
	}

	d.SetId(id.UniqueId())
	d.Set("ids", ids) //nolint:errcheck

	return nil
}
