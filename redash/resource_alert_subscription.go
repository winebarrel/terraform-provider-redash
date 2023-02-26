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
		Description:   "Alert Subscription resource.",
		CreateContext: createAlertSubscription,
		ReadContext:   schema.NoopContext,
		DeleteContext: deleteAlertSubscription,
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
