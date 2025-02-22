package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourceProviderLDAP() *schema.Resource {
	return &schema.Resource{
		Description:   "Applications --- ",
		CreateContext: resourceProviderLDAPCreate,
		ReadContext:   resourceProviderLDAPRead,
		UpdateContext: resourceProviderLDAPUpdate,
		DeleteContext: resourceProviderLDAPDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bind_flow": {
				Type:     schema.TypeString,
				Required: true,
			},
			"unbind_flow": {
				Type:     schema.TypeString,
				Required: true,
			},
			"base_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"certificate": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tls_server_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"uid_start_number": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2000,
			},
			"gid_start_number": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  4000,
			},
			"search_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  api.LDAPAPIACCESSMODE_DIRECT,
			},
			"bind_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  api.LDAPAPIACCESSMODE_DIRECT,
			},
			"mfa_support": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceProviderLDAPSchemaToProvider(d *schema.ResourceData) *api.LDAPProviderRequest {
	r := api.LDAPProviderRequest{
		Name:              d.Get("name").(string),
		AuthorizationFlow: d.Get("bind_flow").(string),
		InvalidationFlow:  d.Get("unbind_flow").(string),
		BaseDn:            api.PtrString(d.Get("base_dn").(string)),
		UidStartNumber:    api.PtrInt32(int32(d.Get("uid_start_number").(int))),
		GidStartNumber:    api.PtrInt32(int32(d.Get("gid_start_number").(int))),
		SearchMode:        api.LDAPAPIAccessMode(d.Get("search_mode").(string)).Ptr(),
		BindMode:          api.LDAPAPIAccessMode(d.Get("bind_mode").(string)).Ptr(),
		MfaSupport:        api.PtrBool(d.Get("mfa_support").(bool)),
	}

	if s, sok := d.GetOk("certificate"); sok && s.(string) != "" {
		r.Certificate.Set(api.PtrString(s.(string)))
	} else {
		r.Certificate.Set(nil)
	}
	if s, sok := d.GetOk("tls_server_name"); sok && s.(string) != "" {
		r.TlsServerName = api.PtrString(s.(string))
	}
	return &r
}

func resourceProviderLDAPCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceProviderLDAPSchemaToProvider(d)

	res, hr, err := c.client.ProvidersApi.ProvidersLdapCreate(ctx).LDAPProviderRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderLDAPRead(ctx, d, m)
}

func resourceProviderLDAPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	res, hr, err := c.client.ProvidersApi.ProvidersLdapRetrieve(ctx, int32(id)).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "bind_flow", res.AuthorizationFlow)
	setWrapper(d, "unbind_flow", res.InvalidationFlow)
	setWrapper(d, "base_dn", res.BaseDn)
	if res.Certificate.IsSet() {
		setWrapper(d, "certificate", res.Certificate.Get())
	}
	setWrapper(d, "tls_server_name", res.TlsServerName)
	setWrapper(d, "uid_start_number", res.UidStartNumber)
	setWrapper(d, "gid_start_number", res.GidStartNumber)
	setWrapper(d, "bind_mode", res.BindMode)
	setWrapper(d, "search_mode", res.SearchMode)
	setWrapper(d, "mfa_support", res.MfaSupport)
	return diags
}

func resourceProviderLDAPUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	app := resourceProviderLDAPSchemaToProvider(d)

	res, hr, err := c.client.ProvidersApi.ProvidersLdapUpdate(ctx, int32(id)).LDAPProviderRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderLDAPRead(ctx, d, m)
}

func resourceProviderLDAPDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	hr, err := c.client.ProvidersApi.ProvidersLdapDestroy(ctx, int32(id)).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
