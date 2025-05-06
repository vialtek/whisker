package whisker

import (
	"log"
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

	for {
	}
}
