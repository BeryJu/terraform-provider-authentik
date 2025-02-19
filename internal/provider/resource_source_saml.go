package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourceSourceSAML() *schema.Resource {
	return &schema.Resource{
		Description:   "Directory --- ",
		CreateContext: resourceSourceSAMLCreate,
		ReadContext:   resourceSourceSAMLRead,
		UpdateContext: resourceSourceSAMLUpdate,
		DeleteContext: resourceSourceSAMLDelete,
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

			"pre_authentication_flow": {
				Type:     schema.TypeString,
				Required: true,
			},
			"issuer": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sso_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"slo_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"allow_idp_initiated": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"name_id_policy": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.NAMEIDPOLICYENUM__2_0NAMEID_FORMATPERSISTENT,
				Description:      EnumToDescription(api.AllowedNameIdPolicyEnumEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedNameIdPolicyEnumEnumValues),
			},
			"binding_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.BINDINGTYPEENUM_REDIRECT,
				Description:      EnumToDescription(api.AllowedBindingTypeEnumEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedBindingTypeEnumEnumValues),
			},
			"signing_kp": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"encryption_kp": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"verification_kp": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"digest_algorithm": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.DIGESTALGORITHMENUM__2001_04_XMLENCSHA256,
				Description:      EnumToDescription(api.AllowedDigestAlgorithmEnumEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedDigestAlgorithmEnumEnumValues),
			},
			"signature_algorithm": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.SIGNATUREALGORITHMENUM__2001_04_XMLDSIG_MORERSA_SHA256,
				Description:      EnumToDescription(api.AllowedSignatureAlgorithmEnumEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedSignatureAlgorithmEnumEnumValues),
			},
			"temporary_user_delete_after": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "days=1",
			},

			"metadata": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SAML Metadata",
			},
		},
	}
}

func resourceSourceSAMLSchemaToSource(d *schema.ResourceData) *api.SAMLSourceRequest {
	r := api.SAMLSourceRequest{
		Name:              d.Get("name").(string),
		Slug:              d.Get("slug").(string),
		Enabled:           api.PtrBool(d.Get("enabled").(bool)),
		UserPathTemplate:  api.PtrString(d.Get("user_path_template").(string)),
		PolicyEngineMode:  api.PolicyEngineMode(d.Get("policy_engine_mode").(string)).Ptr(),
		UserMatchingMode:  api.UserMatchingModeEnum(d.Get("user_matching_mode").(string)).Ptr(),
		GroupMatchingMode: api.GroupMatchingModeEnum(d.Get("group_matching_mode").(string)).Ptr(),

		PreAuthenticationFlow: d.Get("pre_authentication_flow").(string),

		SsoUrl:                   d.Get("sso_url").(string),
		Issuer:                   api.PtrString(d.Get("issuer").(string)),
		AllowIdpInitiated:        api.PtrBool(d.Get("allow_idp_initiated").(bool)),
		TemporaryUserDeleteAfter: api.PtrString(d.Get("temporary_user_delete_after").(string)),
		BindingType:              api.BindingTypeEnum(d.Get("binding_type").(string)).Ptr(),
		DigestAlgorithm:          api.DigestAlgorithmEnum(d.Get("digest_algorithm").(string)).Ptr(),
		SignatureAlgorithm:       api.SignatureAlgorithmEnum(d.Get("signature_algorithm").(string)).Ptr(),
		NameIdPolicy:             api.NameIdPolicyEnum(d.Get("name_id_policy").(string)).Ptr(),
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

	if s, sok := d.GetOk("slo_url"); sok && s.(string) != "" {
		r.SloUrl.Set(api.PtrString(s.(string)))
	} else {
		r.SloUrl.Set(nil)
	}
	if s, sok := d.GetOk("signing_kp"); sok && s.(string) != "" {
		r.SigningKp.Set(api.PtrString(s.(string)))
	} else {
		r.SigningKp.Set(nil)
	}
	if s, sok := d.GetOk("encryption_kp"); sok && s.(string) != "" {
		r.EncryptionKp.Set(api.PtrString(s.(string)))
	} else {
		r.EncryptionKp.Set(nil)
	}
	if s, sok := d.GetOk("verification_kp"); sok && s.(string) != "" {
		r.VerificationKp.Set(api.PtrString(s.(string)))
	} else {
		r.VerificationKp.Set(nil)
	}

	return &r
}

func resourceSourceSAMLCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceSourceSAMLSchemaToSource(d)

	res, hr, err := c.client.SourcesApi.SourcesSamlCreate(ctx).SAMLSourceRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Slug)
	return resourceSourceSAMLRead(ctx, d, m)
}

func resourceSourceSAMLRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)
	res, hr, err := c.client.SourcesApi.SourcesSamlRetrieve(ctx, d.Id()).Execute()
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

	setWrapper(d, "pre_authentication_flow", res.PreAuthenticationFlow)
	setWrapper(d, "issuer", res.Issuer)
	setWrapper(d, "sso_url", res.SsoUrl)
	if res.SloUrl.IsSet() {
		setWrapper(d, "slo_url", res.SloUrl.Get())
	}
	setWrapper(d, "allow_idp_initiated", res.AllowIdpInitiated)
	setWrapper(d, "name_id_policy", res.NameIdPolicy)
	setWrapper(d, "binding_type", res.BindingType)
	if res.SigningKp.IsSet() {
		setWrapper(d, "signing_kp", res.SigningKp.Get())
	}
	if res.EncryptionKp.IsSet() {
		setWrapper(d, "encryption_kp", res.EncryptionKp.Get())
	}
	if res.VerificationKp.IsSet() {
		setWrapper(d, "verification_kp", res.VerificationKp.Get())
	}
	setWrapper(d, "digest_algorithm", res.DigestAlgorithm)
	setWrapper(d, "signature_algorithm", res.SignatureAlgorithm)
	setWrapper(d, "temporary_user_delete_after", res.TemporaryUserDeleteAfter)

	meta, hr, err := c.client.SourcesApi.SourcesSamlMetadataRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	setWrapper(d, "metadata", meta.Metadata)
	return diags
}

func resourceSourceSAMLUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	app := resourceSourceSAMLSchemaToSource(d)

	res, hr, err := c.client.SourcesApi.SourcesSamlUpdate(ctx, d.Id()).SAMLSourceRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Slug)
	return resourceSourceSAMLRead(ctx, d, m)
}

func resourceSourceSAMLDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.SourcesApi.SourcesSamlDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
