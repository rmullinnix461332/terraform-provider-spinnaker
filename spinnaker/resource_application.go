package spinnaker

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApplication() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"application": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateApplicationName,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"platform_health_only": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"platform_health_only_show_override": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},
		Create: resourceApplicationCreate,
		Read:   resourceApplicationRead,
		Update: resourceApplicationUpdate,
		Delete: resourceApplicationDelete,
		Exists: resourceApplicationExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

type applicationRead struct {
	Name       string `json:"name"`
	Attributes struct {
		Email         string `json:"email"`
		Accounts      string `json:"accounts"`
		CloudProvider string `json:"cloudProviders"`
	} `json:"attributes"`
}

func resourceApplicationCreate(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	application := data.Get("application").(string)
	email := data.Get("email").(string)
	description := data.Get("description").(string)
	platform_health_only := data.Get("platform_health_only").(bool)
	platform_health_only_show_override := data.Get("platform_health_only_show_override").(bool)

	if err := client.CreateApplication(application, email, description, platform_health_only, platform_health_only_show_override); err != nil {
		return err
	}

	return resourceApplicationRead(data, meta)
}

func resourceApplicationRead(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	applicationName := data.Get("application").(string)
	fmt.Println(clientConfig, client, applicationName)
	var app applicationRead
	if err := client.GetApplication("docta", &app); err != nil {
		return err
	}

	return readApplication(data, app)
}

func resourceApplicationUpdate(data *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceApplicationDelete(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	applicationName := data.Get("application").(string)

	return client.DeleteAppliation(applicationName)
}

func resourceApplicationExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	applicationName := data.Get("application").(string)
	fmt.Println(clientConfig, client, applicationName)
	var app applicationRead
	if err := client.GetApplication("docta", &app); err != nil {
		errmsg := err.Error()
		if strings.Contains(errmsg, "not found") {
			return false, nil
		}
		return false, err
	}

	if app.Name == "" {
		return false, nil
	}

	return true, nil
}

func readApplication(data *schema.ResourceData, application applicationRead) error {
	data.SetId(application.Name)
	return nil
}
