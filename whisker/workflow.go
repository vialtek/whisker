package whisker

import (
	"gopkg.in/yaml.v3"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/vialtek/whisker/model"
)

func loadWorkflows() []*model.Workflow {
	var workflows []*model.Workflow
	log.Println("Loading workflows...")

	for _, path := range findInDirectory(GetConfig().WorkflowDirPath, ".yaml") {
		workflows = append(workflows, parseWorkflowYaml(path))
	}

	log.Println("Workflows loaded:", len(workflows))
	if len(workflows) == 0 {
		log.Println("Warning: no workflow loaded.")
	}

	return workflows
}

func parseWorkflowYaml(path string) *model.Workflow {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		log.Printf("parseWorkflowYaml could not open file:", err)
	}

	newWorkflow := &model.Workflow{}
	err = yaml.Unmarshal(yamlFile, newWorkflow)
	if err != nil {
		log.Fatalf("parseWorkflowYaml could not unmarshal:", err)
	}

	return newWorkflow
}

func findInDirectory(root, ext string) []string {
	var files []string
	filepath.WalkDir(root, func(s string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(d.Name()) == ext {
			files = append(files, s)
		}
		return nil
	})
	return files
}
