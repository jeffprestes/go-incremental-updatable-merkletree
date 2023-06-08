package goincrementalupdatablemerkletree

import (
	"fmt"
	"log"
	"math"
	"math/big"

	"github.com/iden3/go-iden3-crypto/poseidon"
)

const TREE_MAX_DEPTH = 32

type LeafData struct {
	LeafID      int64
	LeafRawData []byte
	LeafHash    *big.Int
}

type IncrementalAndUpdatableMerkletree struct {
	Depth          uint
	Root           *big.Int
	NumberOfLeaves int
	Zeros          map[uint]*big.Int
	LastSubtrees   map[uint][2]*big.Int
	BaseItems      map[uint]*big.Int
}

func NewIncrementalAndUpdatableMerkletree(depth uint, zero *big.Int) (tree IncrementalAndUpdatableMerkletree, err error) {
	err = nil
	if depth > TREE_MAX_DEPTH {
		err = fmt.Errorf("NewIncrementalAndUpdatableMerkletree: tree depth must be between 1 and 32")
		return
	}
	tree.Depth = depth
	tree.Zeros = make(map[uint]*big.Int)
	tree.BaseItems = make(map[uint]*big.Int)
	tree.LastSubtrees = make(map[uint][2]*big.Int)
	var i uint
	for ; i < tree.Depth; i++ {
		tree.Zeros[i] = zero
		var hashInputs []*big.Int
		hashInputs = append(hashInputs, zero, zero)
		zero, err = poseidon.Hash(hashInputs)
		if err != nil {
			return tree, err
		}
	}
	tree.Root = zero
	return
}

func (tree *IncrementalAndUpdatableMerkletree) InsertLeaf(leaf LeafData) (err error) {
	if tree.NumberOfLeaves == TREE_MAX_DEPTH {
		err = fmt.Errorf("tree reached max leaves")
		return
	}
	index := tree.NumberOfLeaves
	hash := leaf.LeafHash
	var i uint
	for ; i < tree.Depth; i++ {
		var hashInputs []*big.Int
		if index&1 == 0 {
			// LEFT
			tree.LastSubtrees[i] = [2]*big.Int{hash, tree.Zeros[i]}
			hashInputs = append(hashInputs, hash, tree.Zeros[i])
		} else {
			// RIGHT
			tree.LastSubtrees[i] = [2]*big.Int{tree.LastSubtrees[i][0], hash}
			hashInputs = append(hashInputs, tree.LastSubtrees[i][0], hash)
		}
		hash, err = poseidon.Hash(hashInputs)
		if err != nil {
			return
		}
		index >>= 1
	}
	tree.BaseItems[uint(tree.NumberOfLeaves)] = leaf.LeafHash
	tree.NumberOfLeaves++
	tree.Root = hash
	return
}

func (tree *IncrementalAndUpdatableMerkletree) UpdateLeaf(index int64, leaf LeafData, logDebug bool) (err error) {
	if tree.NumberOfLeaves == TREE_MAX_DEPTH {
		err = fmt.Errorf("tree reached max leaves")
		return
	}
	pos := uint(index)
	hash := leaf.LeafHash
	tree.BaseItems[pos] = leaf.LeafHash
	allLeaves, err := tree.BuildTree()
	if err != nil {
		log.Printf("UpdateLeaf - tree.BuildTree - Error: %s\n", err.Error())
		return
	}
	hashInputs := []*big.Int{allLeaves[tree.Depth-1][0], allLeaves[tree.Depth-1][1]}
	tree.Root, err = poseidon.Hash(hashInputs)
	if err != nil {
		log.Printf("UpdateLeaf - tree.Root - poseidon.Hash - Error: %s\n", err.Error())
		return
	}
	if logDebug {
		log.Println("tree *IncrementalBinaryTree - UpdateLeaf - end - index:", index, " - pos:", pos, " - leaf:", leaf.LeafHash.String(), " - root:", hash.String())
	}
	return
}

func (tree *IncrementalAndUpdatableMerkletree) GenerateMerkleProofPath(index uint, logDebug bool) (path []uint8, siblings []*big.Int, err error) {
	if index > uint(len(tree.LastSubtrees)-1) {
		err = fmt.Errorf("invalid leaf index. it is greather than actual leaves inserted")
		return
	}
	var i, z, pos uint
	allLeaves, err := tree.BuildTree()
	if err != nil {
		log.Printf("GenerateMerkleProofPath - tree.BuildTree - Error: %s\n", err.Error())
		return nil, nil, err
	}
	if logDebug {
		log.Println("# GenerateMerkleProofPath", "Actual Leaves:")
		for i = 0; int(i) < len(allLeaves); i++ {
			for z = 0; int(z) < len(allLeaves[i]); z++ {
				item := allLeaves[i][z]
				log.Println("#  GenerateMerkleProofPath", "- level: ", i, " - leaf ", z, ": ", item.String())
			}
		}
		hashInputs := []*big.Int{allLeaves[3][0], allLeaves[3][1]}
		item, err := poseidon.Hash(hashInputs)
		if err != nil {
			log.Printf("GenerateMerkleProofPath - Actual Root Leaves - poseidon.Hash - Error: %s\n", err.Error())
			return nil, nil, err
		}
		log.Println("# GenerateMerkleProofPath", "Actual Root Leaves:", item.String())
		log.Println("# GenerateMerkleProofPath", "Actual Basetree:")
		for i = 0; int(i) < len(tree.BaseItems); i++ {
			item := tree.BaseItems[uint(i)]
			log.Println("#  GenerateMerkleProofPath", "item", item.String())
		}
	}
	pos = uint(index)
	for i = 0; i < tree.Depth; i++ {
		if index&1 == 0 {
			// LEFT
			path = append(path, 0)
			siblings = append(siblings, allLeaves[i][pos+1])
		} else {
			// RIGHT
			path = append(path, 1)
			siblings = append(siblings, allLeaves[i][pos-1])
		}
		index >>= 1
		pos = uint(math.Pow(float64(index), 1))
		if logDebug {
			log.Println("# GenerateMerkleProofPath", " | Path: ", path[len(path)-1], " | Sibling: ", siblings[len(siblings)-1], " | Index: ", index, " | Pos: ", pos)
		}
	}
	if logDebug {
		log.Println("# GenerateMerkleProofPath", "Actual Root:", tree.Root)
	}
	return
}

func (tree *IncrementalAndUpdatableMerkletree) FindIndexByValue(hashLeaf *big.Int) (index int64, err error) {
	for idx, item := range tree.BaseItems {
		if item.Cmp(hashLeaf) == 0 {
			index = int64(idx)
			return
		}
	}
	err = fmt.Errorf("index not found")
	return
}

func (tree *IncrementalAndUpdatableMerkletree) BuildTree() (allLeaves map[uint]map[uint]*big.Int, err error) {
	var i, z, numLeaves, treeDepth uint
	var item *big.Int
	allLeaves = make(map[uint]map[uint]*big.Int)
	treeDepth = uint(tree.Depth)
	for i = 0; i < uint(tree.Depth); i++ {
		allLeaves[i] = make(map[uint]*big.Int)
		log.Println(" === BuildTree === ")
		log.Println(float64(treeDepth - i))
		log.Println(math.Pow(2, float64(treeDepth-i)))
		log.Println(uint(math.Pow(2, float64(treeDepth-i))))
		numLeaves = uint(math.Pow(2, float64(treeDepth-i)))
		for z = 0; z < numLeaves; z++ {
			if i < 1 {
				if tree.BaseItems[uint(z)] != nil {
					item = tree.BaseItems[uint(z)]
				} else {
					item = tree.Zeros[i]
				}
			} else {
				hashInputs := []*big.Int{allLeaves[i-1][z*2], allLeaves[i-1][((z * 2) + 1)]}
				item, err = poseidon.Hash(hashInputs)
				if err != nil {
					log.Printf("BuildTree - poseidon.Hash - Error: %s\n", err.Error())
					return nil, err
				}
			}
			allLeaves[i][z] = item
		}
	}
	return
}

func (tree *IncrementalAndUpdatableMerkletree) PopulateTreeWithZeros(debug bool) (err error) {
	var i, z, numLeaves, treeDepth uint
	var hash *big.Int
	treeDepth = uint(tree.Depth)
	if debug {
		log.Println(" === PopulateTreeWithZeros === ")
	}
	for i = 0; i < uint(tree.Depth); i++ {
		if debug {
			log.Println("Layer: ", float64(treeDepth-i))
			log.Println("Number of leaves: ", uint(math.Pow(2, float64(treeDepth-i))))
		}
		numLeaves = uint(math.Pow(2, float64(treeDepth-i)))
		index := 0
		if i < 1 {
			hash = tree.Zeros[i]
		}
		for z = 0; z < numLeaves; z++ {
			if i < 1 {
				tree.BaseItems[z] = tree.Zeros[i]
			} else {
				var hashInputs []*big.Int
				if index&1 == 0 {
					// LEFT
					tree.LastSubtrees[i] = [2]*big.Int{hash, tree.Zeros[i]}
					hashInputs = append(hashInputs, hash, tree.Zeros[i])
				} else {
					// RIGHT
					tree.LastSubtrees[i] = [2]*big.Int{tree.LastSubtrees[i][0], hash}
					hashInputs = append(hashInputs, tree.LastSubtrees[i][0], hash)
				}
				hash, err = poseidon.Hash(hashInputs)
				if err != nil {
					return
				}
				index >>= 1
				tree.Root = hash
			}
		}
	}

	if debug {
		log.Println(" === End PopulateTreeWithZeros === ")
	}
	return
}

/*
index := tree.NumberOfLeaves
	hash := leaf.LeafHash
	var i uint
	for ; i < tree.Depth; i++ {
		var hashInputs []*big.Int
		if index&1 == 0 {
			// LEFT
			tree.LastSubtrees[i] = [2]*big.Int{hash, tree.Zeros[i]}
			hashInputs = append(hashInputs, hash, tree.Zeros[i])
		} else {
			// RIGHT
			tree.LastSubtrees[i] = [2]*big.Int{tree.LastSubtrees[i][0], hash}
			hashInputs = append(hashInputs, tree.LastSubtrees[i][0], hash)
		}
		hash, err = poseidon.Hash(hashInputs)
		if err != nil {
			return
		}
		index >>= 1
	}
	tree.BaseItems[uint(tree.NumberOfLeaves)] = leaf.LeafHash
	tree.NumberOfLeaves++
	tree.Root = hash
	return
*/
