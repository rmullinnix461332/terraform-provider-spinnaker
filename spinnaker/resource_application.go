package spinnaker

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApplication() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateApplicationName,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cloud_providers": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "kubernetes",
			},
			"accounts": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  80,
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

func resourceApplicationCreate(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client

	application := data.Get("name").(string)
	email := data.Get("email").(string)
	description := data.Get("description").(string)
	port := data.Get("port").(int)
	cloudProviders := data.Get("cloud_providers").(string)

	if err := client.CreateApplication(application, email, description, port, cloudProviders); err != nil {
		return err
	}

	data.SetId(application)

	return nil
	// return resourceApplicationRead(data, meta)
}

func resourceApplicationRead(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client

	applicationName := data.Id()

	var app applicationRead
	if err := client.GetApplication(applicationName, &app); err != nil {
		return err
	}

	data.Set("email", app.Attributes.Email)
	data.Set("description", app.Attributes.Description)
	data.Set("cloud_providers", app.Attributes.CloudProviders)
	data.Set("port", app.Attributes.InstancePort)
	data.Set("accounts", app.Attributes.Accounts)

	data.SetId(app.Name)

	return nil
}

func resourceApplicationUpdate(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client

	applicationName := data.Id()

	if data.HasChanges("email", "description", "port") {
		var app applicationRead
		if err := client.GetApplication(applicationName, &app); err != nil {
			return err
		}

		email := data.Get("email").(string)
		description := data.Get("description").(string)
		port := data.Get("port").(int)
		cloudProviders := app.Attributes.CloudProviders

		if err := client.CreateApplication(applicationName, email, description, port, cloudProviders); err != nil {
			return err
		}
	}

	data.SetId(applicationName)

	return resourceApplicationRead(data, meta)
}

func resourceApplicationDelete(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client

	applicationName := data.Id()

	return client.DeleteAppliation(applicationName)
}

func resourceApplicationExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client

	applicationName := data.Id()

	var app applicationRead

	if err := client.GetApplication(applicationName, &app); err != nil {
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
