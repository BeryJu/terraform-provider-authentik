package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourceEventRule() *schema.Resource {
	return &schema.Resource{
		Description:   "Events --- ",
		CreateContext: resourceEventRuleCreate,
		ReadContext:   resourceEventRuleRead,
		UpdateContext: resourceEventRuleUpdate,
		DeleteContext: resourceEventRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"transports": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"severity": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.SEVERITYENUM_WARNING,
				Description:      EnumToDescription(api.AllowedSeverityEnumEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedSeverityEnumEnumValues),
			},
			"webhook_mapping": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceEventRuleSchemaToModel(d *schema.ResourceData) (*api.NotificationRuleRequest, diag.Diagnostics) {
	m := api.NotificationRuleRequest{
		Name:     d.Get("name").(string),
		Severity: api.SeverityEnum(d.Get("severity").(string)).Ptr(),
	}

	if w, ok := d.Get("group").(string); ok {
		m.Group.Set(&w)
	} else {
		m.Group.Set(nil)
	}

	m.Transports = castSlice[string](d.Get("transports").([]interface{}))
	return &m, nil
}

func resourceEventRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app, diags := resourceEventRuleSchemaToModel(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.EventsApi.EventsRulesCreate(ctx).NotificationRuleRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceEventRuleRead(ctx, d, m)
}

func resourceEventRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.EventsApi.EventsRulesRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "group", res.Group.Get())
	localTransports := castSlice[string](d.Get("transports").([]interface{}))
	setWrapper(d, "transports", listConsistentMerge(localTransports, res.Transports))
	setWrapper(d, "severity", res.Severity)
	return diags
}

func resourceEventRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app, di := resourceEventRuleSchemaToModel(d)
	if di != nil {
		return di
	}
	res, hr, err := c.client.EventsApi.EventsRulesUpdate(ctx, d.Id()).NotificationRuleRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceEventRuleRead(ctx, d, m)
}

func resourceEventRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.EventsApi.EventsRulesDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
