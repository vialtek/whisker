package whisker

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/vialtek/whisker/model"
	"github.com/vialtek/whisker/remote"
)

var client *remote.Client

type NodeState struct {
	NodeName   string
	Busy       bool
	CurrentJob *model.Job
	Workflows  []*model.Workflow
	Datasets   []*model.Dataset
}

func NewNode() *NodeState {
	return &NodeState{
		NodeName:   GetConfig().NodeName,
		Busy:       false,
		CurrentJob: nil,
		Workflows:  loadWorkflows(),
		Datasets:   loadDatasets(),
	}
}

func (s *NodeState) Init() {
	log.Println("Initializing Whisker...")

	client := remote.NewClient(GetConfig().JobServerURL)
	client.SendHeartbeat(s.Status())
}

func (s *NodeState) Run() {
	log.Println("Whisker is running!")

	checkWorkTicker := time.NewTicker(time.Duration(GetConfig().JobFetchInterval) * time.Second)
	heartbeatTicker := time.NewTicker(time.Duration(GetConfig().HeartbeatInterval) * time.Second)

	for {
		select {
		case <-checkWorkTicker.C:
			s.manageWorkload()
		case <-heartbeatTicker.C:
			client.SendHeartbeat(s.Status())
		}
	}
}

func (s *NodeState) manageWorkload() {
	if s.Busy {
		return
	}

	job := s.pickNewJob()
	if job == nil {
		return
	}

	// TODO: async
	s.takeJob(job)
	result := s.executeJob(job)
	s.releaseCurrentJob(result)
}

func (s *NodeState) takeJob(job *model.Job) {
	log.Println("Starting job:", job.Guid)

	s.Busy = true
	s.CurrentJob = job

	client := remote.NewClient(GetConfig().JobServerURL)
	client.ChangeJobState(job.Guid, "accept")
}

func (s *NodeState) releaseCurrentJob(result *Result) {
	client := remote.NewClient(GetConfig().JobServerURL)
	client.SendJobOutput(s.CurrentJob.Guid, result.Output)

	if result.Success {
		client.ChangeJobState(s.CurrentJob.Guid, "finished")
	} else {
		client.ChangeJobState(s.CurrentJob.Guid, "failed")
	}

	s.CurrentJob = nil
	s.Busy = false

	elapsed := result.EndedAt.Sub(result.StartedAt)
	log.Println("Job ended. Elapsed time:", elapsed)
}

func (s *NodeState) Status() map[string]string {
	result := make(map[string]string)

	result["node_name"] = s.NodeName
	result["busy"] = strconv.FormatBool(s.Busy)

	var workflowNames []string
	for _, wf := range s.Workflows {
		if wf != nil {
			workflowNames = append(workflowNames, wf.Workflow)
		}
	}
	result["workflows"] = strings.Join(workflowNames, ",")

	var datasetNames []string
	for _, ds := range s.Datasets {
		if ds != nil {
			datasetNames = append(datasetNames, ds.Name)
		}
	}
	result["datasets"] = strings.Join(datasetNames, ",")

	if s.CurrentJob != nil {
		job := s.CurrentJob
		result["current_job_guid"] = job.Guid
		result["current_job_workflow"] = job.Workflow
		result["current_job_dataset"] = job.Dataset
	}

	return result
}
