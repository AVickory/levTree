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

func TestMakeBranchKeyChainOnBranchOnBranch (t *testing.T) {
	grandParent, err := Root.MakeChildBranch()

	if err != nil {
		t.Error("error making branch")
	}

	parent, err := grandParent.MakeChildBranch()

	if err != nil {
		t.Error("error making branch")
	}

	child, err := parent.MakeChildBranch()

	if err != nil {
		t.Error("error making branch")
	}

	cId := child.Id

	if cId.Height != 3 {
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

func TestMakeSiblingTree (t *testing.T) {
	tree1, err := Root.MakeChildTree()

	if err != nil {
		t.Error("error making tree1: ", err)
	}

	tree2, err := tree1.MakeSibling()

	if err != nil {
		t.Error("error making sibling tree: ", err)
	}

	if !tree2.GetParentLoc().Equal(Root.GetLoc()) {
		t.Error("sibling should have had same parent as original")
	}

	if !tree1.GetSiblingBucket().Equal(tree2.GetSiblingBucket()) {
		t.Error("sibling trees sibling buckets are not the same", 
			"\ntree1: ", tree1,//.GetSiblingBucket(),
			"\ntree2: ", tree2,//.GetSiblingBucket(),)
			)
	}

	if !tree2.GetSiblingBucket().Equal(Root.GetChildBucket()) {
		t.Error("sibling tree's sibling bucket is not the parent's child bucket")
	}
}

func TestMakeSiblingBranch (t *testing.T) {
	branch1, err := Root.MakeChildBranch()

	if err != nil {
		t.Error("error making branch")
	}

	branch2, err := branch1.MakeChildBranch()

	if err != nil {
		t.Error("error making branch")
	}

	branch3, err := branch2.MakeChildBranch()

	if err != nil {
		t.Error("error making branch")
	}

	branch1S, err := branch1.MakeSibling()

	if err != nil {
		t.Error("error making branch")
	}

	branch2S, err := branch2.MakeSibling()

	if err != nil {
		t.Error("error making branch")
	}

	branch3S, err := branch3.MakeSibling()

	if err != nil {
		t.Error("error making branch")
	}

	if !branch1S.GetParentLoc().Equal(Root.GetLoc()) {
		t.Error("branch1S sibling has wrong parent: ", branch1S.GetParentLoc(),
			"\nbranch1S ParentLoc: ", branch1S.GetParentLoc(),
			"\nparent's loc: ", Root.GetLoc())
	}

	if !branch2S.GetParentLoc().Equal(branch1.GetLoc()) {
		t.Error("branch2S sibling has wrong parent: ", branch2S.GetParentLoc())
	}

	if !branch3S.GetParentLoc().Equal(branch2.GetLoc()) {
		t.Error("branch3S sibling has wrong parent: ", branch2S.GetParentLoc())
	}

	if !branch1S.GetSiblingBucket().Equal(Root.GetChildBucket()) {
		t.Error("branch1S sibling bucket is not in parent's bucket: ", 
			"\nchild: ", branch1S.GetSiblingBucket(),
			"\nparent: ", Root.GetChildBucket(),)
	}

	if !branch2S.GetSiblingBucket().Equal(branch1.GetChildBucket()) {
		t.Error("branch2S sibling is not in parent's bucket: ", branch2S.GetSiblingBucket())
	}

	if !branch3S.GetSiblingBucket().Equal(branch2.GetChildBucket()) {
		t.Error("branch3S sibling is not in parent's bucket: ", branch2S.GetSiblingBucket())
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

func TestParentIsTree (t *testing.T) {
	tree, err := Root.MakeChildTree()

	if err != nil {
		t.Error("error making tree: ", err)
	}

	parent, err := tree.MakeChildBranch()

	if err != nil {
		t.Error("error making branch: ", err)
	}

	child, err := parent.MakeChildBranch()

	if err != nil {
		t.Error("error making branch: ", err)
	}

	if !parent.ParentIsTree() {
		t.Error("parent's parent is a tree but ParentIsTree returned false.")
	}

	if child.ParentIsTree() {
		t.Error("child's parent is not a tree but ParentIsTree returned true.")
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

	if !bytes.Equal(child.ParentKey(), child.GetParentLoc().Key()) {
		t.Error("keychain ParentKey is not passing through correctly")
	}

	if !bytes.Equal(child.ParentKey(), parent.Key()) {
		t.Error("child key equal to parent key!")
	}	

	if !bytes.Equal(child.ChildKeyPrefix(), child.GetChildBucket().Key()) {
		t.Error("keyChain ChildKeyPrefix is not passing through correctly")
	}

	if !bytes.Equal(parent.ChildKeyPrefix(), child.GetLoc()[:len(child.GetLoc()) - 1].Key()) {
		t.Error("child's key does not have parent's child prefix")
	}

	if child.KeyString() != child.GetLoc().KeyString() {
		t.Error("child keystring not equal to child's locations keystring")
	}

	if child.KeyString() == parent.KeyString() {
		t.Error("child keyString equal to parent keyString!")
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