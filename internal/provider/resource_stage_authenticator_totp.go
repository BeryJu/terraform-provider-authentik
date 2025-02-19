package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourceStageAuthenticatorTOTP() *schema.Resource {
	return &schema.Resource{
		Description:   "Flows & Stages --- ",
		CreateContext: resourceStageAuthenticatorTOTPCreate,
		ReadContext:   resourceStageAuthenticatorTOTPRead,
		UpdateContext: resourceStageAuthenticatorTOTPUpdate,
		DeleteContext: resourceStageAuthenticatorTOTPDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"friendly_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"configure_flow": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"digits": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.DIGITSENUM__6,
				Description:      EnumToDescription(api.AllowedDigitsEnumEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedDigitsEnumEnumValues),
			},
		},
	}
}

func resourceStageAuthenticatorTOTPSchemaToProvider(d *schema.ResourceData) *api.AuthenticatorTOTPStageRequest {
	r := api.AuthenticatorTOTPStageRequest{
		Name:   d.Get("name").(string),
		Digits: api.DigitsEnum(d.Get("digits").(string)),
	}

	if fn, fnSet := d.GetOk("friendly_name"); fnSet {
		r.FriendlyName.Set(api.PtrString(fn.(string)))
	} else {
		r.FriendlyName.Set(nil)
	}
	if h, hSet := d.GetOk("configure_flow"); hSet {
		r.ConfigureFlow.Set(api.PtrString(h.(string)))
	} else {
		r.ConfigureFlow.Set(nil)
	}
	return &r
}

func resourceStageAuthenticatorTOTPCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStageAuthenticatorTOTPSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorTotpCreate(ctx).AuthenticatorTOTPStageRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageAuthenticatorTOTPRead(ctx, d, m)
}

func resourceStageAuthenticatorTOTPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorTotpRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "digits", res.Digits)
	setWrapper(d, "friendly_name", res.FriendlyName.Get())
	if res.ConfigureFlow.IsSet() {
		setWrapper(d, "configure_flow", res.ConfigureFlow.Get())
	}
	return diags
}

func resourceStageAuthenticatorTOTPUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStageAuthenticatorTOTPSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorTotpUpdate(ctx, d.Id()).AuthenticatorTOTPStageRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageAuthenticatorTOTPRead(ctx, d, m)
}

func resourceStageAuthenticatorTOTPDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesAuthenticatorTotpDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
