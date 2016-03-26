package keyChain

import (
	"testing"
	"bytes"
)

func TestMakeTreeLocH1 (t *testing.T) {
	treeLoc, err := makeTreeLoc(rootLoc)

	if err != nil {
		t.Error("error making location", err)
	}

	if len(treeLoc) != 2 {
		t.Error("tree's length should have been to, but was: ", len(treeLoc))
	}

	if treeLoc[1].Height != 1 {
		t.Error("tree's Id had the wrong height: ", treeLoc[1].Height)
	}

	if !treeLoc[0].Equal(rootId) {
		t.Error("the id in the tree that should have been the root, was not the root: ", treeLoc[0])
	}
}

func TestMakeTreeLocH2 (t *testing.T) {
	treeLoc, err := makeTreeLoc(rootLoc)

	if err != nil {
		t.Error("error making location", err)
	}

	childTreeLoc, err := makeTreeLoc(treeLoc)

	if err != nil {
		t.Error("error making location: ", err)
	}

	if len(childTreeLoc) != 3 {
		t.Error("tree's length should have been to, but was: ", len(childTreeLoc))
	}

	if childTreeLoc[2].Height != 2 {
		t.Error("tree's Id had the wrong height: ", childTreeLoc[1].Height)
	}

	if !childTreeLoc[0].Equal(rootId) {
		t.Error("the id in the child tree that should have been the root Id was not the root Id: ", childTreeLoc[0])
	}

	if !childTreeLoc[1].Equal(treeLoc[1]) {
		t.Error("the id in the child tree that should have been the parent tree's Id was not the parent tree Id.",
			"\nparent: ", treeLoc[1],
			"\nchild: ", childTreeLoc[1],
			)
	}
}

func TestMakeBranchLocH1 (t *testing.T) {
	branchLoc, err := makeBranchLoc(rootLoc, rootLoc)

	if err != nil {
		t.Error("error making location: ", err)
	}

	if len(branchLoc) != 3 {
		t.Error("branch length should have been 3, but was: ", len(branchLoc))
	}

	if branchLoc[2].Height != 1 {
		t.Error("branch's id was not set correctly")
	}

	if !branchLoc[0].Equal(rootId) {
		t.Error("branch should have been placed in rootnameSpace: ", branchLoc[0])
	}

	if !branchLoc[1].Equal(rootId) {
		t.Error("branch's parent should have been root: ", branchLoc[1])
	}

}

func TestMakeBranchOnBranch (t *testing.T) {
	branchLoc, err := makeBranchLoc(rootLoc, rootLoc)

	if err != nil {
		t.Error("error making location: ", err)
	}

	child, err := makeBranchLoc(rootLoc, branchLoc)

	if err != nil {
		t.Error("error making location: ", err)
	}

	if len(child) != len(branchLoc) { //len(childBranch) == len(parentBranch) but not len(parentTree)
		t.Error("branch's length should have been the same as it's parent, but was: ", len(child))
	}

	if child[2].Height != 2 {
		t.Error("child's Height should have been 2, but was: ", child[2].Height)
	}

	if !child[0].Equal(rootId) {
		t.Error("child branch should have been in the same nameSpace as it's parent, but was in: ", child[0])
	}

	if !child[1].Equal(branchLoc[2]) {
		t.Error("the element before child branch's id should have been it's parent Id.",
			"\nparent: ", branchLoc[2],
			"\nchild: ", child[1],
			)
	}
}

func TestMakeBranchOnTree (t *testing.T) {
	treeLoc, err := makeTreeLoc(rootLoc)

	if err != nil {
		t.Error("error making tree: ", err)
	}

	child, err := makeBranchLoc(treeLoc, treeLoc)

	if err != nil {
		t.Error("error making branch: ", err)
	}

	if len(child) != len(treeLoc) + 2 {
		t.Error("branch length should have been 4, but was: ", child)
	}

	if child[3].Height != 2 {
		t.Error("child's Height should have been 2, but was: ", child[3].Height)
	}

	if !child[0].Equal(rootId) {
		t.Error("child branch should have been in the same nameSpace as it's parent, but was in: ", child[0])
	}

	if !child[1].Equal(treeLoc[1]) {
		t.Error("the element before child branch's id should have been it's parent Id.",
			"\nparent: ", treeLoc[1],
			"\nchild: ", child[1],
			)
	}

	if !child[2].Equal(treeLoc[1]) {
		t.Error("the element before child branch's id should have been it's parent Id.",
			"\nparent: ", treeLoc[1],
			"\nchild: ", child[2],
			)
	}
}

func TestMakeBranchOnBranchOnTree (t *testing.T) {
	treeLoc, err := makeTreeLoc(rootLoc)

	if err != nil {
		t.Error("error making tree: ", err)
	}

	parent, err := makeBranchLoc(treeLoc, treeLoc)

	if err != nil {
		t.Error("error making branch: ", err)
	}

	child, err := makeBranchLoc(treeLoc, parent)

	if err != nil {
		t.Error("error making branch: ", err)
	}

	if len(child) != len(parent) {
		t.Error("branch length should have been 4, but was: ", child)
	}

	if child[3].Height != 3 {
		t.Error("child's Height should have been 2, but was: ", child[3].Height)
	}

	if !child[0].Equal(rootId) {
		t.Error("child branch should have been in the same nameSpace as it's parent, but was in: ", child[0])
	}

	if !child[1].Equal(parent[1]) {
		t.Error("the child and parent should be in the same nameSpace",
			"\nparent: ", parent[1],
			"\nchild: ", child[1],
			)
	}

	if !child[2].Equal(parent[3]) {
		t.Error("the element before child branch's id should have been it's parent Id.",
			"\nparent: ", parent[3],
			"\nchild: ", child[2],
			)
	}
}

func TestGetId (t *testing.T) {
	rId := rootLoc.GetId()

	if !rId.Equal(rootId) {
		t.Error("rootLoc's Id was not root Id: ", rId)
	}


	tree, err := makeTreeLoc(rootLoc)

	if err != nil {
		t.Error("error making tree", err)
	}

	id := tree.GetId()

	if !id.Equal(tree[1]) {
		t.Error("GetId got the wrong id: ", id)
	}
}

func TestLocKey (t *testing.T) {
	rootKey := rootLoc.Key()

	if !bytes.Equal(rootKey, []byte{}) {
		t.Error("rootLoc's key should have been empty, but was: ", rootKey)
	}

	tree1, err := makeTreeLoc(rootLoc)

	if err != nil {
		t.Error("error making tree: ", err)
	}

	treeKey1 := tree1.Key()

	if !bytes.Equal(treeKey1, tree1[1].Key()) {
		t.Error("first tree's key should have been the key of it's id: ", treeKey1)
	}

	tree2, err := makeTreeLoc(tree1)

	if err != nil {
		t.Error("error making tree: ", err)
	}

	treeKey2 := tree2.Key()

	if !bytes.Equal(treeKey2, append(tree1.Key(), tree2[2].Key()...)) {
		t.Error("second tree's key should have been the combination of it's namespace and Id")
	}
}

func TestLocKeyString (t *testing.T) {
	rootKey := rootLoc.KeyString()

	if rootKey != "" {
		t.Error("rootLoc's key should have been empty, but was: ", rootKey)
	}

	tree1, err := makeTreeLoc(rootLoc)

	if err != nil {
		t.Error("error making tree: ", err)
	}

	treeKey1 := tree1.KeyString()

	if treeKey1 != string(tree1[1].Key()) {
		t.Error("first tree's key should have been the key of it's id: ", treeKey1)
	}

	tree2, err := makeTreeLoc(tree1)

	if err != nil {
		t.Error("error making tree: ", err)
	}

	treeKey2 := tree2.KeyString()

	if treeKey2 != string(append(tree1.Key(), tree2[2].Key()...)) {
		t.Error("second tree's key should have been the combination of it's namespace and Id")
	}	
}

func TestTreeLocEqual (t *testing.T) {
	if !rootLoc.Equal(rootLoc) {
		t.Error("rootLoc is not equal to rootLoc... wut?", rootLoc)
	}

	tree1, err := makeTreeLoc(rootLoc)

	if err != nil {
		t.Error("error making tree: ", err)
	}

	tree2, err := makeTreeLoc(rootLoc)

	if err != nil {
		t.Error("error making tree: ", err)
	}

	tree3, err := makeTreeLoc(tree1)

	if !tree1.Equal(tree1) {
		t.Error("tree 1 is not equal to tree 1: ", tree1)
	}

	if !tree3.Equal(tree3) {
		t.Error("tree 3 is not equal to tree 3: ", tree3)
	}

	if tree1.Equal(tree2) {
		t.Error("tree 1 should not be equal to tree 2 (unless the same uuid got generated twice in a row XD)")
	}

	if tree1.Equal(tree3) {
		t.Error("tree3 should not be equal to tree1")
	}
}

func TestBranchLocEqual (t *testing.T) {
	tree, err := makeTreeLoc(rootLoc)

	if err != nil {
		t.Error("error making tree: ", err)
	}

	branchT1, err := makeBranchLoc(rootLoc, tree)

	if err != nil {
		t.Error("error making branch: ", err)
	}

	branchT2, err := makeBranchLoc(rootLoc, tree)

	if err != nil {
		t.Error("error making branch: ", err)
	}

	branch1, err := makeBranchLoc(rootLoc, rootLoc)

	if err != nil {
		t.Error("error making branch: ", err)
	}

	branch11, err := makeBranchLoc(rootLoc, branch1)

	if err != nil {
		t.Error("error making branch: ", err)
	}

	branch12, err := makeBranchLoc(rootLoc, branch1)

	if err != nil {
		t.Error("error making branch: ", err)
	}

	branch2, err := makeBranchLoc(rootLoc, rootLoc)

	if err != nil {
		t.Error("error making branch: ", err)
	}

	if branchT1.Equal(tree) {
		t.Error("T1 and tree were equal")
	}

	if branchT1.Equal(branchT2) {
		t.Error("T1 and T2 were equal")
	}

	if branch1.Equal(tree) {
		t.Error("1 and tree were equal")
	}

	if branch1.Equal(branchT1) {
		t.Error("1 and T1 were equal")
	}

	if branch1.Equal(branch2) {
		t.Error("1 and 2 were equal")
	}

	if branch1.Equal(branch11) {
		t.Error("1 and 11 were equal")
	}

	if branch11.Equal(branch12) {
		t.Error("12 and 11 were equal")
	}

}