package redash

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redashgo "github.com/winebarrel/redash-go/v2"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: createUser,
		ReadContext:   readUser,
		UpdateContext: updateUser,
		DeleteContext: deleteUser,
		Importer: &schema.ResourceImporter{
			StateContext: importUser,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func createUser(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*redashgo.Client)

	input := &redashgo.CreateUsersInput{
		Name:  d.Get("name").(string),
		Email: d.Get("email").(string),
	}

	user, err := client.CreateUser(ctx, input)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(user.ID))

	return readUser(ctx, d, meta)
}

func readUser(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	err := readUser0(ctx, d, meta)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func readUser0(ctx context.Context, d *schema.ResourceData, meta any) error {
	id, err := strconv.Atoi(d.Id())

	if err != nil {
		return err
	}

	client := meta.(*redashgo.Client)
	user, err := client.GetUser(ctx, id)

	if err != nil {
		return err
	}

	d.Set("name", user.Name)   //nolint:errcheck
	d.Set("email", user.Email) //nolint:errcheck

	return nil
}

func updateUser(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, _ := strconv.Atoi(d.Id())
	client := meta.(*redashgo.Client)

	input := &redashgo.UpdateUserInput{
		Name:  d.Get("name").(string),
		Email: d.Get("email").(string),
	}

	_, err := client.UpdateUser(ctx, id, input)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func deleteUser(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, _ := strconv.Atoi(d.Id())
	client := meta.(*redashgo.Client)

	err := client.DeleteUser(ctx, id)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func importUser(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	err := readUser0(ctx, d, meta)

	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
