package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourceSourcePlex() *schema.Resource {
	return &schema.Resource{
		Description:   "Directory --- ",
		CreateContext: resourceSourcePlexCreate,
		ReadContext:   resourceSourcePlexRead,
		UpdateContext: resourceSourcePlexUpdate,
		DeleteContext: resourceSourcePlexDelete,
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

			"client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"allowed_servers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"allow_friends": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"plex_token": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceSourcePlexSchemaToSource(d *schema.ResourceData) *api.PlexSourceRequest {
	r := api.PlexSourceRequest{
		Name:              d.Get("name").(string),
		Slug:              d.Get("slug").(string),
		Enabled:           api.PtrBool(d.Get("enabled").(bool)),
		UserPathTemplate:  api.PtrString(d.Get("user_path_template").(string)),
		PolicyEngineMode:  api.PolicyEngineMode(d.Get("policy_engine_mode").(string)).Ptr(),
		UserMatchingMode:  api.UserMatchingModeEnum(d.Get("user_matching_mode").(string)).Ptr(),
		GroupMatchingMode: api.GroupMatchingModeEnum(d.Get("group_matching_mode").(string)).Ptr(),

		ClientId:     api.PtrString(d.Get("client_id").(string)),
		AllowFriends: api.PtrBool(d.Get("allow_friends").(bool)),
		PlexToken:    d.Get("plex_token").(string),
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

	r.AllowedServers = castSlice[string](d.Get("allowed_servers").([]interface{}))
	return &r
}

func resourceSourcePlexCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceSourcePlexSchemaToSource(d)

	res, hr, err := c.client.SourcesApi.SourcesPlexCreate(ctx).PlexSourceRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Slug)
	return resourceSourcePlexRead(ctx, d, m)
}

func resourceSourcePlexRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)
	res, hr, err := c.client.SourcesApi.SourcesPlexRetrieve(ctx, d.Id()).Execute()
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
	setWrapper(d, "group_matching_mode", res.GroupMatchingMode)

	setWrapper(d, "client_id", res.ClientId)
	localServers := castSlice[string](d.Get("allowed_servers").([]interface{}))
	setWrapper(d, "allowed_servers", listConsistentMerge(localServers, res.AllowedServers))
	setWrapper(d, "allow_friends", res.AllowFriends)
	setWrapper(d, "plex_token", res.PlexToken)
	return diags
}

func resourceSourcePlexUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	app := resourceSourcePlexSchemaToSource(d)

	res, hr, err := c.client.SourcesApi.SourcesPlexUpdate(ctx, d.Id()).PlexSourceRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Slug)
	return resourceSourcePlexRead(ctx, d, m)
}

func resourceSourcePlexDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.SourcesApi.SourcesPlexDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
