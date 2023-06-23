package gateclient

import (
	"fmt"
	"net/http"

	"github.com/mitchellh/mapstructure"
)

const (
	ErrCodeNoSuchEntityException = "NoSuchEntityException"
)

func (m *GatewayClient) CreatePipelineTemplate(template interface{}) error {
	resp, err := m.PipelineTemplatesControllerApi.CreateUsingPOST(m.Context, template)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("Encountered an error saving template, status code: %d\n", resp.StatusCode)
	}

	return nil
}

func (m *GatewayClient) GetPipelineTemplate(templateID string, dest interface{}) error {
	successPayload, resp, err := m.PipelineTemplatesControllerApi.GetUsingGET(m.Context, templateID)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("%s", ErrCodeNoSuchEntityException)
		}
		return fmt.Errorf("Encountered an error getting pipeline template %s, %s\n",
			templateID,
			err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Encountered an error getting pipeline template %s, status code: %d\n",
			templateID,
			resp.StatusCode,
		)
	}

	if successPayload == nil {
		return fmt.Errorf(ErrCodeNoSuchEntityException)
	}

	if err := mapstructure.Decode(successPayload, dest); err != nil {
		return err
	}

	return nil
}

func (m *GatewayClient) DeletePipelineTemplate(templateID string) error {
	_, resp, err := m.PipelineTemplatesControllerApi.DeleteUsingDELETE(m.Context, templateID, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("Encountered an error deleting pipeline template %s, status code: %d\n",
			templateID,
			resp.StatusCode)
	}

	return nil
}

func (m *GatewayClient) UpdatePipelineTemplate(templateID string, template interface{}) error {
	resp, err := m.PipelineTemplatesControllerApi.UpdateUsingPOST(m.Context, templateID, template, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("Encountered an error updating pipeline template %s, status code: %d\n",
			templateID,
			resp.StatusCode)
	}

	return nil
}
