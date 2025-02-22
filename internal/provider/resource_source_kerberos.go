package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourceSourceKerberos() *schema.Resource {
	return &schema.Resource{
		Description:   "Directory --- ",
		CreateContext: resourceSourceKerberosCreate,
		ReadContext:   resourceSourceKerberosRead,
		UpdateContext: resourceSourceKerberosUpdate,
		DeleteContext: resourceSourceKerberosDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_path_template": {
				Type:     schema.TypeString,
				Default:  "goauthentik.io/sources/%(slug)s",
				Optional: true,
			},
			"authentication_flow": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enrollment_flow": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"policy_engine_mode": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.POLICYENGINEMODE_ANY,
				Description:      EnumToDescription(api.AllowedPolicyEngineModeEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedPolicyEngineModeEnumValues),
			},
			"user_matching_mode": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.USERMATCHINGMODEENUM_IDENTIFIER,
				Description:      EnumToDescription(api.AllowedUserMatchingModeEnumEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedUserMatchingModeEnumEnumValues),
			},
			"group_matching_mode": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.GROUPMATCHINGMODEENUM_IDENTIFIER,
				Description:      EnumToDescription(api.AllowedGroupMatchingModeEnumEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedGroupMatchingModeEnumEnumValues),
			},

			"realm": {
				Description: "Kerberos realm",
				Type:        schema.TypeString,
				Required:    true,
			},
			"krb5_conf": {
				Description: "Custom krb5.conf to use. Uses the system one by default",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"sync_users": {
				Description: "Sync users from Kerberos into authentik",
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
			},
			"sync_users_password": {
				Description: "When a user changes their password, sync it back to Kerberos",
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
			},
			"sync_principal": {
				Description: "Principal to authenticate to kadmin for sync.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"sync_password": {
				Description: "Password to authenticate to kadmin for sync",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
			},
			"sync_keytab": {
				Description: "Keytab to authenticate to kadmin for sync. Must be base64-encoded or in the form TYPE:residual",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
			},
			"sync_ccache": {
				Description: "Credentials cache to authenticate to kadmin for sync. Must be in the form TYPE:residual",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"spnego_server_name": {
				Description: "Force the use of a specific server name for SPNEGO",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"spnego_keytab": {
				Description: "SPNEGO keytab base64-encoded or path to keytab in the form FILE:path",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
			},
			"spnego_ccache": {
				Description: "Credential cache to use for SPNEGO in form type:residual",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"password_login_update_internal_password": {
				Description: "If enabled, the authentik-stored password will be updated upon login with the Kerberos password backend",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func resourceSourceKerberosSchemaToSource(d *schema.ResourceData) (*api.KerberosSourceRequest, diag.Diagnostics) {
	r := api.KerberosSourceRequest{
		Name:             d.Get("name").(string),
		Slug:             d.Get("slug").(string),
		Enabled:          api.PtrBool(d.Get("enabled").(bool)),
		UserPathTemplate: api.PtrString(d.Get("user_path_template").(string)),

		PolicyEngineMode:  api.PolicyEngineMode(d.Get("policy_engine_mode").(string)).Ptr(),
		UserMatchingMode:  api.UserMatchingModeEnum(d.Get("user_matching_mode").(string)).Ptr(),
		GroupMatchingMode: api.GroupMatchingModeEnum(d.Get("group_matching_mode").(string)).Ptr(),

		Realm:                               d.Get("realm").(string),
		Krb5Conf:                            api.PtrString(d.Get("krb5_conf").(string)),
		SyncUsers:                           api.PtrBool(d.Get("sync_users").(bool)),
		SyncUsersPassword:                   api.PtrBool(d.Get("sync_users_password").(bool)),
		SyncPrincipal:                       api.PtrString(d.Get("sync_principal").(string)),
		SyncPassword:                        api.PtrString(d.Get("sync_password").(string)),
		SyncKeytab:                          api.PtrString(d.Get("sync_keytab").(string)),
		SyncCcache:                          api.PtrString(d.Get("sync_ccache").(string)),
		SpnegoServerName:                    api.PtrString(d.Get("spnego_server_name").(string)),
		SpnegoKeytab:                        api.PtrString(d.Get("spnego_keytab").(string)),
		SpnegoCcache:                        api.PtrString(d.Get("spnego_ccache").(string)),
		PasswordLoginUpdateInternalPassword: api.PtrBool(d.Get("password_login_update_internal_password").(bool)),
	}

	if ak, ok := d.GetOk("authentication_flow"); ok {
		r.AuthenticationFlow.Set(api.PtrString(ak.(string)))
	} else {
		r.AuthenticationFlow.Set(nil)
	}
	if ef, ok := d.GetOk("enrollment_flow"); ok {
		r.EnrollmentFlow.Set(api.PtrString(ef.(string)))
	} else {
		r.EnrollmentFlow.Set(nil)
	}

	return &r, nil
}

func resourceSourceKerberosCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r, diags := resourceSourceKerberosSchemaToSource(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.SourcesApi.SourcesKerberosCreate(ctx).KerberosSourceRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Slug)
	return resourceSourceKerberosRead(ctx, d, m)
}

func resourceSourceKerberosRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)
	res, hr, err := c.client.SourcesApi.SourcesKerberosRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "slug", res.Slug)
	setWrapper(d, "uuid", res.Pk)
	setWrapper(d, "user_path_template", res.UserPathTemplate)
	if res.AuthenticationFlow.IsSet() {
		setWrapper(d, "authentication_flow", res.AuthenticationFlow.Get())
	}
	if res.EnrollmentFlow.IsSet() {
		setWrapper(d, "enrollment_flow", res.EnrollmentFlow.Get())
	}
	setWrapper(d, "enabled", res.Enabled)
	setWrapper(d, "policy_engine_mode", res.PolicyEngineMode)
	setWrapper(d, "user_matching_mode", res.UserMatchingMode)
	setWrapper(d, "group_matching_mode", res.UserMatchingMode)

	setWrapper(d, "realm", res.Realm)
	setWrapper(d, "krb5_conf", res.Krb5Conf)
	setWrapper(d, "sync_users", res.SyncUsers)
	setWrapper(d, "sync_users_password", res.SyncUsersPassword)
	setWrapper(d, "sync_principal", res.SyncPrincipal)
	setWrapper(d, "sync_ccache", res.SyncCcache)
	setWrapper(d, "spnego_server_name", res.SpnegoServerName)
	setWrapper(d, "spnego_ccache", res.SpnegoCcache)
	setWrapper(d, "password_login_update_internal_password", res.PasswordLoginUpdateInternalPassword)
	return diags
}

func resourceSourceKerberosUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	app, diags := resourceSourceKerberosSchemaToSource(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.SourcesApi.SourcesKerberosUpdate(ctx, d.Id()).KerberosSourceRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Slug)
	return resourceSourceKerberosRead(ctx, d, m)
}

func resourceSourceKerberosDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.SourcesApi.SourcesKerberosDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
