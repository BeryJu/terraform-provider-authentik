package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourceProviderOAuth2() *schema.Resource {
	return &schema.Resource{
		Description:   "Applications --- ",
		CreateContext: resourceProviderOAuth2Create,
		ReadContext:   resourceProviderOAuth2Read,
		UpdateContext: resourceProviderOAuth2Update,
		DeleteContext: resourceProviderOAuth2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"authentication_flow": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"authorization_flow": {
				Type:     schema.TypeString,
				Required: true,
			},
			"invalidation_flow": {
				Type:     schema.TypeString,
				Required: true,
			},
			"property_mappings": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"client_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.CLIENTTYPEENUM_CONFIDENTIAL,
				Description:      EnumToDescription(api.AllowedClientTypeEnumEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedClientTypeEnumEnumValues),
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"client_secret": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				Computed:  true,
			},
			"access_code_validity": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "minutes=1",
			},
			"access_token_validity": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "minutes=10",
			},
			"refresh_token_validity": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "days=30",
			},
			"include_claims_in_id_token": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"signing_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"encryption_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"allowed_redirect_uris": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
				},
			},
			"sub_mode": {
				Type:             schema.TypeString,
				Default:          api.SUBMODEENUM_HASHED_USER_ID,
				Optional:         true,
				Description:      EnumToDescription(api.AllowedSubModeEnumEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedSubModeEnumEnumValues),
			},
			"issuer_mode": {
				Type:             schema.TypeString,
				Default:          api.ISSUERMODEENUM_PER_PROVIDER,
				Optional:         true,
				Description:      EnumToDescription(api.AllowedIssuerModeEnumEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedIssuerModeEnumEnumValues),
			},
			"jwks_sources": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Deprecated. Use `jwt_federation_sources` instead.",
			},
			"jwt_federation_sources": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "JWTs issued by keys configured in any of the selected sources can be used to authenticate on behalf of this provider.",
			},
			"jwt_federation_providers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "JWTs issued by any of the configured providers can be used to authenticate on behalf of this provider.",
			},
		},
	}
}

func resourceProviderOAuth2SchemaToProvider(d *schema.ResourceData) *api.OAuth2ProviderRequest {
	r := api.OAuth2ProviderRequest{
		Name:                   d.Get("name").(string),
		AuthorizationFlow:      d.Get("authorization_flow").(string),
		InvalidationFlow:       d.Get("invalidation_flow").(string),
		AccessCodeValidity:     api.PtrString(d.Get("access_code_validity").(string)),
		AccessTokenValidity:    api.PtrString(d.Get("access_token_validity").(string)),
		RefreshTokenValidity:   api.PtrString(d.Get("refresh_token_validity").(string)),
		IncludeClaimsInIdToken: api.PtrBool(d.Get("include_claims_in_id_token").(bool)),
		ClientId:               api.PtrString(d.Get("client_id").(string)),
		IssuerMode:             api.IssuerModeEnum(d.Get("issuer_mode").(string)).Ptr(),
		SubMode:                api.SubModeEnum(d.Get("sub_mode").(string)).Ptr(),
		ClientType:             api.ClientTypeEnum(d.Get("client_type").(string)).Ptr(),
		PropertyMappings:       castSlice[string](d.Get("property_mappings").([]interface{})),
		JwtFederationSources:   castSlice[string](d.Get("jwt_federation_sources").([]interface{})),
	}

	if s, sok := d.GetOk("authentication_flow"); sok && s.(string) != "" {
		r.AuthenticationFlow.Set(api.PtrString(s.(string)))
	} else {
		r.AuthenticationFlow.Set(nil)
	}
	if s, sok := d.GetOk("client_secret"); sok && s.(string) != "" {
		r.ClientSecret = api.PtrString(s.(string))
	}

	if s, sok := d.GetOk("signing_key"); sok && s.(string) != "" {
		r.SigningKey.Set(api.PtrString(s.(string)))
	} else {
		r.SigningKey.Set(nil)
	}
	if s, sok := d.GetOk("encryption_key"); sok && s.(string) != "" {
		r.EncryptionKey.Set(api.PtrString(s.(string)))
	} else {
		r.EncryptionKey.Set(nil)
	}

	r.RedirectUris = listToRedirectURIsRequest(d.Get("allowed_redirect_uris").([]interface{}))

	providers := d.Get("jwt_federation_providers").([]interface{})
	r.JwtFederationProviders = make([]int32, len(providers))
	for i, prov := range providers {
		r.JwtFederationProviders[i] = int32(prov.(int))
	}
	return &r
}

func listToRedirectURIsRequest(raw []interface{}) []api.RedirectURIRequest {
	rus := []api.RedirectURIRequest{}
	for _, rr := range raw {
		rd := rr.(map[string]interface{})
		rus = append(rus, api.RedirectURIRequest{
			MatchingMode: api.MatchingModeEnum(rd["matching_mode"].(string)),
			Url:          rd["url"].(string),
		})
	}
	return rus
}

func listToRedirectURIs(raw []interface{}) []api.RedirectURI {
	rus := []api.RedirectURI{}
	for _, rr := range raw {
		rd := rr.(map[string]interface{})
		rus = append(rus, api.RedirectURI{
			MatchingMode: api.MatchingModeEnum(rd["matching_mode"].(string)),
			Url:          rd["url"].(string),
		})
	}
	return rus
}

func redirectURIsToList(raw []api.RedirectURI) []map[string]interface{} {
	rus := []map[string]interface{}{}
	for _, rr := range raw {
		rus = append(rus, map[string]interface{}{
			"matching_mode": string(rr.MatchingMode),
			"url":           rr.Url,
		})
	}
	return rus
}

func resourceProviderOAuth2Create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceProviderOAuth2SchemaToProvider(d)

	res, hr, err := c.client.ProvidersApi.ProvidersOauth2Create(ctx).OAuth2ProviderRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderOAuth2Read(ctx, d, m)
}

func resourceProviderOAuth2Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	res, hr, err := c.client.ProvidersApi.ProvidersOauth2Retrieve(ctx, int32(id)).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "authentication_flow", res.AuthenticationFlow.Get())
	setWrapper(d, "authorization_flow", res.AuthorizationFlow)
	setWrapper(d, "invalidation_flow", res.InvalidationFlow)
	setWrapper(d, "client_id", res.ClientId)
	setWrapper(d, "client_secret", res.ClientSecret)
	setWrapper(d, "client_type", res.ClientType)
	setWrapper(d, "include_claims_in_id_token", res.IncludeClaimsInIdToken)
	setWrapper(d, "issuer_mode", res.IssuerMode)
	localMappings := castSlice[string](d.Get("property_mappings").([]interface{}))
	setWrapper(d, "property_mappings", listConsistentMerge(localMappings, res.PropertyMappings))
	localRedirectURIs := listToRedirectURIs(d.Get("allowed_redirect_uris").([]interface{}))
	setWrapper(d, "allowed_redirect_uris", redirectURIsToList(castSlice[api.RedirectURI](listConsistentMerge(localRedirectURIs, res.RedirectUris))))
	if res.SigningKey.IsSet() {
		setWrapper(d, "signing_key", res.SigningKey.Get())
	}
	if res.EncryptionKey.IsSet() {
		setWrapper(d, "encryption_key", res.EncryptionKey.Get())
	}
	setWrapper(d, "sub_mode", res.SubMode)
	setWrapper(d, "access_code_validity", res.AccessCodeValidity)
	setWrapper(d, "access_token_validity", res.AccessTokenValidity)
	setWrapper(d, "refresh_token_validity", res.RefreshTokenValidity)
	localJWKSProviders := castSlice[int](d.Get("jwt_federation_providers").([]interface{}))
	setWrapper(d, "jwt_federation_providers", listConsistentMerge(localJWKSProviders, slice32ToInt(res.JwtFederationProviders)))
	localJWKSSources := castSlice[string](d.Get("jwt_federation_sources").([]interface{}))
	setWrapper(d, "jwt_federation_sources", listConsistentMerge(localJWKSSources, res.JwtFederationSources))
	return diags
}

func resourceProviderOAuth2Update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	app := resourceProviderOAuth2SchemaToProvider(d)

	res, hr, err := c.client.ProvidersApi.ProvidersOauth2Update(ctx, int32(id)).OAuth2ProviderRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderOAuth2Read(ctx, d, m)
}

func resourceProviderOAuth2Delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	hr, err := c.client.ProvidersApi.ProvidersOauth2Destroy(ctx, int32(id)).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
