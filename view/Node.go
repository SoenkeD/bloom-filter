package view

import (
	"bloom/bloom"
	"log"
)

type NodeInput struct {
	Id       string
	Setup    func(bloomNode *bloom.Node) error
	Shutdown func() error
}

type Node struct {
	id       string
	setup    func(bloomNode *bloom.Node) error
	shutdown func() error
}

func initNode(input *NodeInput) (node *Node) {
	node = &Node{
		id:       input.Id,
		setup:    input.Setup,
		shutdown: input.Shutdown,
	}
	return
}

type NodeManager struct {
	bloomNode *bloom.Node
	viewNodes []*Node
}

func (manager *NodeManager) Setup() (err error) {

	for _, node := range manager.viewNodes {
		go func(node *Node) {
			setupErr := node.setup(manager.bloomNode)
			if setupErr != nil {
				log.Printf("failed to setup view=%s with err=%s", node.id, setupErr)
			}
		}(node)

	}

	return
}

func (manager *NodeManager) Shutdown() (err error) {

	for _, node := range manager.viewNodes {
		err = node.shutdown()
		if err != nil {
			return
		}
	}

	return
}

func CreateNodes(inputs []*NodeInput, nodeInit bloom.NodeInitInput) (manager *NodeManager) {
	bloomNode := bloom.InitNode(nodeInit)

	var nodes []*Node
	for _, input := range inputs {
		nodes = append(nodes, initNode(input))
	}

	manager = &NodeManager{
		bloomNode: bloomNode,
		viewNodes: nodes,
	}

	return
}
