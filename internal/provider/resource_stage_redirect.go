package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourceStageRedirect() *schema.Resource {
	return &schema.Resource{
		Description:   "Flows & Stages --- ",
		CreateContext: resourceStageRedirectCreate,
		ReadContext:   resourceStageRedirectRead,
		UpdateContext: resourceStageRedirectUpdate,
		DeleteContext: resourceStageRedirectDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mode": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.REDIRECTSTAGEMODEENUM_FLOW,
				Description:      EnumToDescription(api.AllowedRedirectStageModeEnumEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedRedirectStageModeEnumEnumValues),
			},
			"keep_context": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"target_static": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"target_flow": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceStageRedirectSchemaToProvider(d *schema.ResourceData) *api.RedirectStageRequest {
	r := api.RedirectStageRequest{
		Name:        d.Get("name").(string),
		Mode:        api.RedirectStageModeEnum(d.Get("mode").(string)),
		KeepContext: api.PtrBool(d.Get("keep_context").(bool)),
	}

	if target, targetSet := d.GetOk("target_static"); targetSet {
		r.TargetStatic = api.PtrString(target.(string))
	}
	if target, targetSet := d.GetOk("target_flow"); targetSet {
		r.TargetFlow.Set(api.PtrString(target.(string)))
	} else {
		r.TargetFlow.Set(nil)
	}
	return &r
}

func resourceStageRedirectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStageRedirectSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesRedirectCreate(ctx).RedirectStageRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageRedirectRead(ctx, d, m)
}

func resourceStageRedirectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesRedirectRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "mode", res.Mode)
	setWrapper(d, "keep_context", res.KeepContext)
	setWrapper(d, "target_flow", res.TargetFlow.Get())
	setWrapper(d, "target_static", res.TargetStatic)
	return diags
}

func resourceStageRedirectUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStageRedirectSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesRedirectUpdate(ctx, d.Id()).RedirectStageRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageRedirectRead(ctx, d, m)
}

func resourceStageRedirectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesRedirectDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
