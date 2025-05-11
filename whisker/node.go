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

func (s *NodeState) Run() {
	log.Println("Whisker is running!")

	remote.SendHeartbeat(s.Status())
	heartbeatTicker := time.NewTicker(1 * time.Minute)

	for {
		select {
		case <-heartbeatTicker.C:
			remote.SendHeartbeat(s.Status())
		}
	}
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
