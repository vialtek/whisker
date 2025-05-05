package whisker

import (
	"log"
)

type NodeState struct {
	NodeName string
	Busy     bool
}

func NewNode() *NodeState {
	return &NodeState{
		NodeName: "Node",
		Busy:     false,
	}
}

func (s *NodeState) Run() {
	log.Println("Whisker is running!")

	for {
	}
}
