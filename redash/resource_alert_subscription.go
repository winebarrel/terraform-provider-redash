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

func resourceAlertSubscription() *schema.Resource {
	return &schema.Resource{
		CreateContext: createAlertSubscription,
		ReadContext:   schema.NoopContext,
		DeleteContext: deleteAlertSubscription,
		Importer: &schema.ResourceImporter{
			StateContext: importAlertSubscription,
		},

		Schema: map[string]*schema.Schema{
			"alert_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"alert_destination_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func createAlertSubscription(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*redashgo.Client)
	alertId := d.Get("alert_id").(int)
	destId := d.Get("alert_destination_id").(int)
	subs, err := client.AddAlertSubscription(ctx, alertId, destId)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d/%d", alertId, subs.ID))

	return nil
}

func deleteAlertSubscription(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	alertSubsId := strings.SplitN(d.Id(), "/", 2)
	alertId, _ := strconv.Atoi(alertSubsId[0])
	subsId, _ := strconv.Atoi(alertSubsId[1])
	client := meta.(*redashgo.Client)

	err := client.RemoveAlertSubscription(ctx, alertId, subsId)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func importAlertSubscription(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	alertDestId := strings.SplitN(d.Id(), "/", 2)
	alertId, _ := strconv.Atoi(alertDestId[0])
	destId, _ := strconv.Atoi(alertDestId[1])
	client := meta.(*redashgo.Client)

	subsList, err := client.ListAlertSubscriptions(ctx, alertId)

	if err != nil {
		return nil, err
	}

	for _, s := range subsList {
		if s.Destination.ID == destId {
			d.SetId(fmt.Sprintf("%d/%d", alertId, s.ID))
			d.Set("alert_id", s.AlertID)                    //nolint:errcheck
			d.Set("alert_destination_id", s.Destination.ID) //nolint:errcheck
			return []*schema.ResourceData{d}, nil
		}
	}

	return nil, fmt.Errorf("Alert Subscription (alert_id=%d, alert_destination_id=%d) not found", alertId, destId)
}
