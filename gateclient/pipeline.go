package gateclient

import (
	"fmt"
	"net/http"

	"github.com/mitchellh/mapstructure"
)

func (m *GatewayClient) CreatePipeline(pipeline interface{}) error {
	resp, err := m.PipelineControllerApi.SavePipelineUsingPOST(m.Context, pipeline, nil)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Encountered an error saving pipeline, status code: %d\n", resp.StatusCode)
	}

	return nil
}

func (m *GatewayClient) GetPipeline(applicationName, pipelineName string, dest interface{}) (map[string]interface{}, error) {
	jsonMap, resp, err := m.ApplicationControllerApi.GetPipelineConfigUsingGET(m.Context,
		applicationName,
		pipelineName)

	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return jsonMap, fmt.Errorf("%s", ErrCodeNoSuchEntityException)
		}
		return jsonMap, fmt.Errorf("Encountered an error getting pipeline %s, %s\n",
			pipelineName,
			err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return jsonMap, fmt.Errorf("Encountered an error getting pipeline in pipeline %s with name %s, status code: %d\n",
			applicationName,
			pipelineName,
			resp.StatusCode)
	}

	if jsonMap == nil {
		return jsonMap, fmt.Errorf(ErrCodeNoSuchEntityException)
	}

	if err := mapstructure.Decode(jsonMap, dest); err != nil {
		return jsonMap, err
	}

	return jsonMap, nil
}

func (m *GatewayClient) UpdatePipeline(pipelineID string, pipeline interface{}) error {
	_, resp, err := m.PipelineControllerApi.UpdatePipelineUsingPUT(m.Context, pipelineID, pipeline)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Encountered an error saving pipeline, status code: %d\n", resp.StatusCode)
	}

	return nil
}

func (m *GatewayClient) DeletePipeline(applicationName, pipelineName string) error {
	resp, err := m.PipelineControllerApi.DeletePipelineUsingDELETE(m.Context, applicationName, pipelineName)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Encountered an error deleting pipeline, status code: %d\n", resp.StatusCode)
	}

	return nil
}
