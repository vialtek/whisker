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
	NodeName  string
	Busy      bool
	Workflows []*model.Workflow
	Datasets  []*model.Dataset
}

func NewNode() *NodeState {
	return &NodeState{
		NodeName:  GetConfig().NodeName,
		Busy:      false,
		Workflows: loadWorkflows(),
		Datasets:  loadDatasets(),
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
	s.Busy = true
	log.Println("Starting job:", job.Guid)

	client := remote.NewClient(GetConfig().JobServerURL)
	client.AcceptJob(job.Guid)

	result := s.executeJob(job)
	client.SendJobOutput(job.Guid, result.Output)

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

	return result
}
