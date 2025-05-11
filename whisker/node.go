package whisker

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/vialtek/whisker/remote"
)

type NodeState struct {
	NodeName  string
	Busy      bool
	Workflows []*Workflow
	Datasets  []*Dataset
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

	remote.SendHeartbeat(nodeStatus(s))
	heartbeatTicker := time.NewTicker(1 * time.Minute)

	for {
		select {
		case <-heartbeatTicker.C:
			remote.SendHeartbeat(nodeStatus(s))
		}
	}
}

func nodeStatus(state *NodeState) map[string]string {
	result := make(map[string]string)

	result["node_name"] = state.NodeName
	result["busy"] = strconv.FormatBool(state.Busy)

	var workflowNames []string
	for _, wf := range state.Workflows {
		if wf != nil {
			workflowNames = append(workflowNames, wf.Workflow)
		}
	}
	result["workflows"] = strings.Join(workflowNames, ",")

	var datasetNames []string
	for _, ds := range state.Datasets {
		if ds != nil {
			datasetNames = append(datasetNames, ds.Name)
		}
	}
	result["datasets"] = strings.Join(datasetNames, ",")

	return result
}
