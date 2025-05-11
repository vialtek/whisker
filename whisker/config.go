package whisker

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/vialtek/whisker/model"
)

var configInstance *model.Config

func GetConfig() *model.Config {
	if configInstance == nil {
		configInstance = parseConfig()

		if configInstance == nil {
			log.Println("Warning: config file not found, using default config attributes.")
			configInstance = defaultConfig()
		}
	}

	return configInstance
}

func parseConfig() *model.Config {
	var config *model.Config
	log.Println("Loading config file...")

	configFilePath := getConfigFilePath()
	if configFilePath == "" {
		return nil
	}

	yamlFile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatalf("Error: yamlFile Config file read: %v", err)
		return nil
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Error: yamlFile Config Unmarshal: %v", err)
		return nil
	}

	return config
}

func getConfigFilePath() string {
	ex, _ := os.Executable()
	localConfigPath := filepath.Dir(ex) + "/config.yaml"
	globalConfigPath := "/etc/whisker/config.yaml"

	_, err := os.Stat(localConfigPath)
	if err == nil {
		log.Println("Using config on: ", localConfigPath)
		return localConfigPath
	}

	_, glob_err := os.Stat(globalConfigPath)
	if glob_err == nil {
		log.Println("Using config on:", globalConfigPath)
		return globalConfigPath
	}

	// TODO: should return error when not found
	return ""
}

func defaultConfig() *model.Config {
	return &model.Config{
		NodeName:        "New Node",
		WorkflowDirPath: "./examples/workflows",
		DatasetDirPath:  "./examples/datasets",
		JobServerURL:    "http://localhost:4567",
	}
}
