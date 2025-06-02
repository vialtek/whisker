package whisker

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/vialtek/whisker/model"
	"github.com/vialtek/whisker/utils"
)

type Result struct {
	Output    []string
	Error     string
	Success   bool
	StartedAt time.Time
	EndedAt   time.Time
}

func (s *NodeState) pickNewJob() *model.Job {
	jobs := GetClient().AvailableJobs()

	if len(jobs) > 0 {
		log.Println("Picking new job, jobs available:", len(jobs))
	}

	// Finds first job that require workflow and dataset available on our node
	for _, job := range jobs {
		if s.workflowByName(job.Workflow) != nil && s.datasetByName(job.Dataset) != nil {
			return job
		}
	}

	return nil
}

func (s *NodeState) executeJob(job *model.Job) *Result {
	result := &Result{StartedAt: time.Now()}

	workflow := s.workflowByName(job.Workflow)
	if workflow == nil {
		errMsg := "Error: workflow not found."

		result.Error = errMsg
		result.Success = false
		result.EndedAt = time.Now()

		log.Println(errMsg)
		return result
	}

	err := prepareJobDir(job)
	if err != nil {
		result.Error = err.Error()
		result.Success = false
		result.EndedAt = time.Now()

		log.Println(result.Error)
		return result
	}

	for _, step := range workflow.Steps {
		tokens := utils.TokenizeStep(step)

		if tokens[0] == "echo" {
			result.Output = append(result.Output, tokens[1])
		} else if tokens[0] == "run" {
			recipe := s.recipeByName(tokens[1])

			if recipe == nil {
				errMsg := "Error: recipe not found."

				result.Error = errMsg
				result.Success = false
				result.EndedAt = time.Now()

				log.Println(errMsg)
				return result
			}

			// TODO: should report failure
			execRecipe(recipe, result)
		} else {
			errMsg := "Error: unsupported action in step: " + step

			result.Error = errMsg
			result.Success = false
			result.EndedAt = time.Now()

			log.Println(errMsg)
			return result
		}
	}

	result.Success = true
	result.EndedAt = time.Now()

	return result
}

func prepareJobDir(job *model.Job) error {
	jobDir := filepath.Join(GetConfig().JobsDirPath, job.Guid)

	// Directory already exists, remove it
	if _, err := os.Stat(jobDir); err == nil {
		if err := os.RemoveAll(jobDir); err != nil {
			return fmt.Errorf("failed to remove directory %s: %w", jobDir, err)
		}
	}

	if err := os.MkdirAll(jobDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", jobDir, err)
	}

	// Exporting info to the directory
	if job.Params != nil {
		paramsData, _ := json.Marshal(job.Params)
		fullPath := filepath.Join(jobDir, "params.json")
		if err := os.WriteFile(fullPath, paramsData, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", fullPath, err)
		}
	}

	return nil
}
