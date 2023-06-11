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
	mktree, err := NewIncrementalAndUpdatableMerkletree(26, TestZeroValue)
	if err != nil {
		t.Fatal("could not create merkletree. Error: " + err.Error())
	}
	t.Log("Merkletree root: " + mktree.Root.String())
	t.Logf("Merkletree root in hex: 0x%064s", mktree.Root.Text(16))
}

func TestMktreeCreationPopulated(t *testing.T) {
	zeros := make(map[uint]*big.Int)
	zeros[0], _ = big.NewInt(0).SetString("2e99dc37b0a4f107b20278c26562b55df197e0b3eb237ec672f4cf729d159b69", 16)
	zeros[1], _ = big.NewInt(0).SetString("225624653ac89fe211c0c3d303142a4caf24eb09050be08c33af2e7a1e372a0f", 16)
	zeros[2], _ = big.NewInt(0).SetString("276c76358db8af465e2073e4b25d6b1d83f0b9b077f8bd694deefe917e2028d7", 16)
	zeros[3], _ = big.NewInt(0).SetString("09df92f4ade78ea54b243914f93c2da33414c22328a73274b885f32aa9dea718", 16)
	zeros[4], _ = big.NewInt(0).SetString("1c78b565f2bfc03e230e0cf12ecc9613ab8221f607d6f6bc2a583ccd690ecc58", 16)
	zeros[5], _ = big.NewInt(0).SetString("2879d62c83d6a3af05c57a4aee11611a03edec5ff8860b07de77968f47ff1c5f", 16)
	zeros[6], _ = big.NewInt(0).SetString("28ad970560de01e93b613aabc930fcaf087114743909783e3770a1ed07c2cde6", 16)
	zeros[7], _ = big.NewInt(0).SetString("27ca60def9dd0603074444029cbcbeaa9dbe77668479ac1db738bb892d9f3b6d", 16)
	zeros[8], _ = big.NewInt(0).SetString("28e4c1e90bbfa69de93abf6cbdc7cd1c0753a128e83b2b3afe34e0471a13ff55", 16)
	zeros[9], _ = big.NewInt(0).SetString("1b89c44a9f153266ad5bf754d4b252c26acba7d21fc661b94dc0618c6a82f49c", 16)
	zeros[10], _ = big.NewInt(0).SetString("0a5e5ec37bd8f9a21a1c2192e7c37d86bf975d947c2b38598b00babe567191c9", 16)
	zeros[11], _ = big.NewInt(0).SetString("21fb04b171b68944c640020a3a464602ec8d02495c44f1e403d9be4a97128e49", 16)
	zeros[12], _ = big.NewInt(0).SetString("19151c748859974805eb30feac7a301266dec9f67e23e285fe750f86448a2af9", 16)
	zeros[13], _ = big.NewInt(0).SetString("18fb0b755218eaa809681eb87e45925faa9197507d368210d73b5836ebf139e4", 16)
	zeros[14], _ = big.NewInt(0).SetString("1e294375b42dfd97795e07e1fe8bd6cefcb16c3bbb71b30bed950f8965861244", 16)
	zeros[15], _ = big.NewInt(0).SetString("0d3e4235db275d9bab0808dd9ade8789d46d0e1f1c9a99ce73fefca51dc92f4a", 16)
	zeros[16], _ = big.NewInt(0).SetString("075ab2ca945c4dc5ea40a9f1c66d5bf3c367cef1e04e73aa17c2bc747eb5fc87", 16)
	zeros[17], _ = big.NewInt(0).SetString("26f0f533a8ea2210001aeb8f8306c7c70656ba6afe145c6540bd4ed2c967a230", 16)
	zeros[18], _ = big.NewInt(0).SetString("24be7e64f680326e6e3621e5862d7b6b1f31e9e183a0bf5dd04e823be84e6af9", 16)
	zeros[19], _ = big.NewInt(0).SetString("212b13c9cbf421942ae3e3c62a3c072903c2a745a220cfb3c43cd520f55f44bf", 16)
	mktree, err := NewPopulatedIncrementalAndUpdatableMerkletree(20, zeros)
	if err != nil {
		t.Fatal("could not create merkletree. Error: " + err.Error())
	}
	t.Log("Merkletree root: " + mktree.Root.String())
	t.Logf("Merkletree root in hex: 0x%064s", mktree.Root.Text(16))
}

func TestMktreeFullZeroValue(t *testing.T) {
	t.Log("Starting TestMktreeFullZeroValue...")
	t.Log("Creating merkletree...")
	mktree, err := NewIncrementalAndUpdatableMerkletree(26, TestZeroValue)
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
	mktree.PopulateTreeWithZeros(false)
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

func TestMktreePopulatedUpdateLastItem(t *testing.T) {
	t.Log("Starting TestMktreePopulatedUpdateLastItem...")
	t.Log("Creating merkletree...")
	mktree, err := NewIncrementalAndUpdatableMerkletree(26, TestZeroValue)
	if err != nil {
		t.Fatal("could not create merkletree. Error: " + err.Error())
	}
	t.Log("Merkletree has been created.")
	t.Log("Merkletree root: " + mktree.Root.String())
	t.Logf("Merkletree root in hex: 0x%064s", mktree.Root.Text(16))
	t.Log("Updating Merkletree ...")
	ok := false
	leaf := LeafData{}
	// leaf.LeafHash, ok = big.NewInt(0).SetString("22fbeb55d785f2d0723d9325756a56c25f2e4006f1a0427ab17516fd033ed3f5", 16)
	leaf.LeafHash, ok = big.NewInt(0).SetString("0fb2bc182d159c8c843b6f69c2fed99955c6007787c4c573ce65edb4f15cd9df", 16)
	if !ok {
		t.Fatal("could not convert to big int the new leaf hash. Error: " + err.Error())
	}
	err = mktree.UpdateLeafAtLevel(0, 6, leaf, false)
	if err != nil {
		t.Fatal("could not Update Leaf. Error: " + err.Error())
	}
	t.Log("Merkletree updated.")
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
	leaf.LeafHash, ok = big.NewInt(0).SetString("22fbeb55d785f2d0723d9325756a56c25f2e4006f1a0427ab17516fd033ed3f5", 16)
	if !ok {
		t.Fatal("could not convert to big int the new leaf hash. Error: " + err.Error())
	}
	mktree.InsertLeaf(leaf)
	t.Log("Merkletree root: " + mktree.Root.String())
	t.Logf("Merkletree root in hex: 0x%064s", mktree.Root.Text(16))
}

func TestMktreeProofGeneration(t *testing.T) {
	t.Log("Starting TestMktreeProofGeneration...")
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
	log.Println("Starting GenerateMerkleProofPath...")
	path, siblings, err := mktree.GenerateMerkleProofPath(0, false)
	if err != nil {
		t.Fatal("could not Generate Merkle Proof Path. Error: " + err.Error())
	}
	log.Println("Finished GenerateMerkleProofPath.")
	t.Logf("Merkletree proof path: %+v\n", path)
	t.Logf("Merkletree proof siblings: %+v\n", siblings)
}
