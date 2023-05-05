package bloom

import (
	"fmt"
	"testing"
)

func nodeBeforeEach() (node *Node) {
	input := NodeInitInput{HashPrefix: []byte("some_prefix"), MinEntryCount: 1000000000}
	node = InitNode(input)
	return
}

func TestInitNode(t *testing.T) {
	node := nodeBeforeEach()
	if node.hashPrefix == nil {
		t.Errorf("expected hash to be initialized")
	}
	if node.lock == nil {
		t.Errorf("expected lock to be initialized")
	}
	if node.storage == nil {
		t.Errorf("expected storage to be initialized")
	}
}

func TestNode_ItemAdd_Simple(t *testing.T) {
	node := nodeBeforeEach()

	nodeVal1 := []byte("node1")
	contains, err := node.ItemPossiblyContains(nodeVal1)
	if err != nil {
		t.Error(err)
	}
	if contains {
		t.Errorf("expected to not contain the item")
	}

	err = node.ItemAdd(nodeVal1)
	if err != nil {
		t.Error(err)
	}

	contains, err = node.ItemPossiblyContains(nodeVal1)
	if err != nil {
		t.Error(err)
	}
	if !contains {
		t.Errorf("expected to contain the item")
	}
}

func TestNode_ItemAdd(t *testing.T) {
	node := nodeBeforeEach()

	for idx := 0; idx < 10000000; idx++ {
		node2Hash := []byte(fmt.Sprintf("node_%d", idx))

		err := node.ItemAdd(node2Hash)
		if err != nil {
			t.Error(err)
		}

		contains, err := node.ItemPossiblyContains(node2Hash)
		if err != nil {
			t.Error(err)
		}
		if !contains {
			t.Errorf("expected to not contain the item")
		}

		if idx%100000 == 0 {
			t.Logf("added %d items", idx)
		}
	}
}
