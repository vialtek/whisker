package whisker

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/vialtek/whisker/model"
	"github.com/vialtek/whisker/remote"
)

type NodeState struct {
	NodeName   string
	Busy       bool
	CurrentJob *model.Job
	Workflows  []*model.Workflow
	Datasets   []*model.Dataset
	Recipes    []*model.Recipe
}

func NewNode() *NodeState {
	return &NodeState{
		NodeName:   GetConfig().NodeName,
		Busy:       false,
		CurrentJob: nil,
		Workflows:  loadWorkflows(),
		Datasets:   loadDatasets(),
		Recipes:    loadRecipes(),
	}
}

var clientInstance *remote.Client

func GetClient() *remote.Client {
	if clientInstance == nil {
		clientInstance = remote.NewClient(GetConfig().JobServerURL)
	}

	return clientInstance
}

func (s *NodeState) Init() {
	log.Println("Initializing Whisker...")

	GetClient().SendHeartbeat(s.Status())
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
			GetClient().SendHeartbeat(s.Status())
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

	GetClient().ChangeJobState(job.Guid, "accept")
}

func (s *NodeState) releaseCurrentJob(result *Result) {
	GetClient().SendJobOutput(s.CurrentJob.Guid, result.Output)

	if result.Success {
		GetClient().ChangeJobState(s.CurrentJob.Guid, "finished")
	} else {
		GetClient().ChangeJobState(s.CurrentJob.Guid, "failed")
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
