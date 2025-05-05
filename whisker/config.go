package whisker

type Config struct {
	NodeName         string `yaml:"nodeName"`
	WorkflowDirPath  string `yaml:"workflowDirPath"`
}

var configInstance *Config

// TODO: parse the config.yaml file
func GetConfig() *Config {
	if configInstance == nil {
		configInstance = defaultConfig()
	}

	return configInstance
}

func defaultConfig() *Config {
	return &Config{
		WorkflowDirPath:  "./examples/workflows",
	}
}


