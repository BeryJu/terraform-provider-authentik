package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourceStageAuthenticatorSms() *schema.Resource {
	return &schema.Resource{
		Description:   "Flows & Stages --- ",
		CreateContext: resourceStageAuthenticatorSmsCreate,
		ReadContext:   resourceStageAuthenticatorSmsRead,
		UpdateContext: resourceStageAuthenticatorSmsUpdate,
		DeleteContext: resourceStageAuthenticatorSmsDelete,
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
			"sms_provider": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.PROVIDERENUM_TWILIO,
				Description:      EnumToDescription(api.AllowedProviderEnumEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedProviderEnumEnumValues),
			},
			"from_number": {
				Type:     schema.TypeString,
				Required: true,
			},
			"account_sid": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"auth": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"auth_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.AUTHTYPEENUM_BASIC,
				Description:      EnumToDescription(api.AllowedAuthTypeEnumEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedAuthTypeEnumEnumValues),
			},
			"auth_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"mapping": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"verify_only": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceStageAuthenticatorSmsSchemaToProvider(d *schema.ResourceData) *api.AuthenticatorSMSStageRequest {
	r := api.AuthenticatorSMSStageRequest{
		Name:       d.Get("name").(string),
		Provider:   api.ProviderEnum(d.Get("sms_provider").(string)),
		FromNumber: d.Get("from_number").(string),
		AccountSid: d.Get("account_sid").(string),
		AuthType:   api.AuthTypeEnum(d.Get("auth_type").(string)).Ptr(),
		Auth:       d.Get("auth").(string),
	}

	if fn, fnSet := d.GetOk("friendly_name"); fnSet {
		r.FriendlyName.Set(api.PtrString(fn.(string)))
	} else {
		r.FriendlyName.Set(nil)
	}
	if h, hSet := d.GetOk("auth_password"); hSet {
		r.AuthPassword = api.PtrString(h.(string))
	}
	if h, hSet := d.GetOk("configure_flow"); hSet {
		r.ConfigureFlow.Set(api.PtrString(h.(string)))
	} else {
		r.ConfigureFlow.Set(nil)
	}
	if h, hSet := d.GetOk("mapping"); hSet {
		r.Mapping.Set(api.PtrString(h.(string)))
	} else {
		r.Mapping.Set(nil)
	}
	if verify, verifySet := d.GetOk("verify_only"); verifySet {
		r.VerifyOnly = api.PtrBool(verify.(bool))
	}
	return &r
}

func resourceStageAuthenticatorSmsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStageAuthenticatorSmsSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorSmsCreate(ctx).AuthenticatorSMSStageRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageAuthenticatorSmsRead(ctx, d, m)
}

func resourceStageAuthenticatorSmsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorSmsRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "sms_provider", res.Provider)
	setWrapper(d, "from_number", res.FromNumber)
	setWrapper(d, "account_sid", res.AccountSid)
	setWrapper(d, "auth", res.Auth)
	setWrapper(d, "auth_password", res.AuthPassword)
	setWrapper(d, "auth_type", res.AuthType)
	setWrapper(d, "verify_only", res.VerifyOnly)
	setWrapper(d, "mapping", res.Mapping.Get())
	setWrapper(d, "friendly_name", res.FriendlyName.Get())
	if res.ConfigureFlow.IsSet() {
		setWrapper(d, "configure_flow", res.ConfigureFlow.Get())
	}
	return diags
}

func resourceStageAuthenticatorSmsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStageAuthenticatorSmsSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorSmsUpdate(ctx, d.Id()).AuthenticatorSMSStageRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageAuthenticatorSmsRead(ctx, d, m)
}

func resourceStageAuthenticatorSmsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesAuthenticatorSmsDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
