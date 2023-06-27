package spinnaker

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rmullinnix461332/terraform-provider-spinnaker/gateclient"
)

func New() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"server": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "URL for Gate",
				DefaultFunc: schema.EnvDefaultFunc("GATE_URL", nil),
			},
			"x509_cert": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "x509 Certificate for authenticating to Gate",
				DefaultFunc: schema.EnvDefaultFunc("GATE_X509_CERT", nil),
			},
			"x509_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "x509 private key for authenticating to Gate",
				DefaultFunc: schema.EnvDefaultFunc("GATE_X509_KEY", nil),
			},
			"ignore_cert_errors": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Ignore certificate errors from Gate",
				Default:     true,
			},
			"default_headers": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Headers to be passed to the gate endpoint by the client on each request",
				Default:     "",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"spinnaker_application":              resourceApplication(),
			"spinnaker_pipeline":                 resourcePipeline(),
			"spinnaker_pipeline_template":        resourcePipelineTemplate(),
			"spinnaker_pipeline_template_config": resourcePipelineTemplateConfig(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"spinnaker_pipeline": datasourcePipeline(),
		},
		ConfigureContextFunc: providerConfigure,
	}

	return p
}

type gateConfig struct {
	server string
	client *gateclient.GatewayClient
}

func providerConfigure(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
	server := data.Get("server").(string)
	x509_cert := data.Get("x509_cert").(string)
	x509_key := data.Get("x509_key").(string)
	ignoreCertErrors := data.Get("ignore_cert_errors").(bool)
	defaultHeaders := data.Get("default_headers").(string)

	client, err := gateclient.NewGateClient(server, defaultHeaders, x509_cert, x509_key, ignoreCertErrors)

	if err != nil {
		fmt.Println("config error", err)
	}
	return gateConfig{
		server: server,
		client: client,
	}, diag.Diagnostics{}
}
