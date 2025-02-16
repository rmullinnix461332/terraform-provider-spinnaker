package spinnaker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/ghodss/yaml"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rmullinnix461332/terraform-provider-spinnaker/gateclient"
)

func resourcePipelineTemplate() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateTemplateName,
			},
			"template": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppressEquivalentPipelineTemplateDiffs,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Create: resourcePipelineTemplateCreate,
		Read:   resourcePipelineTemplateRead,
		Update: resourcePipelineTemplateUpdate,
		Delete: resourcePipelineTemplateDelete,
		Exists: resourcePipelineTemplateExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourcePipelineTemplateCreate(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client

	template := data.Get("template").(string)
	templateName := data.Get("name").(string)

	d, err := yaml.YAMLToJSON([]byte(template))
	if err != nil {
		return err
	}

	var jsonContent map[string]interface{}
	if err = json.NewDecoder(bytes.NewReader(d)).Decode(&jsonContent); err != nil {
		return fmt.Errorf("Error decoding json: %s", err.Error())
	}

	if _, ok := jsonContent["schema"]; !ok {
		return fmt.Errorf("Pipeline save command currently only supports pipeline template configurations")
	}

	jsonContent["id"] = templateName
	//templateName := jsonContent["id"].(string)

	log.Println("[DEBUG] Making request to spinnaker")
	if err := client.CreatePipelineTemplate(jsonContent); err != nil {
		log.Printf("[DEBUG] Error response from spinnaker: %s", err.Error())
		return err
	}

	log.Printf("[DEBUG] Created template successfully")

	data.SetId(templateName)

	return nil
}

func resourcePipelineTemplateRead(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client

	templateName := data.Id()

	t := make(map[string]interface{})
	if err := client.GetPipelineTemplate(templateName, &t); err != nil {
		if err.Error() == gateclient.ErrCodeNoSuchEntityException {
			data.SetId("")
			return nil
		}
		return err
	}

	// Remove timestamp from response
	delete(t, "updateTs")
	delete(t, "lastModifiedBy")

	jsonContent, err := json.Marshal(t)
	if err != nil {
		return err
	}

	raw, err := yaml.JSONToYAML(jsonContent)
	if err != nil {
		return err
	}
	data.Set("name", templateName)
	data.Set("template", string(raw))
	data.Set("url", fmt.Sprintf("spinnaker://%s", t["id"].(string)))
	data.SetId(templateName)

	return nil
}

func resourcePipelineTemplateUpdate(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	var templateName string
	template := data.Get("template").(string)

	d, err := yaml.YAMLToJSON([]byte(template))
	if err != nil {
		return err
	}

	var jsonContent map[string]interface{}
	if err = json.NewDecoder(bytes.NewReader(d)).Decode(&jsonContent); err != nil {
		return fmt.Errorf("Error decoding json: %s", err.Error())
	}

	if _, ok := jsonContent["schema"]; !ok {
		return fmt.Errorf("Pipeline save command currently only supports pipeline template configurations")
	}

	templateName = jsonContent["id"].(string)

	if err := client.UpdatePipelineTemplate(templateName, jsonContent); err != nil {
		return err
	}

	data.SetId(templateName)
	return resourcePipelineTemplateRead(data, meta)
}

func resourcePipelineTemplateDelete(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	templateName := data.Id()

	if err := client.DeletePipelineTemplate(templateName); err != nil {
		return err
	}

	data.SetId("")
	return nil
}

func resourcePipelineTemplateExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	templateName := data.Id()

	t := &templateRead{}
	if err := client.GetPipelineTemplate(templateName, t); err != nil {
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

func suppressEquivalentPipelineTemplateDiffs(k, old, new string, d *schema.ResourceData) bool {
	equivalent, err := areEqualJSON(old, new)

	if err != nil {
		return false
	}

	return equivalent
}

func areEqualJSON(s1, s2 string) (bool, error) {
	var o1 interface{}
	var o2 interface{}

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

	var x1 interface{}
	bytesA, _ := json.Marshal(o1)
	_ = json.Unmarshal(bytesA, &x1)
	var x2 interface{}
	bytesB, _ := json.Marshal(o2)
	_ = json.Unmarshal(bytesB, &x2)

	return reflect.DeepEqual(x1, x2), nil
}
