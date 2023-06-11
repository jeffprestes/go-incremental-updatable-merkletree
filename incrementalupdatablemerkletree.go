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

func NewPopulatedIncrementalAndUpdatableMerkletree(depth uint, zeros map[uint]*big.Int) (tree IncrementalAndUpdatableMerkletree, err error) {
	err = nil
	if depth > TREE_MAX_DEPTH {
		err = fmt.Errorf("NewIncrementalAndUpdatableMerkletree: tree depth must be between 1 and 32")
		return
	}
	tree.Depth = depth
	tree.Zeros = zeros
	tree.BaseItems = make(map[uint]*big.Int)
	tree.LastSubtrees = make(map[uint][2]*big.Int)
	var i uint
	var zero *big.Int
	for ; i < tree.Depth; i++ {
		var hashInputs []*big.Int
		hashInputs = append(hashInputs, tree.Zeros[i], tree.Zeros[i])
		zero, err = poseidon.Hash(hashInputs)
		if err != nil {
			return tree, err
		}
	}
	tree.Root = zero
	return
}

func (tree *IncrementalAndUpdatableMerkletree) InsertLeaf(leaf LeafData) (err error) {
	err = tree.InsertLeafWithDebug(leaf, false)
	return
}

func (tree *IncrementalAndUpdatableMerkletree) InsertLeafWithDebug(leaf LeafData, debug bool) (err error) {
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
			if debug {
				log.Printf("InsertLeaf - Pos: Left - Level: %d - Hash: %s - Zero: %s\n", i, hash.Text(16), tree.Zeros[i].Text(16))
			}
		} else {
			// RIGHT
			tree.LastSubtrees[i] = [2]*big.Int{tree.LastSubtrees[i][0], hash}
			hashInputs = append(hashInputs, tree.LastSubtrees[i][0], hash)
			if debug {
				log.Printf("InsertLeaf - Pos: Right - Level: %d - LastSubtrees: %s - Hash: %s\n", i, tree.LastSubtrees[i][0].Text(16), hash.Text(16))
			}
		}
		hash, err = poseidon.Hash(hashInputs)
		if err != nil {
			return
		}
		if debug {
			log.Printf("InsertLeaf - Level: %d - Calculated Hash: %s\n", i, hash.Text(16))
		}
		index >>= 1
	}
	tree.BaseItems[uint(tree.NumberOfLeaves)] = leaf.LeafHash
	tree.NumberOfLeaves++
	if debug {
		log.Printf("InsertLeaf - Root Hash: %s\n", hash.Text(16))
	}
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
	allLeaves, err := tree.BuildTree(logDebug)
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

func (tree *IncrementalAndUpdatableMerkletree) UpdateLeafAtLevel(index int64, level uint, leaf LeafData, logDebug bool) (err error) {
	if tree.NumberOfLeaves == TREE_MAX_DEPTH {
		err = fmt.Errorf("tree reached max leaves")
		return
	}
	levelsAbove := tree.Depth - level
	if logDebug {
		log.Printf("UpdateLeafAtLevel - LevelsAbove: %d\n", levelsAbove)
	}
	zeros := make(map[uint]*big.Int)
	i := uint(0)
	for ; level < tree.Depth; level++ {
		zeros[i] = tree.Zeros[level]
		i++
	}
	if logDebug {
		log.Printf("UpdateLeafAtLevel - zeros: %+v\n", zeros)
	}
	tmpMk, err := NewPopulatedIncrementalAndUpdatableMerkletree(levelsAbove, zeros)
	if err != nil {
		return
	}
	if logDebug {
		log.Printf("UpdateLeafAtLevel - Before - Num of leaves: %+v - Root: %s\n", tmpMk.NumberOfLeaves, tmpMk.Root.Text(16))
	}
	err = tmpMk.InsertLeafWithDebug(leaf, false)
	if err != nil {
		return
	}
	if logDebug {
		log.Printf("UpdateLeafAtLevel - After - Num of leaves: %+v - Item Zero: %s - Root: %s\n", tmpMk.NumberOfLeaves, tmpMk.BaseItems[0].Text(16), tmpMk.Root.Text(16))
	}
	tree.Root.SetBytes(tmpMk.Root.Bytes())
	return
}

func (tree *IncrementalAndUpdatableMerkletree) GenerateMerkleProofPath(index uint, logDebug bool) (path []uint8, siblings []*big.Int, err error) {
	if index >= uint(tree.NumberOfLeaves) {
		err = fmt.Errorf("invalid leaf index. it is greather than actual leaves inserted")
		return
	}
	pos := uint(index)
	var i uint
	for ; i < tree.Depth; i++ {
		if index&1 == 0 {
			// LEFT
			path = append(path, 0)
			// tree.LastSubtrees[i] = [2]*big.Int{hash, tree.Zeros[i]}
			siblings = append(siblings, tree.Zeros[i])
		} else {
			// RIGHT
			path = append(path, 1)
			// tree.LastSubtrees[i] = [2]*big.Int{tree.LastSubtrees[i][0], hash}
			siblings = append(siblings, tree.LastSubtrees[i][0])
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
	totalLeaves := uint(tree.NumberOfLeaves)
	var idx uint
	for ; idx < totalLeaves; idx++ {
		item := tree.BaseItems[idx]
		if item.Cmp(hashLeaf) == 0 {
			index = int64(idx)
			return
		}
	}
	err = fmt.Errorf("index not found")
	return
}

func (tree *IncrementalAndUpdatableMerkletree) BuildTree(debug bool) (allLeaves map[uint]map[uint]*big.Int, err error) {
	var i, z, numLeaves, treeDepth uint
	var item *big.Int
	allLeaves = make(map[uint]map[uint]*big.Int)
	treeDepth = uint(tree.Depth)
	if debug {
		log.Println(" === Starting BuildTree === ")
	}
	for i = 0; i < uint(tree.Depth); i++ {
		allLeaves[i] = make(map[uint]*big.Int)
		numLeaves = uint(math.Pow(2, float64(treeDepth-i)))
		if debug {
			log.Println(" Level: ", float64(treeDepth-i))
			log.Println(" Number of leaves: ", numLeaves)
		}
		for z = 0; z < numLeaves; z++ {
			if i < 1 {
				if debug {
					log.Printf(" Is Base Item %d null ? %t\n", z, tree.BaseItems[uint(z)] == nil)
				}
				if tree.BaseItems[uint(z)] != nil {
					item = tree.BaseItems[uint(z)]
				} else {
					item = tree.Zeros[i]
				}
				if debug {
					if item.Cmp(tree.Zeros[i]) != 0 {
						log.Printf(" Item value %s is different of what Zero value is %s\n", item.Text(10), tree.Zeros[i].Text(10))
					}
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
	if debug {
		log.Println(" === End BuildTree === ")
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
		numLeaves = uint(math.Pow(2, float64(treeDepth-i)))
		if debug {
			log.Println("Layer: ", float64(treeDepth-i))
			log.Println("Number of leaves: ", numLeaves)
		}
		index := 0
		if i < 1 {
			hash = tree.Zeros[i]
			tree.NumberOfLeaves = int(numLeaves)
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
