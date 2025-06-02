package model

type Recipe struct {
	Name       string `yaml:"name"`
	Runner     string `yaml:"runner"`
	Entrypoint string `yaml:"entrypoint"`
	Pwd        string
}
