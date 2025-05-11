package model

type Workflow struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Workflow    string   `yaml:"workflow"`
	Steps       []string `yaml:"steps"`
}
