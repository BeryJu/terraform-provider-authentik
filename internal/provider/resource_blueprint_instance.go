package provider

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourceBlueprintInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBlueprintInstanceCreate,
		ReadContext:   resourceBlueprintInstanceRead,
		UpdateContext: resourceBlueprintInstanceUpdate,
		DeleteContext: resourceBlueprintInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"context": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "{}",
			},
		},
	}
}

func resourceBlueprintInstanceSchemaToModel(d *schema.ResourceData, c *APIClient) (*api.BlueprintInstanceRequest, diag.Diagnostics) {
	m := api.BlueprintInstanceRequest{
		Name:    d.Get("name").(string),
		Path:    d.Get("path").(string),
		Enabled: boolToPointer(d.Get("enabled").(bool)),
	}

	ctx := make(map[string]interface{})
	if l, ok := d.Get("context").(string); ok {
		if l != "" {
			err := json.NewDecoder(strings.NewReader(l)).Decode(&ctx)
			if err != nil {
				return nil, diag.FromErr(err)
			}
		}
	}
	m.Context = ctx
	return &m, nil
}

func resourceBlueprintInstanceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app, diags := resourceBlueprintInstanceSchemaToModel(d, c)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.ManagedApi.ManagedBlueprintsCreate(ctx).BlueprintInstanceRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)

	return resourceBlueprintInstanceRead(ctx, d, m)
}

func resourceBlueprintInstanceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.ManagedApi.ManagedBlueprintsRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "path", res.Path)
	setWrapper(d, "enabled", res.Enabled)
	b, err := json.Marshal(res.Context)
	if err != nil {
		return diag.FromErr(err)
	}
	setWrapper(d, "context", string(b))
	return diags
}

func resourceBlueprintInstanceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app, diags := resourceBlueprintInstanceSchemaToModel(d, c)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.ManagedApi.ManagedBlueprintsUpdate(ctx, d.Id()).BlueprintInstanceRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceBlueprintInstanceRead(ctx, d, m)
}

func resourceBlueprintInstanceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.ManagedApi.ManagedBlueprintsDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
