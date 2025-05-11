package whisker

import (
	"log"
	"time"

	"github.com/vialtek/whisker/model"
	"github.com/vialtek/whisker/remote"
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
	if s.Busy {
		return nil
	}

	client := remote.NewClient(GetConfig().JobServerURL)
	jobs := client.AvailableJobs()

	for _, job := range jobs {
		// TODO: we should also check the dataset
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

	for _, step := range workflow.Steps {
		tokens := utils.TokenizeStep(step)

		if tokens[0] == "echo" {
			result.Output = append(result.Output, tokens[1])
		} else if tokens[0] == "run" {
			// TODO: implement run command
			result.Output = append(result.Output, "TBD")
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
