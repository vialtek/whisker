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

func (s *NodeState) loadJobs() {
	if s.Busy {
		return
	}

	// TODO: make singleton
	client := remote.NewClient(GetConfig().JobServerURL)
	jobs := client.AvailableJobs()

	if len(jobs) == 0 {
		return
	}

	// TODO: pick job to be executed
}

func ProcessJob(job *model.Job, workflow *model.Workflow) *Result {
	log.Println("Starting job:", job.Guid)

	result := &Result{StartedAt: time.Now()}

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
