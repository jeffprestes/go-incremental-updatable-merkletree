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
	t.Log("Starting TestMktreeZeroValue...")
	t.Log("Creating merkletree...")
	mktree, err := NewIncrementalAndUpdatableMerkletree(20, TestZeroValue)
	if err != nil {
		t.Fatal("could not create merkletree. Error: " + err.Error())
	}
	t.Log("Merkletree has been created.")
	t.Log("Populating Merkletree with zeros...")
	mktree.PopulateTreeWithZeros(true)
	t.Log("Merkletree filled with zeros.")
	t.Log("Merkletree root: " + mktree.Root.String())
	t.Logf("Merkletree root in hex: 0x%064s", mktree.Root.Text(16))
}
