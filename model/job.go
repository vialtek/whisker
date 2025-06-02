package model

type Job struct {
	Guid     string                 `json:"guid"`
	Dataset  string                 `json:"dataset"`
	Workflow string                 `json:"workflow"`
	Status   string                 `json:"status"`
	Params   map[string]interface{} `json:"params"`
}
