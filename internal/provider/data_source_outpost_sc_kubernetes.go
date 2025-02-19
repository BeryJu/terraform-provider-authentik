package provider

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataOutpostServiceConnectionsKubernetes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataOutpostServiceConnectionsKubernetesRead,
		Description: "Applications --- Get Service Connections for Kubernetes",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"kubeconfig": {
				Type:      schema.TypeString,
				Computed:  true,
				Optional:  true,
				Sensitive: true,
			},
			"local": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"verify_ssl": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func dataOutpostServiceConnectionsKubernetesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	req := c.client.OutpostsApi.OutpostsServiceConnectionsKubernetesList(ctx)
	if s, ok := d.GetOk("name"); ok {
		req = req.Name(s.(string))
	}

	res, hr, err := req.Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	if len(res.Results) < 1 {
		return diag.Errorf("No Kubernetes Outpost Service Connections found")
	}
	f := res.Results[0]
	d.SetId(f.Pk)
	setWrapper(d, "name", f.Name)
	setWrapper(d, "local", f.Local)
	setWrapper(d, "verify_ssl", f.VerifySsl)
	b, err := json.Marshal(f.Kubeconfig)
	if err != nil {
		return diag.FromErr(err)
	}
	setWrapper(d, "kubeconfig", string(b))
	return diags
}
