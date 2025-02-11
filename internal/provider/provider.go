package provider

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/getsentry/sentry-go"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
		desc := s.Description
		if s.Default != nil {
			desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
		}
		if s.Computed {
			desc += " Generated."
		}
		return strings.TrimSpace(desc)
	}
}

// Provider -
func Provider(version string, testing bool) *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AUTHENTIK_URL", nil),
				Description: "The authentik API endpoint, can optionally be passed as `AUTHENTIK_URL` environmental variable",
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AUTHENTIK_INSECURE", false),
				Description: "Whether to skip TLS verification, can optionally be passed as `AUTHENTIK_INSECURE` environmental variable",
			},
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AUTHENTIK_TOKEN", nil),
				Sensitive:   true,
				Description: "The authentik API token, can optionally be passed as `AUTHENTIK_TOKEN` environmental variable",
			},
			"headers": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Sensitive:   true,
				Description: "Optional HTTP headers sent with every request",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"authentik_application_entitlement":                    tr(resourceApplicationEntitlement),
			"authentik_application":                                tr(resourceApplication),
			"authentik_blueprint":                                  tr(resourceBlueprintInstance),
			"authentik_brand":                                      tr(resourceBrand),
			"authentik_certificate_key_pair":                       tr(resourceCertificateKeyPair),
			"authentik_enterprise_license":                         tr(resourceEnterpriseLicense),
			"authentik_event_rule":                                 tr(resourceEventRule),
			"authentik_event_transport":                            tr(resourceEventTransport),
			"authentik_flow_stage_binding":                         tr(resourceFlowStageBinding),
			"authentik_flow":                                       tr(resourceFlow),
			"authentik_group":                                      tr(resourceGroup),
			"authentik_outpost":                                    tr(resourceOutpost),
			"authentik_policy_binding":                             tr(resourcePolicyBinding),
			"authentik_policy_dummy":                               tr(resourcePolicyDummy),
			"authentik_policy_event_matcher":                       tr(resourcePolicyEventMatcher),
			"authentik_policy_expiry":                              tr(resourcePolicyExpiry),
			"authentik_policy_expression":                          tr(resourcePolicyExpression),
			"authentik_policy_geoip":                               tr(resourcePolicyGeoIP),
			"authentik_policy_password":                            tr(resourcePolicyPassword),
			"authentik_policy_reputation":                          tr(resourcePolicyReputation),
			"authentik_property_mapping_notification":              tr(resourcePropertyMappingNotification),
			"authentik_property_mapping_provider_google_workspace": tr(resourcePropertyMappingProviderGoogleWorkspace),
			// TODO: Remove after 2024.10
			"authentik_property_mapping_google_workspace":         tr(markDeprecated(resourcePropertyMappingProviderGoogleWorkspace, "authentik_property_mapping_provider_google_workspace")),
			"authentik_property_mapping_provider_microsoft_entra": tr(resourcePropertyMappingProviderMicrosoftEntra),
			// TODO: Remove after 2024.10
			"authentik_property_mapping_microsoft_entra": tr(markDeprecated(resourcePropertyMappingProviderMicrosoftEntra, "authentik_property_mapping_provider_microsoft_entra")),
			"authentik_property_mapping_provider_rac":    tr(resourcePropertyMappingProviderRAC),
			// TODO: Remove after 2024.10
			"authentik_property_mapping_rac":             tr(markDeprecated(resourcePropertyMappingProviderRAC, "authentik_property_mapping_provider_rac")),
			"authentik_property_mapping_provider_radius": tr(resourcePropertyMappingProviderRadius),
			// TODO: Remove after 2024.10
			"authentik_property_mapping_radius":        tr(markDeprecated(resourcePropertyMappingProviderRadius, "authentik_property_mapping_provider_radius")),
			"authentik_property_mapping_provider_saml": tr(resourcePropertyMappingProviderSAML),
			// TODO: Remove after 2024.10
			"authentik_property_mapping_saml":          tr(markDeprecated(resourcePropertyMappingProviderSAML, "authentik_property_mapping_provider_saml")),
			"authentik_property_mapping_provider_scim": tr(resourcePropertyMappingProviderSCIM),
			// TODO: Remove after 2024.10
			"authentik_property_mapping_scim":           tr(markDeprecated(resourcePropertyMappingProviderSCIM, "authentik_property_mapping_provider_scim")),
			"authentik_property_mapping_provider_scope": tr(resourcePropertyMappingProviderScope),
			"authentik_property_mapping_source_ldap":    tr(resourcePropertyMappingSourceLDAP),
			// TODO: Remove after 2024.10
			"authentik_property_mapping_ldap":            tr(markDeprecated(resourcePropertyMappingSourceLDAP, "authentik_property_mapping_source_ldap")),
			"authentik_property_mapping_source_oauth":    tr(resourcePropertyMappingSourceOAuth),
			"authentik_property_mapping_source_plex":     tr(resourcePropertyMappingSourcePlex),
			"authentik_property_mapping_source_saml":     tr(resourcePropertyMappingSourceSAML),
			"authentik_property_mapping_source_scim":     tr(resourcePropertyMappingSourceSCIM),
			"authentik_property_mapping_source_kerberos": tr(resourcePropertyMappingSourceKerberos),
			"authentik_provider_google_workspace":        tr(resourceProviderGoogleWorkspace),
			"authentik_provider_ldap":                    tr(resourceProviderLDAP),
			"authentik_provider_microsoft_entra":         tr(resourceProviderMicrosoftEntra),
			"authentik_provider_oauth2":                  tr(resourceProviderOAuth2),
			"authentik_provider_proxy":                   tr(resourceProviderProxy),
			"authentik_provider_rac":                     tr(resourceProviderRAC),
			"authentik_provider_radius":                  tr(resourceProviderRadius),
			"authentik_provider_saml":                    tr(resourceProviderSAML),
			"authentik_provider_scim":                    tr(resourceProviderSCIM),
			"authentik_rac_endpoint":                     tr(resourceRACEndpoint),
			"authentik_rbac_permission_role":             tr(resourceRBACRoleObjectPermission),
			"authentik_rbac_permission_user":             tr(resourceRBACUserObjectPermission),
			"authentik_rbac_role":                        tr(resourceRBACRole),
			// TODO: Remove after 2024.10
			"authentik_scope_mapping":                     tr(markDeprecated(resourcePropertyMappingProviderScope, "authentik_property_mapping_provider_scope")),
			"authentik_service_connection_docker":         tr(resourceServiceConnectionDocker),
			"authentik_service_connection_kubernetes":     tr(resourceServiceConnectionKubernetes),
			"authentik_source_kerberos":                   tr(resourceSourceKerberos),
			"authentik_source_ldap":                       tr(resourceSourceLDAP),
			"authentik_source_oauth":                      tr(resourceSourceOAuth),
			"authentik_source_plex":                       tr(resourceSourcePlex),
			"authentik_source_saml":                       tr(resourceSourceSAML),
			"authentik_stage_authenticator_duo":           tr(resourceStageAuthenticatorDuo),
			"authentik_stage_authenticator_endpoint_gdtc": tr(resourceStageAuthenticatorEndpointGDTC),
			"authentik_stage_authenticator_sms":           tr(resourceStageAuthenticatorSms),
			"authentik_stage_authenticator_static":        tr(resourceStageAuthenticatorStatic),
			"authentik_stage_authenticator_totp":          tr(resourceStageAuthenticatorTOTP),
			"authentik_stage_authenticator_validate":      tr(resourceStageAuthenticatorValidate),
			"authentik_stage_authenticator_webauthn":      tr(resourceStageAuthenticatorWebAuthn),
			"authentik_stage_captcha":                     tr(resourceStageCaptcha),
			"authentik_stage_consent":                     tr(resourceStageConsent),
			"authentik_stage_deny":                        tr(resourceStageDeny),
			"authentik_stage_dummy":                       tr(resourceStageDummy),
			"authentik_stage_email":                       tr(resourceStageEmail),
			"authentik_stage_identification":              tr(resourceStageIdentification),
			"authentik_stage_invitation":                  tr(resourceStageInvitation),
			"authentik_stage_password":                    tr(resourceStagePassword),
			"authentik_stage_prompt_field":                tr(resourceStagePromptField),
			"authentik_stage_prompt":                      tr(resourceStagePrompt),
			"authentik_stage_redirect":                    tr(resourceStageRedirect),
			"authentik_stage_source":                      tr(resourceStageSource),
			"authentik_stage_user_delete":                 tr(resourceStageUserDelete),
			"authentik_stage_user_login":                  tr(resourceStageUserLogin),
			"authentik_stage_user_logout":                 tr(resourceStageUserLogout),
			"authentik_stage_user_write":                  tr(resourceStageUserWrite),
			"authentik_system_settings":                   tr(resourceSystemSettings),
			"authentik_token":                             tr(resourceToken),
			"authentik_user":                              tr(resourceUser),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"authentik_brand":                            td(dataSourceBrand),
			"authentik_certificate_key_pair":             td(dataSourceCertificateKeyPair),
			"authentik_flow":                             td(dataSourceFlow),
			"authentik_group":                            td(dataSourceGroup),
			"authentik_groups":                           td(dataSourceGroups),
			"authentik_policy":                           td(dataSourcePolicy),
			"authentik_property_mapping_provider_rac":    td(dataSourcePropertyMappingProviderRAC),
			"authentik_property_mapping_provider_radius": td(dataSourcePropertyMappingProviderRadius),
			"authentik_property_mapping_provider_saml":   td(dataSourcePropertyMappingProviderSAML),
			"authentik_property_mapping_provider_scim":   td(dataSourcePropertyMappingProviderSCIM),
			"authentik_property_mapping_provider_scope":  td(dataSourcePropertyMappingProviderScope),
			"authentik_property_mapping_source_ldap":     td(dataSourcePropertyMappingSourceLDAP),
			"authentik_provider_oauth2_config":           td(dataSourceProviderOAuth2Config),
			"authentik_provider_saml_metadata":           td(dataSourceProviderSAMLMetadata),
			"authentik_source":                           td(dataSourceSource),
			"authentik_stage":                            td(dataSourceStage),
			"authentik_user":                             td(dataSourceUser),
			"authentik_users":                            td(dataSourceUsers),
			"authentik_webauthn_device_type":             td(dataSourceWebAuthnDeviceType),
		},
		ConfigureContextFunc: providerConfigure(version, testing),
	}
}

// APIClient Hold the API Client and any relevant configuration
type APIClient struct {
	client *api.APIClient
}

func providerConfigure(version string, testing bool) schema.ConfigureContextFunc {
	return func(c context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		apiURL := d.Get("url").(string)
		token := d.Get("token").(string)
		insecure := d.Get("insecure").(bool)

		// Warning or errors can be collected in a slice type
		var diags diag.Diagnostics

		akURL, err := url.Parse(apiURL)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		config := api.NewConfiguration()
		config.Debug = true
		config.UserAgent = fmt.Sprintf("authentik-terraform@%s", version)
		config.Host = akURL.Host
		config.Scheme = akURL.Scheme
		if testing {
			config.HTTPClient = &http.Client{
				Transport: NewTestingTransport(GetTLSTransport(insecure)),
			}
		} else {
			config.HTTPClient = &http.Client{
				Transport: GetTLSTransport(insecure),
			}
		}

		config.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", token))
		if _headers, ok := d.GetOk("headers"); ok {
			headers := _headers.(map[string]any)
			for headerName, headerValue := range headers {
				config.AddDefaultHeader(headerName, headerValue.(string))
			}
		}
		apiClient := api.NewAPIClient(config)

		rootConfig, _, err := apiClient.RootApi.RootConfigRetrieve(context.Background()).Execute()
		if err == nil && rootConfig.ErrorReporting.Enabled {
			dsn := ""
			// Customisable Sentry DSN was added in 2022.11, so only use that DSN when its set
			if rootConfig.ErrorReporting.SentryDsn != "" {
				dsn = rootConfig.ErrorReporting.SentryDsn
			}
			if envDsn, found := os.LookupEnv("SENTRY_DSN"); found {
				dsn = envDsn
			}
			err := sentry.Init(sentry.ClientOptions{
				Dsn:              dsn,
				EnableTracing:    true,
				Environment:      rootConfig.ErrorReporting.Environment,
				TracesSampleRate: float64(rootConfig.ErrorReporting.TracesSampleRate),
				Release:          fmt.Sprintf("terraform-provider-authentik@%s", version),
			})
			if err != nil {
				fmt.Printf("Error during sentry init: %v\n", err)
			} else {
				config.HTTPClient.Transport = NewTracingTransport(context.Background(), config.HTTPClient.Transport)
				apiClient = api.NewAPIClient(config)
			}
		}

		return &APIClient{
			client: apiClient,
		}, diags
	}
}

// TestingTransport Transport used for testing, always returns a 400 Response
type TestingTransport struct {
	inner http.RoundTripper
}

// NewTestingTransport Get a HTTP Transport that fails all requests
func NewTestingTransport(inner http.RoundTripper) *TestingTransport {
	return &TestingTransport{inner}
}

// RoundTrip HTTP Transport
func (tt *TestingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	body := "mock-failed-request"
	return &http.Response{
		Status:        "400 Bad Request",
		StatusCode:    400,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          io.NopCloser(bytes.NewBufferString(body)),
		ContentLength: int64(len(body)),
		Request:       r,
		Header:        make(http.Header),
	}, nil
}

// GetTLSTransport Get a TLS transport instance, that skips verification if configured via environment variables.
func GetTLSTransport(insecure bool) http.RoundTripper {
	tlsTransport, err := httptransport.TLSTransport(httptransport.TLSClientOptions{
		InsecureSkipVerify: insecure,
	})
	if err != nil {
		panic(err)
	}
	return tlsTransport
}

type tracingTransport struct {
	inner http.RoundTripper
	ctx   context.Context
}

func NewTracingTransport(ctx context.Context, inner http.RoundTripper) *tracingTransport {
	return &tracingTransport{inner, ctx}
}

func (tt *tracingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	span := sentry.StartSpan(tt.ctx, "authentik.go.http_request")
	r.Header.Set("sentry-trace", span.ToSentryTrace())
	span.Description = fmt.Sprintf("%s %s", r.Method, r.URL.String())
	span.SetTag("url", r.URL.String())
	span.SetTag("method", r.Method)
	defer span.Finish()
	res, err := tt.inner.RoundTrip(r.WithContext(span.Context()))
	return res, err
}
