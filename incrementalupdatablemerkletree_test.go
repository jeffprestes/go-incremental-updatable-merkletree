package goincrementalupdatablemerkletree

import (
	"log"
	"math/big"
	"testing"
)

var TestZeroValue *big.Int

func TestMain(m *testing.M) {
	var ok bool
	// TestZeroValue, ok = big.NewInt(0).SetString("667764c376602b72ef22218e1673c2cc8546201f9a77807570b3e5de137680d", 16)
	TestZeroValue, ok = big.NewInt(0).SetString("0667764c376602b72ef22218e1673c2cc8546201f9a77807570b3e5de137680d", 16)
	if !ok {
		log.Fatalln("could not compute zero value")
		return
	}
	m.Run()
}

func TestMktreeCreation(t *testing.T) {
	mktree, err := NewIncrementalAndUpdatableMerkletree(20, TestZeroValue)
	if err != nil {
		t.Fatal("could not create merkletree. Error: " + err.Error())
	}
	t.Log("Merkletree root: " + mktree.Root.String())
	t.Logf("Merkletree root in hex: 0x%064s", mktree.Root.Text(16))
}

func TestMktreeFullZeroValue(t *testing.T) {
	t.Log("Starting TestMktreeFullZeroValue...")
	t.Log("Creating merkletree...")
	mktree, err := NewIncrementalAndUpdatableMerkletree(20, TestZeroValue)
	if err != nil {
		t.Fatal("could not create merkletree. Error: " + err.Error())
	}
	t.Log("Merkletree has been created.")
	t.Log("Populating Merkletree with zeros...")
	mktree.PopulateTreeWithZeros(true)
	t.Log("Merkletree filled with zeros.")
	t.Logf("Total of Merkletree's base leaves: %d\n", mktree.NumberOfLeaves)
	t.Log("Merkletree root: " + mktree.Root.String())
	t.Logf("Merkletree root in hex: 0x%064s", mktree.Root.Text(16))
}

func TestMktreeUpdateLastItem(t *testing.T) {
	t.Log("Starting TestMktreeUpdateLastItem...")
	t.Log("Creating merkletree...")
	mktree, err := NewIncrementalAndUpdatableMerkletree(20, TestZeroValue)
	if err != nil {
		t.Fatal("could not create merkletree. Error: " + err.Error())
	}
	t.Log("Merkletree has been created.")
	t.Log("Populating Merkletree with zeros...")
	mktree.PopulateTreeWithZeros(true)
	t.Log("Merkletree filled with zeros.")
	t.Logf("Total of Merkletree's base leaves: %d\n", mktree.NumberOfLeaves)
	t.Log("Updating last leaf...")
	ok := false
	leaf := LeafData{}
	leaf.LeafHash, ok = big.NewInt(0).SetString("1bdded415724018275c7fcc2f564f64db01b5bbeb06d65700564b05c3c59c9e6", 16)
	if !ok {
		t.Fatal("could not convert to big int the new leaf hash. Error: " + err.Error())
	}
	index := int64(mktree.NumberOfLeaves - 1)
	mktree.UpdateLeaf(index, leaf, true)
	t.Log("Merkletree root: " + mktree.Root.String())
	t.Logf("Merkletree root in hex: 0x%064s", mktree.Root.Text(16))
}

func TestMktreeUpdateLastItemWithEmptyTree(t *testing.T) {
	t.Log("Starting TestMktreeUpdateLastItem...")
	t.Log("Creating merkletree...")
	mktree, err := NewIncrementalAndUpdatableMerkletree(20, TestZeroValue)
	if err != nil {
		t.Fatal("could not create merkletree. Error: " + err.Error())
	}
	t.Log("Merkletree has been created.")
	t.Log("Updating last leaf...")
	ok := false
	leaf := LeafData{}
	leaf.LeafHash, ok = big.NewInt(0).SetString("1bdded415724018275c7fcc2f564f64db01b5bbeb06d65700564b05c3c59c9e6", 16)
	if !ok {
		t.Fatal("could not convert to big int the new leaf hash. Error: " + err.Error())
	}
	mktree.InsertLeaf(leaf)
	t.Log("Merkletree root: " + mktree.Root.String())
	t.Logf("Merkletree root in hex: 0x%064s", mktree.Root.Text(16))
}
