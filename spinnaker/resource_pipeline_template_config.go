package spinnaker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rmullinnix461332/terraform-provider-spinnaker/gateclient"
)

func resourcePipelineTemplateConfig() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"pipeline_config": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppressEquivalentPipelineTemplateDiffs,
			},
			"application": {
				Type:     schema.TypeString,
				Required: true,
			},
			"template_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
		},
		Create: resourcePipelineTemplateConfigCreate,
		Read:   resourcePipelineTemplateConfigRead,
		Update: resourcePipelineTemplateConfigUpdate,
		Delete: resourcePipelineTemplateConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourcePipelineTemplateConfigCreate(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client

	pConfig, err := buildConfig(data)

	if err != nil {
		return err
	}

	if err := client.CreatePipeline(*pConfig); err != nil {
		fmt.Printf("[DEBUG] Error response from spinnaker: %s", err.Error())
		return err
	}

	data.SetId(pConfig.Application + ":" + pConfig.Name)

	return nil
}

func resourcePipelineTemplateConfigRead(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client

	id := data.Id()

	parts := strings.Split(id, ":")
	application := parts[0]
	name := parts[1]

	p := PipelineConfig{}
	if _, err := client.GetPipeline(application, name, &p); err != nil {
		if err.Error() == gateclient.ErrCodeNoSuchEntityException {
			data.SetId("")
			return nil
		}
		return err
	}
	/*
		raw, err := yaml.Marshal(p)
		if err != nil {
			return err
		}

		data.Set("template_name", p.Name)
		data.Set("application", p.Application)
	*/
	data.SetId(p.ID)

	return nil
}

func resourcePipelineTemplateConfigUpdate(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	pipelineID := data.Id()

	pConfig, err := buildConfig(data)
	if err != nil {
		return err
	}

	pConfig.ID = pipelineID
	if err := client.UpdatePipeline(pipelineID, *pConfig); err != nil {
		return err
	}

	return resourcePipelineTemplateConfigRead(data, meta)
}

func resourcePipelineTemplateConfigDelete(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client

	application := data.Get("application").(string)
	name := data.Get("template_name").(string)

	if err := client.DeletePipeline(application, name); err != nil {
		return err
	}

	data.SetId("")

	return nil
}

func resourcePipelineTemplateConfigExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	templateName := data.Id()

	var t templateRead
	if err := client.GetPipelineTemplate(templateName, &t); err != nil {
		if err.Error() == gateclient.ErrCodeNoSuchEntityException {
			return false, nil
		}
		return false, err
	}

	if t.ID == templateName {
		return true, nil
	}

	return false, nil
}

func buildConfig(data *schema.ResourceData) (*PipelineConfig, error) {
	config := data.Get("pipeline_config").(string)

	d, err := yaml.YAMLToJSON([]byte(config))
	if err != nil {
		return nil, err
	}

	var jsonContent map[string]interface{}
	if err = json.NewDecoder(bytes.NewReader(d)).Decode(&jsonContent); err != nil {
		return nil, fmt.Errorf("Error decoding json: %s", err.Error())
	}

	_, ok := jsonContent["name"].(string)
	if !ok {
		return nil, fmt.Errorf("pipeline name not set in configuration")
	}

	_, ok = jsonContent["application"].(string)
	if !ok {
		return nil, fmt.Errorf("application not set in pipeline configuration")
	}

	var pConfig PipelineConfig
	json.Unmarshal(d, &pConfig)
	pConfig.Type = "templatedPipeline"

	return &pConfig, err
}
