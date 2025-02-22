package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourceStageAuthenticatorDuo() *schema.Resource {
	return &schema.Resource{
		Description:   "Flows & Stages --- ",
		CreateContext: resourceStageAuthenticatorDuoCreate,
		ReadContext:   resourceStageAuthenticatorDuoRead,
		UpdateContext: resourceStageAuthenticatorDuoUpdate,
		DeleteContext: resourceStageAuthenticatorDuoDelete,
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
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"client_secret": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"admin_integration_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"admin_secret_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"api_hostname": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceStageAuthenticatorDuoSchemaToProvider(d *schema.ResourceData) *api.AuthenticatorDuoStageRequest {
	r := api.AuthenticatorDuoStageRequest{
		Name:         d.Get("name").(string),
		ClientId:     d.Get("client_id").(string),
		ClientSecret: d.Get("client_secret").(string),
		ApiHostname:  d.Get("api_hostname").(string),
	}

	if fn, fnSet := d.GetOk("friendly_name"); fnSet {
		r.FriendlyName.Set(api.PtrString(fn.(string)))
	} else {
		r.FriendlyName.Set(nil)
	}
	if h, hSet := d.GetOk("admin_integration_key"); hSet {
		r.AdminIntegrationKey = api.PtrString(h.(string))
	}
	if h, hSet := d.GetOk("admin_secret_key"); hSet {
		r.AdminSecretKey = api.PtrString(h.(string))
	}
	if h, hSet := d.GetOk("configure_flow"); hSet {
		r.ConfigureFlow.Set(api.PtrString(h.(string)))
	} else {
		r.ConfigureFlow.Set(nil)
	}
	return &r
}

func resourceStageAuthenticatorDuoCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStageAuthenticatorDuoSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorDuoCreate(ctx).AuthenticatorDuoStageRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageAuthenticatorDuoRead(ctx, d, m)
}

func resourceStageAuthenticatorDuoRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorDuoRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "client_id", res.ClientId)
	setWrapper(d, "admin_integration_key", res.AdminIntegrationKey)
	setWrapper(d, "api_hostname", res.ApiHostname)
	setWrapper(d, "friendly_name", res.FriendlyName.Get())
	if res.ConfigureFlow.IsSet() {
		setWrapper(d, "configure_flow", res.ConfigureFlow.Get())
	}
	return diags
}

func resourceStageAuthenticatorDuoUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStageAuthenticatorDuoSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorDuoUpdate(ctx, d.Id()).AuthenticatorDuoStageRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageAuthenticatorDuoRead(ctx, d, m)
}

func resourceStageAuthenticatorDuoDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesAuthenticatorDuoDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
