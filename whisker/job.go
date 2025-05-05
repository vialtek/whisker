package whisker

import (
	"log"
	"time"

	"github.com/vialtek/whisker/utils"
)

type Job struct {
	Guid     string `json:"guid"`
	Dataset  string `json:"dataset"`
	Workflow string `json:"workflow"`
}

type Result struct {
	Output    []string
	Error     string
	Success   bool
	StartedAt time.Time
	EndedAt   time.Time
}

func ProcessJob(job *Job, workflow *Workflow) *Result {
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
