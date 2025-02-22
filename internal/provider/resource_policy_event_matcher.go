package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourcePolicyEventMatcher() *schema.Resource {
	return &schema.Resource{
		Description:   "Customization --- ",
		CreateContext: resourcePolicyEventMatcherCreate,
		ReadContext:   resourcePolicyEventMatcherRead,
		UpdateContext: resourcePolicyEventMatcherUpdate,
		DeleteContext: resourcePolicyEventMatcherDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"execution_logging": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"app": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      EnumToDescription(api.AllowedAppEnumEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedAppEnumEnumValues),
			},
			"model": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      EnumToDescription(api.AllowedModelEnumEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedModelEnumEnumValues),
			},
		},
	}
}

func resourcePolicyEventMatcherSchemaToProvider(d *schema.ResourceData) *api.EventMatcherPolicyRequest {
	r := api.EventMatcherPolicyRequest{
		Name:             d.Get("name").(string),
		ExecutionLogging: api.PtrBool(d.Get("execution_logging").(bool)),
	}

	if a, ok := d.Get("action").(string); ok && a != "" {
		r.Action.Set(api.EventActions(a).Ptr())
	} else {
		r.Action.Set(nil)
	}
	if p, ok := d.Get("client_ip").(string); ok && p != "" {
		r.ClientIp.Set(api.PtrString(p))
	} else {
		r.ClientIp.Set(nil)
	}
	if a, ok := d.Get("app").(string); ok && a != "" {
		r.App.Set(api.AppEnum(a).Ptr())
	} else {
		r.App.Set(nil)
	}
	if m, ok := d.Get("model").(string); ok && m != "" {
		r.Model.Set(api.ModelEnum(m).Ptr())
	} else {
		r.Model.Set(nil)
	}
	return &r
}

func resourcePolicyEventMatcherCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourcePolicyEventMatcherSchemaToProvider(d)

	res, hr, err := c.client.PoliciesApi.PoliciesEventMatcherCreate(ctx).EventMatcherPolicyRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePolicyEventMatcherRead(ctx, d, m)
}

func resourcePolicyEventMatcherRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PoliciesApi.PoliciesEventMatcherRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "execution_logging", res.ExecutionLogging)
	if res.HasAction() {
		setWrapper(d, "action", res.Action.Get())
	}
	if res.HasClientIp() {
		setWrapper(d, "client_ip", res.ClientIp.Get())
	}
	if res.HasApp() {
		setWrapper(d, "app", res.App.Get())
	}
	if res.HasModel() {
		setWrapper(d, "model", res.Model.Get())
	}
	return diags
}

func resourcePolicyEventMatcherUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourcePolicyEventMatcherSchemaToProvider(d)

	res, hr, err := c.client.PoliciesApi.PoliciesEventMatcherUpdate(ctx, d.Id()).EventMatcherPolicyRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePolicyEventMatcherRead(ctx, d, m)
}

func resourcePolicyEventMatcherDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PoliciesApi.PoliciesEventMatcherDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
