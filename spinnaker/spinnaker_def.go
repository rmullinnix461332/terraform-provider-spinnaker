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
	InstancePort   string `json:"instancePort"`
}

type pipelineRead struct {
	Name        string `json:"name"`
	Application string `json:"application"`
	ID          string `json:"id"`
}

type PipelineConfig struct {
	ID                   string                   `json:"id,omitempty"`
	Schema               string                   `json:"schema,omitempty"`
	Type                 string                   `json:"type,omitempty"`
	Name                 string                   `json:"name"`
	Application          string                   `json:"application"`
	Description          string                   `json:"description,omitempty"`
	ExecutionEngine      string                   `json:"executionEngine,omitempty"`
	Parallel             bool                     `json:"parallel"`
	LimitConcurrent      bool                     `json:"limitConcurrent"`
	KeepWaitingPipelines bool                     `json:"keepWaitingPipelines"`
	Stages               []map[string]interface{} `json:"stages,omitempty"`
	Triggers             []map[string]interface{} `json:"triggers,omitempty"`
	ExpectedArtifacts    []map[string]interface{} `json:"expectedArtifacts,omitempty"`
	Parameters           []map[string]interface{} `json:"parameterConfig,omitempty"`
	Notifications        []map[string]interface{} `json:"notifications,omitempty"`
	Variables            map[string]interface{}   `json:"variables,omitempty"`
	Template             map[string]interface{}   `json:"template,omitempty"`
	LastModifiedBy       string                   `json:"lastModifiedBy"`
	Config               interface{}              `json:"config,omitempty"`
	UpdateTs             string                   `json:"updateTs"`
}
