package model

type Config struct {
	NodeName          string `yaml:"nodeName"`
	WorkflowDirPath   string `yaml:"workflowDirPath"`
	DatasetDirPath    string `yaml:"datasetDirPath"`
	JobServerURL      string `yaml:"jobServerURL"`
	HeartbeatInterval uint   `yaml:"heartbeatInterval"`
	JobFetchInterval  uint   `yaml:"jobFetchInterval"`
}
