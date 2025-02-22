package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourcePropertyMappingProviderSAML() *schema.Resource {
	return &schema.Resource{
		Description:   "Customization --- Manage SAML Provider Property mappings",
		CreateContext: resourcePropertyMappingProviderSAMLCreate,
		ReadContext:   resourcePropertyMappingProviderSAMLRead,
		UpdateContext: resourcePropertyMappingProviderSAMLUpdate,
		DeleteContext: resourcePropertyMappingProviderSAMLDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"saml_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"friendly_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"expression": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: diffSuppressExpression,
			},
		},
	}
}

func resourcePropertyMappingProviderSAMLSchemaToProvider(d *schema.ResourceData) *api.SAMLPropertyMappingRequest {
	r := api.SAMLPropertyMappingRequest{
		Name:       d.Get("name").(string),
		SamlName:   d.Get("saml_name").(string),
		Expression: d.Get("expression").(string),
	}
	if de, dSet := d.GetOk("friendly_name"); dSet {
		r.FriendlyName.Set(api.PtrString(de.(string)))
	} else {
		r.FriendlyName.Set(nil)
	}
	return &r
}

func resourcePropertyMappingProviderSAMLCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourcePropertyMappingProviderSAMLSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsProviderSamlCreate(ctx).SAMLPropertyMappingRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePropertyMappingProviderSAMLRead(ctx, d, m)
}

func resourcePropertyMappingProviderSAMLRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsProviderSamlRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "expression", res.Expression)
	setWrapper(d, "saml_name", res.SamlName)
	if res.FriendlyName.IsSet() {
		setWrapper(d, "friendly_name", res.FriendlyName.Get())
	} else {
		setWrapper(d, "friendly_name", nil)
	}
	return diags
}

func resourcePropertyMappingProviderSAMLUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourcePropertyMappingProviderSAMLSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsProviderSamlUpdate(ctx, d.Id()).SAMLPropertyMappingRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePropertyMappingProviderSAMLRead(ctx, d, m)
}

func resourcePropertyMappingProviderSAMLDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PropertymappingsApi.PropertymappingsProviderSamlDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
