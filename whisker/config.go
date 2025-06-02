package whisker

import (
	"errors"
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
	}

	return configInstance
}

func parseConfig() *model.Config {
	log.Println("Loading config file...")
	config := defaultConfig()

	configFilePath, err := getConfigFilePath()
	if err != nil {
		log.Println(err)
		return config
	}

	yamlFile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Println("Error: parseConfig - yamlFile Config file read:", err)
		return config
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Println("Error: parseConfig - yamlFile Config Unmarshal:", err)
		return config
	}

	return config
}

func getConfigFilePath() (string, error) {
	localConfigPath := filepath.Join(filepath.Dir(os.Args[0]), "config.yaml")
	globalConfigPath := "/etc/whisker/config.yaml"

	for _, path := range []string{localConfigPath, globalConfigPath} {
		if _, err := os.Stat(path); err == nil {
			log.Println("Using config at:", path)
			return path, nil
		}
	}

	return "", errors.New("No config file not found, using default configuration")
}

func defaultConfig() *model.Config {
	return &model.Config{
		NodeName:          "New Node",
		WorkflowDirPath:   "./examples/workflows",
		DatasetDirPath:    "./examples/datasets",
		RecipeDirPath:     "./examples/recipes",
		JobsDirPath:       "./jobs_dir",
		JobServerURL:      "http://localhost:4567",
		HeartbeatInterval: 60,
		JobFetchInterval:  1,
	}
}
