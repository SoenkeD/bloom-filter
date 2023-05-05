package bloom

import (
	"hash/fnv"
	"sync"
)

type NodeInitInput struct {
	HashPrefix    []byte
	MinEntryCount uint64
}

type Node struct {
	hashPrefix []byte
	lock       *sync.Mutex
	storage    *Storage
}

func InitNode(input NodeInitInput) (node *Node) {

	node = &Node{
		hashPrefix: input.HashPrefix,
		lock:       &sync.Mutex{},
		storage:    InitStorage(StorageInitInput{MinEntryCount: input.MinEntryCount}),
	}

	return
}

func (node *Node) nodeHash(hashInput []byte) (hash uint64, err error) {
	algo := fnv.New64a()
	_, err = algo.Write(append(node.hashPrefix, hashInput...))
	if err != nil {
		return
	}
	hash = algo.Sum64() % uint64(node.storage.GetSlotCount())

	return
}

func (node *Node) ItemAdd(item []byte) (err error) {
	node.lock.Lock()
	defer node.lock.Unlock()

	hash, err := node.nodeHash(item)
	if err != nil {
		return
	}

	err = node.storage.AddItem(hash)

	return
}

func (node *Node) ItemPossiblyContains(item []byte) (contains bool, err error) {
	node.lock.Lock()
	defer node.lock.Unlock()

	hash, err := node.nodeHash(item)
	if err != nil {
		return
	}

	contains = node.storage.PotentiallyKnowItem(hash)

	return
}
