package consistenthash

import (
	"hash/crc32"
	"sort"
)

// HashFunc defines function to generate hash code
type HashFunc func(data []byte) uint32

// NodeMap stores nodes and you can pick node from NodeMap
type NodeMap struct {
	hashFunc    HashFunc
	nodeHashes  []int // sorted
	nodeHashMap map[int]string
}

// NewNodeMap creates a new NodeMap
func NewNodeMap(fn HashFunc) *NodeMap {
	m := &NodeMap{
		hashFunc:    fn,
		nodeHashMap: make(map[int]string),
	}
	if m.hashFunc == nil {
		m.hashFunc = crc32.ChecksumIEEE
	}
	return m
}

// IsEmpty returns if there is no node in NodeMap
func (m *NodeMap) IsEmpty() bool {
	return len(m.nodeHashes) == 0
}

// AddNode add the given nodes into consistent hash circle
func (m *NodeMap) AddNode(keys ...string) {
	for _, key := range keys {
		if key == "" {
			continue
		}
		hash := int(m.hashFunc([]byte(key)))
		m.nodeHashes = append(m.nodeHashes, hash)
		m.nodeHashMap[hash] = key
	}
	sort.Ints(m.nodeHashes)
}

// PickNode gets the closest item in the hash to the provided key.
func (m *NodeMap) PickNode(key string) string {
	if m.IsEmpty() {
		return ""
	}

	hash := int(m.hashFunc([]byte(key)))

	// binary search for appropriate replica
	idx := sort.Search(len(m.nodeHashes), func(i int) bool {
		return hash <= m.nodeHashes[i]
	})

	// cycled back to the first replica.
	if idx == len(m.nodeHashes) {
		idx = 0
	}

	return m.nodeHashMap[m.nodeHashes[idx]]
}
