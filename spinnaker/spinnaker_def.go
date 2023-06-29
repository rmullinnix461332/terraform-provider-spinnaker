package spinnaker

type applicationRead struct {
	Name       string                `json:"name"`
	Attributes applicationAttributes `json:"attributes"`
}

type applicationAttributes struct {
	Description    string `json:"description"`
	Email          string `json:"email"`
	Accounts       string `json:"accounts"`
	CloudProviders string `json:"cloudProviders"`
	InstancePort   int    `json:"instancePort"`
}

type pipelineRead struct {
	Name        string `json:"name"`
	Application string `json:"application"`
	ID          string `json:"id"`
}

type PipelineConfig struct {
	ID          string                   `json:"id,omitempty"`
	Schema      string                   `json:"schema,omitempty"`
	Type        string                   `json:"type,omitempty"`
	Name        string                   `json:"name"`
	Application string                   `json:"application"`
	Description string                   `json:"description,omitempty"`
	Parameters  []map[string]interface{} `json:"parameterConfig,omitempty"`
	Variables   map[string]interface{}   `json:"variables,omitempty"`
	Template    map[string]interface{}   `json:"template,omitempty"`
}
