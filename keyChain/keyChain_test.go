package keyChain

import (
	"bytes"
	"testing"
)

// func testChildTree(t *testing.T, child KeyChain, parent KeyChain) {
// 	if len(parent.NameSpace) + 1 != len(child.NameSpace) {
// 		t.Error("child's namespace was not 1 element longer than parent")
// 	}
// 	if parent.NameSpace.Equal(child.NameSpace[:len(parent.NameSpace)]) {
// 		t.Error("child's namespace does not contain parent's namespace")
// 	}
// 	if !child.Id.Equal(child.NameSpace[len(child.NameSpace) - 1]) {
// 		t.Error("child should be in it's own namespace")
// 	}
// 	if child.Id.Equal(parent.Id) {
// 		t.Error("")
// 	}
// 	if len(child.Ancestors) != 0 {
// 		t.Error("")
// 	}
// }

func TestMakeTreeKeyChainRoot (t *testing.T) {
	tree, err := Root.MakeChildTree()

	if err != nil {
		t.Error("error making tree: ", err)
	}

	tId := tree.Id

	if tree.GetLoc().Equal(rootLoc) {
		t.Error("tree's self is equal to root location: ", tree.GetLoc())
	}

	if tId.Height != 1 {
		t.Error("something went wrong making the tree's Id: ", tId)
	}

	if !tree.GetLoc().Equal(tree.GetChildBucket()) {
		t.Error("tree's self and ChildBucket should have been the same: ", 
			"\nself: ", tree.GetLoc(),
			"\nchildBucket: ", tree.GetChildBucket())
	}

	if !tree.GetParentLoc().Equal(rootLoc) {
		t.Error("tree's parent should have been equal to root Location: ", tree.GetParentLoc())
	}

}

func TestMakeTreeKeyChainOnTree (t *testing.T) {
	parent, err := Root.MakeChildTree()

	if err != nil {
		t.Error("error making parent: ", err)
	}

	child, err := parent.MakeChildTree()

	if err != nil {
		t.Error("error making tree: ", err)
	}

	cId := child.Id

	if child.GetLoc().Equal(parent.GetLoc()) {
		t.Error("tree's self is equal to root location: ", child.GetLoc())
	}

	if cId.Height != 2 {
		t.Error("something went wrong making the tree's Id: ", cId)
	}

	if !child.GetLoc().Equal(child.GetChildBucket()) {
		t.Error("tree's self and ChildBucket should have been the same: ", 
			"\nself: ", child.GetLoc(),
			"\nchildBucket: ", child.GetChildBucket())
	}

	if !child.GetParentLoc().Equal(parent.GetLoc()) {
		t.Error("tree's parent should have been equal to parent Location \nchild: ", child.GetParentLoc(), "\nparent: ", parent.GetLoc())
	}
}

func TestMakeBranchKeyChainRoot (t *testing.T) {
	branch, err := Root.MakeChildBranch()

	if err != nil {
		t.Error("error making branch: ", err)
	}

	if branch.GetLoc().Equal(rootLoc){
		t.Error("branch's self was equal to rootLoc: ", branch.GetLoc())
	}

	if bId := branch.Id; bId.Height != 1 {
		t.Error("branch's Id was set incorrectly: ", branch.GetLoc())
	}

	if !branch.GetChildBucket().Equal(rootLoc.copyAndAppend(branch.Id)) {
		t.Error("branch's child bucket should have been it's parent's self plus the branch's Id: ", branch.GetChildBucket())
	}

	if !branch.GetParentLoc().Equal(rootLoc) {
		t.Error("branch's parent was not equal to it's parent's self: ", 
			"\nbranch: ", branch, 
			"\nbranch ParentLoc: ", branch.GetParentLoc(),
			"\nparent: ", Root)
	}

}

func TestMakeBranchKeyChainOnTree (t *testing.T) {
	parent, err := Root.MakeChildTree()

	if err != nil {
		t.Error("error making parent: ", err)
	}

	child, err := parent.MakeChildBranch()

	if err != nil {
		t.Error("error making tree: ", err)
	}

	cId := child.Id

	if cId.Height != 2 {
		t.Error("something went wrong making the branch's Id: ", cId)
	}

	if child.GetLoc().Equal(parent.GetLoc()) {
		t.Error("branch's self was equal to rootLoc: ", child.GetLoc())
	}

	if !child.GetChildBucket().Equal(parent.NameSpace.copyAndAppend(child.Id)) {
		t.Error("branch's Child Bucket should be parent's Loc plus branch's ID: ", 
			"\nparent: ", parent.GetChildBucket(),
			"\nchild: ", child.GetChildBucket())
	}

	if !child.GetParentLoc().Equal(parent.GetLoc()) {
		t.Error("branch's parent should have been equal to parent's self: ", child.GetParentLoc())
	}

}

func TestMakeBranchKeyChainOnBranch (t *testing.T) {
	parent, err := Root.MakeChildBranch()

	if err != nil {
		t.Error("error making branch")
	}

	child, err := parent.MakeChildBranch()

	if err != nil {
		t.Error("error making branch")
	}

	cId := child.Id

	if cId.Height != 2 {
		t.Error("something went wrong making the branch's Id: ", cId)
	}

	if child.GetLoc().Equal(parent.GetLoc()) {
		t.Error("branch's self was equal to parent's self: ", child.GetLoc())
	}

	if !child.GetChildBucket().Equal(parent.NameSpace.copyAndAppend(child.Id)) {
		t.Error("branch's Child Bucket should be parent's NameSpace plus branch's ID: ", 
			"\nparent: ", parent.GetChildBucket(),
			"\nchild: ", child.GetChildBucket())
	}

	if !child.GetParentLoc().Equal(parent.GetLoc()) {
		t.Error("branch's parent should have been equal to parent's self: ", child.GetParentLoc())
	}

}

func TestIsTree (t *testing.T) {
	tree, err := Root.MakeChildTree()

	if err != nil {
		t.Error("error making tree: ", err)
	}

	treeIsTree := tree.IsTree()

	if !treeIsTree {
		t.Error("tree is not a tree: ", tree)
	}

	branch, err := Root.MakeChildBranch()

	if err != nil {
		t.Error("error making branch: ", err)
	}

	branchIsTree := branch.IsTree()

	if branchIsTree {
		t.Error("branch is a tree: ", branch)
	}

}

func TestPassThroughFunctions (t *testing.T) {
	parent, err := Root.MakeChildTree()

	if err != nil {
		t.Error("error making parent: ", err)
	}

	child, err := parent.MakeChildTree()

	if err != nil {
		t.Error("error making child: ", err)
	}

	if !bytes.Equal(child.Key(), child.GetLoc().Key()) {
		t.Error("keychain Key is not passing through correctly")
	}

	if bytes.Equal(child.Key(), parent.Key()) {
		t.Error("child key equal to parent key!")
	}

	if child.KeyString() != child.GetLoc().KeyString() {

	}

	if child.KeyString() == parent.KeyString() {
		t.Error("child key equal to parent key!")
	}

	if !child.Equal(child) {
		t.Error("child is not equal to child: ", child)
	}

	if !parent.Equal(parent) {
		t.Error("parent is not equal to parent: ", parent)
	}

	if parent.Equal(child) {
		t.Error("parent is equal to child: ", 
			"\nparent: ", parent,
			"\nchild: ", child)
	}


}