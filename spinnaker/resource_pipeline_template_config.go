package spinnaker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
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
				DiffSuppressFunc: suppressEquivalentPipelineConfigDiffs,
			},
			"application": {
				Type:     schema.TypeString,
				Required: true,
			},
			"template_name": {
				Type:     schema.TypeString,
				Required: true,
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

	return resourcePipelineTemplateConfigRead(data, meta)
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

	jsonContent, err := json.Marshal(p)
	if err != nil {
		return err
	}

	raw, err := yaml.JSONToYAML(jsonContent)
	if err != nil {
		return err
	}

	data.Set("template_name", name)
	data.Set("application", application)
	data.Set("pipeline_config", string(raw))

	data.SetId(id)

	return nil
}

func resourcePipelineTemplateConfigUpdate(data *schema.ResourceData, meta interface{}) error {
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

	pConfig, err := buildConfig(data)
	if err != nil {
		return err
	}

	pConfig.ID = p.ID
	if err := client.UpdatePipeline(p.ID, *pConfig); err != nil {
		return err
	}

	jsonContent, err := json.Marshal(pConfig)
	if err != nil {
		return err
	}

	raw, err := yaml.JSONToYAML(jsonContent)
	if err != nil {
		return err
	}

	data.Set("template_name", pConfig.Name)
	data.Set("application", application)
	data.Set("pipeline_config", string(raw))

	data.SetId(application + ":" + pConfig.Name)

	return resourcePipelineTemplateConfigRead(data, meta)
}

func resourcePipelineTemplateConfigDelete(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client

	id := data.Id()

	parts := strings.Split(id, ":")
	application := parts[0]
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

	id := data.Id()

	parts := strings.Split(id, ":")
	application := parts[0]
	name := data.Get("template_name").(string)

	p := PipelineConfig{}
	if _, err := client.GetPipeline(application, name, &p); err != nil {
		if err.Error() == gateclient.ErrCodeNoSuchEntityException {
			return false, nil
		}
		return false, err
	}

	if p.Name == name {
		return true, nil
	}

	return false, nil
}

func buildConfig(data *schema.ResourceData) (*PipelineConfig, error) {
	config := data.Get("pipeline_config").(string)
	tName := data.Get("template_name").(string)

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
	pConfig.Name = tName

	return &pConfig, err
}

func suppressEquivalentPipelineConfigDiffs(k, old, new string, d *schema.ResourceData) bool {
	equivalent, err := areRoughlyEqualJSON(old, new)
	if err != nil {
		return false
	}

	return equivalent
}

func areRoughlyEqualJSON(s1 string, s2 string) (bool, error) {
	var o1 PipelineConfig
	var o2 PipelineConfig

	var err error
	log.Printf("[DEBUG] s1: %s", s1)
	err = yaml.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string 1 :: %s", err.Error())
	}
	log.Printf("[DEBUG] s2: %s", s2)
	err = yaml.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string 2 :: %s", err.Error())
	}

	equal := true
	if len(o1.ID) > 0 && len(o2.ID) > 0 {
		equal = o1.ID == o2.ID
	}

	equal = equal && o1.Application == o2.Application
	equal = equal && o1.Description == o2.Description
	equal = equal && reflect.DeepEqual(o1.Variables, o2.Variables)
	equal = equal && reflect.DeepEqual(o1.Template, o2.Template)

	return equal, nil
}
