package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourceEventTransport() *schema.Resource {
	return &schema.Resource{
		Description:   "Events --- ",
		CreateContext: resourceEventTransportCreate,
		ReadContext:   resourceEventTransportRead,
		UpdateContext: resourceEventTransportUpdate,
		DeleteContext: resourceEventTransportDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mode": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      EnumToDescription(api.AllowedNotificationTransportModeEnumEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedNotificationTransportModeEnumEnumValues),
			},
			"webhook_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"webhook_mapping": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"send_once": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceEventTransportSchemaToModel(d *schema.ResourceData) (*api.NotificationTransportRequest, diag.Diagnostics) {
	m := api.NotificationTransportRequest{
		Name:     d.Get("name").(string),
		SendOnce: api.PtrBool(d.Get("send_once").(bool)),
		Mode:     api.NotificationTransportModeEnum(d.Get("mode").(string)).Ptr(),
	}

	if w, ok := d.Get("webhook_url").(string); ok {
		m.WebhookUrl = &w
	}

	if w, ok := d.Get("webhook_mapping").(string); ok {
		m.WebhookMapping.Set(&w)
	} else {
		m.WebhookMapping.Set(nil)
	}
	return &m, nil
}

func resourceEventTransportCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app, diags := resourceEventTransportSchemaToModel(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.EventsApi.EventsTransportsCreate(ctx).NotificationTransportRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceEventTransportRead(ctx, d, m)
}

func resourceEventTransportRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.EventsApi.EventsTransportsRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "mode", res.Mode)
	setWrapper(d, "send_once", res.SendOnce)
	setWrapper(d, "webhook_url", res.WebhookUrl)
	setWrapper(d, "webhook_mapping", res.WebhookMapping.Get())
	return diags
}

func resourceEventTransportUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app, di := resourceEventTransportSchemaToModel(d)
	if di != nil {
		return di
	}
	res, hr, err := c.client.EventsApi.EventsTransportsUpdate(ctx, d.Id()).NotificationTransportRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceEventTransportRead(ctx, d, m)
}

func resourceEventTransportDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.EventsApi.EventsTransportsDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
