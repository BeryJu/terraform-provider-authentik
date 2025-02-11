package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePolicyRead,
		Description: "Customization --- Get policy by name",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourcePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	req := c.client.PoliciesApi.PoliciesAllList(ctx)
	if s, ok := d.GetOk("name"); ok {
		req = req.Search(s.(string))
	}

	res, hr, err := req.Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	if len(res.Results) < 1 {
		return diag.Errorf("No matching policy found")
	}
	f := res.Results[0]
	d.SetId(f.Pk)
	setWrapper(d, "name", f.Name)
	return diags
}
