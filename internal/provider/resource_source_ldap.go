package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourceSourceLDAP() *schema.Resource {
	return &schema.Resource{
		Description:   "Directory --- ",
		CreateContext: resourceSourceLDAPCreate,
		ReadContext:   resourceSourceLDAPRead,
		UpdateContext: resourceSourceLDAPUpdate,
		DeleteContext: resourceSourceLDAPDelete,
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
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"server_uri": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bind_cn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bind_password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"start_tls": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"sni": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"base_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"additional_user_dn": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"additional_group_dn": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"user_object_filter": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "(objectClass=person)",
			},
			"group_object_filter": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "(objectClass=group)",
			},
			"group_membership_field": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "member",
			},
			"object_uniqueness_field": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "objectSid",
			},
			"sync_users": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"sync_users_password": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"sync_groups": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"sync_parent_group": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password_login_update_internal_password": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"property_mappings": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"property_mappings_group": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceSourceLDAPSchemaToSource(d *schema.ResourceData) *api.LDAPSourceRequest {
	r := api.LDAPSourceRequest{
		Name:             d.Get("name").(string),
		Slug:             d.Get("slug").(string),
		Enabled:          api.PtrBool(d.Get("enabled").(bool)),
		UserPathTemplate: api.PtrString(d.Get("user_path_template").(string)),

		BaseDn:       d.Get("base_dn").(string),
		ServerUri:    d.Get("server_uri").(string),
		BindCn:       api.PtrString(d.Get("bind_cn").(string)),
		BindPassword: api.PtrString(d.Get("bind_password").(string)),
		StartTls:     api.PtrBool(d.Get("start_tls").(bool)),
		Sni:          api.PtrBool(d.Get("sni").(bool)),

		AdditionalUserDn:      api.PtrString(d.Get("additional_user_dn").(string)),
		AdditionalGroupDn:     api.PtrString(d.Get("additional_group_dn").(string)),
		UserObjectFilter:      api.PtrString(d.Get("user_object_filter").(string)),
		GroupObjectFilter:     api.PtrString(d.Get("group_object_filter").(string)),
		GroupMembershipField:  api.PtrString(d.Get("group_membership_field").(string)),
		ObjectUniquenessField: api.PtrString(d.Get("object_uniqueness_field").(string)),

		SyncUsers:                           api.PtrBool(d.Get("sync_users").(bool)),
		SyncUsersPassword:                   api.PtrBool(d.Get("sync_users_password").(bool)),
		SyncGroups:                          api.PtrBool(d.Get("sync_groups").(bool)),
		PasswordLoginUpdateInternalPassword: api.PtrBool(d.Get("password_login_update_internal_password").(bool)),
	}

	if s, sok := d.GetOk("sync_parent_group"); sok && s.(string) != "" {
		r.SyncParentGroup.Set(api.PtrString(s.(string)))
	} else {
		r.SyncParentGroup.Set(nil)
	}

	r.UserPropertyMappings = castSlice[string](d.Get("property_mappings").([]interface{}))
	r.GroupPropertyMappings = castSlice[string](d.Get("property_mappings_group").([]interface{}))
	return &r
}

func resourceSourceLDAPCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceSourceLDAPSchemaToSource(d)

	res, hr, err := c.client.SourcesApi.SourcesLdapCreate(ctx).LDAPSourceRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Slug)
	return resourceSourceLDAPRead(ctx, d, m)
}

func resourceSourceLDAPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)
	res, hr, err := c.client.SourcesApi.SourcesLdapRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "slug", res.Slug)
	setWrapper(d, "uuid", res.Pk)
	setWrapper(d, "enabled", res.Enabled)
	setWrapper(d, "user_path_template", res.UserPathTemplate)

	setWrapper(d, "base_dn", res.BaseDn)
	setWrapper(d, "server_uri", res.ServerUri)
	setWrapper(d, "bind_cn", res.BindCn)
	setWrapper(d, "start_tls", res.StartTls)
	setWrapper(d, "sni", res.Sni)
	setWrapper(d, "base_dn", res.BaseDn)
	setWrapper(d, "additional_user_dn", res.AdditionalUserDn)
	setWrapper(d, "additional_group_dn", res.AdditionalGroupDn)
	setWrapper(d, "user_object_filter", res.UserObjectFilter)
	setWrapper(d, "group_object_filter", res.GroupObjectFilter)
	setWrapper(d, "group_membership_field", res.GroupMembershipField)
	setWrapper(d, "object_uniqueness_field", res.ObjectUniquenessField)
	setWrapper(d, "sync_users", res.SyncUsers)
	setWrapper(d, "sync_users_password", res.SyncUsersPassword)
	setWrapper(d, "sync_groups", res.SyncGroups)
	setWrapper(d, "password_login_update_internal_password", res.PasswordLoginUpdateInternalPassword)
	if res.SyncParentGroup.IsSet() {
		setWrapper(d, "sync_parent_group", res.SyncParentGroup.Get())
	}
	localMappings := castSlice[string](d.Get("property_mappings").([]interface{}))
	setWrapper(d, "property_mappings", listConsistentMerge(localMappings, res.UserPropertyMappings))
	localGroupMappings := castSlice[string](d.Get("property_mappings_group").([]interface{}))
	setWrapper(d, "property_mappings_group", listConsistentMerge(localGroupMappings, res.GroupPropertyMappings))
	return diags
}

func resourceSourceLDAPUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	app := resourceSourceLDAPSchemaToSource(d)

	res, hr, err := c.client.SourcesApi.SourcesLdapUpdate(ctx, d.Id()).LDAPSourceRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Slug)
	return resourceSourceLDAPRead(ctx, d, m)
}

func resourceSourceLDAPDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.SourcesApi.SourcesLdapDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
