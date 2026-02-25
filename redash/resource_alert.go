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

var (
	// cf. https://github.com/getredash/redash/blob/v10.1.0/redash/models/__init__.py#L923-L934
	alertOperators = []string{
		">",
		">=",
		"<",
		"<=",
		"==",
		"!=",
		"greater than",
		"less than",
		"equals",
	}
)

func resourceAlert() *schema.Resource {
	return &schema.Resource{
		CreateContext: createAlert,
		ReadContext:   readAlert,
		UpdateContext: updateAlert,
		DeleteContext: deleteAlert,
		Importer: &schema.ResourceImporter{
			StateContext: importAlert,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"query_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"options": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"column": {
							Type:     schema.TypeString,
							Required: true,
						},
						"op": {
							Type:     schema.TypeString,
							Required: true,
							// cf. https://github.com/getredash/redash/blob/v10.1.0/redash/models/__init__.py#L923
							ValidateFunc: func(val any, key string) (warns []string, errs []error) {
								v := val.(string)

								for _, op := range alertOperators {
									if op == v {
										return
									}
								}

								errs = append(errs, fmt.Errorf("must be a valid operator (%s), got: %s", strings.Join(alertOperators, ","), v))

								return
							},
						},
						"value": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"custom_subject": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"custom_body": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"template": {
							Type:       schema.TypeString,
							Optional:   true,
							Default:    "",
							Deprecated: "This attribute is for backward compatibility.",
						},
					},
				},
			},
			"rearm": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
				ValidateFunc: func(val any, key string) (warns []string, errs []error) {
					v := val.(int)

					if v < 0 {
						errs = append(errs, fmt.Errorf("%q must be >= 0, got: %d", key, v))
					}

					return
				},
			},
			"muted": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func createAlert(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*redashgo.Client)

	optionsList := d.Get("options").([]any)
	options := optionsList[0].(map[string]any)

	input := &redashgo.CreateAlertInput{
		Name:    d.Get("name").(string),
		QueryId: d.Get("query_id").(int),
		Rearm:   d.Get("rearm").(int),
		Options: redashgo.CreateAlertOptions{
			Column:        options["column"].(string),
			Op:            options["op"].(string),
			Value:         float64(options["value"].(int)),
			CustomSubject: options["custom_subject"].(string),
			CustomBody:    options["custom_body"].(string),
			Template:      options["template"].(string),
		},
	}

	alert, err := client.CreateAlert(ctx, input)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(alert.ID))

	if d.Get("muted").(bool) {
		err = client.MuteAlert(ctx, alert.ID)

		if err != nil {
			return diag.FromErr(err)
		}
	}

	return readAlert(ctx, d, meta)
}

func readAlert(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	err := readAlert0(ctx, d, meta)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func readAlert0(ctx context.Context, d *schema.ResourceData, meta any) error {
	id, err := strconv.Atoi(d.Id())

	if err != nil {
		return err
	}

	client := meta.(*redashgo.Client)
	alert, err := client.GetAlert(ctx, id)

	if err != nil {
		return err
	}

	d.Set("name", alert.Name)           //nolint:errcheck
	d.Set("query_id", alert.Query.ID)   //nolint:errcheck
	d.Set("rearm", alert.Rearm)         //nolint:errcheck
	d.Set("muted", alert.Options.Muted) //nolint:errcheck

	options := map[string]any{
		"column":         alert.Options.Column,
		"op":             alert.Options.Op,
		"value":          alert.Options.Value,
		"custom_subject": alert.Options.CustomSubject,
		"custom_body":    alert.Options.CustomBody,
		"template":       alert.Options.Template, //nolint:staticcheck
	}

	d.Set("options", []any{options}) //nolint:errcheck

	return nil
}

func updateAlert(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, _ := strconv.Atoi(d.Id())
	client := meta.(*redashgo.Client)

	input := &redashgo.UpdateAlertInput{
		Name:    d.Get("name").(string),
		QueryId: d.Get("query_id").(int),
		Rearm:   d.Get("rearm").(int),
	}

	optionsList := d.Get("options").([]any)
	options := optionsList[0].(map[string]any)

	input.Options = &redashgo.UpdateAlertOptions{
		Column:        options["column"].(string),
		Op:            options["op"].(string),
		Value:         float64(options["value"].(int)),
		CustomSubject: options["custom_subject"].(string),
		CustomBody:    options["custom_body"].(string),
		Template:      options["template"].(string),
	}

	_, err := client.UpdateAlert(ctx, id, input)

	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("muted") {
		if d.Get("muted").(bool) {
			err = client.MuteAlert(ctx, id)
		} else {
			err = client.UnmuteAlert(ctx, id)
		}

		if err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func deleteAlert(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, _ := strconv.Atoi(d.Id())
	client := meta.(*redashgo.Client)

	err := client.DeleteAlert(ctx, id)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func importAlert(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	err := readAlert0(ctx, d, meta)

	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
