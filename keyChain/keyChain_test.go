package keyChain

import (
	"bytes"
	"testing"
)

func TestMakeTreeKeyChainRoot (t *testing.T) {
	tree, err := MakeTreeKeyChain(Root)

	if err != nil {
		t.Error("error making tree: ", err)
	}

	tId := tree.Self.GetId()

	if tree.Self.Equal(rootLoc) {
		t.Error("tree's self is equal to root location: ", tree.Self)
	}

	if tId.Height != 1 {
		t.Error("something went wrong making the tree's Id: ", tId)
	}

	if !tree.Self.Equal(tree.ChildBucket) {
		t.Error("tree's self and ChildBucket should have been the same: ", 
			"\nself: ", tree.Self,
			"\nchildBucket: ", tree.ChildBucket)
	}

	if !tree.Parent.Equal(rootLoc) {
		t.Error("tree's parent should have been equal to root Location: ", tree.Parent)
	}

}

func TestMakeTreeKeyChainOnTree (t *testing.T) {
	parent, err := MakeTreeKeyChain(Root)

	if err != nil {
		t.Error("error making parent: ", err)
	}

	child, err := MakeTreeKeyChain(parent)

	if err != nil {
		t.Error("error making tree: ", err)
	}

	cId := child.Self.GetId()

	if child.Self.Equal(parent.Self) {
		t.Error("tree's self is equal to root location: ", child.Self)
	}

	if cId.Height != 2 {
		t.Error("something went wrong making the tree's Id: ", cId)
	}

	if !child.Self.Equal(child.ChildBucket) {
		t.Error("tree's self and ChildBucket should have been the same: ", 
			"\nself: ", child.Self,
			"\nchildBucket: ", child.ChildBucket)
	}

	if !child.Parent.Equal(parent.Self) {
		t.Error("tree's parent should have been equal to root Location: ", child.Parent)
	}
}

func TestMakeBranchKeyChainRoot (t *testing.T) {
	branch, err := MakeBranchKeyChain(Root)

	if err != nil {
		t.Error("error making branch: ", err)
	}

	if branch.Self.Equal(rootLoc){
		t.Error("branch's self was equal to rootLoc: ", branch.Self)
	}

	if bId := branch.Self.GetId(); bId.Height != 1 {
		t.Error("branch's Id was set incorrectly: ", branch.Self)
	}

	if !branch.ChildBucket.Equal(rootLoc) {
		t.Error("branch's child bucket should have been it's parent's self: ", branch.ChildBucket)
	}

	if !branch.Parent.Equal(rootLoc) {
		t.Error("branch's parent was not equal to it's parent's self: ", branch.Parent)
	}

}

func TestMakeBranchKeyChainOnTree (t *testing.T) {
	parent, err := MakeTreeKeyChain(Root)

	if err != nil {
		t.Error("error making parent: ", err)
	}

	child, err := MakeBranchKeyChain(parent)

	if err != nil {
		t.Error("error making tree: ", err)
	}

	cId := child.Self.GetId()

	if cId.Height != 2 {
		t.Error("something went wrong making the branch's Id: ", cId)
	}

	if child.Self.Equal(parent.Self) {
		t.Error("branch's self was equal to rootLoc: ", child.Self)
	}

	if !child.ChildBucket.Equal(parent.ChildBucket) {
		t.Error("branch's and parent's ChildBuckets should have been the same: ", 
			"\nparent: ", parent.ChildBucket,
			"\nchild: ", child.ChildBucket)
	}

	if !child.Parent.Equal(parent.Self) {
		t.Error("branch's parent should have been equal to parent's self: ", child.Parent)
	}

}

func TestMakeBranchKeyChainOnBranch (t *testing.T) {
	parent, err := MakeBranchKeyChain(Root)

	if err != nil {
		t.Error("error making branch")
	}

	child, err := MakeBranchKeyChain(parent)

	if err != nil {
		t.Error("error making branch")
	}

	cId := child.Self.GetId()

	if cId.Height != 2 {
		t.Error("something went wrong making the branch's Id: ", cId)
	}

	if child.Self.Equal(parent.Self) {
		t.Error("branch's self was equal to parent's self: ", child.Self)
	}

	if !child.ChildBucket.Equal(parent.ChildBucket) {
		t.Error("branch's and parent's ChildBuckets should have been the same: ", 
			"\nparent: ", parent.ChildBucket,
			"\nchild: ", child.ChildBucket)
	}

	if !child.Parent.Equal(parent.Self) {
		t.Error("branch's parent should have been equal to parent's self: ", child.Parent)
	}

}

func TestIsTree (t *testing.T) {
	tree, err := MakeTreeKeyChain(Root)

	if err != nil {
		t.Error("error making tree: ", err)
	}

	treeIsTree := tree.IsTree()

	if !treeIsTree {
		t.Error("tree is not a tree: ", tree)
	}

	branch, err := MakeBranchKeyChain(Root)

	if err != nil {
		t.Error("error making branch: ", err)
	}

	branchIsTree := branch.IsTree()

	if branchIsTree {
		t.Error("branch is a tree: ", branch)
	}

}

func TestPassThroughFunctions (t *testing.T) {
	parent, err := MakeTreeKeyChain(Root)

	if err != nil {
		t.Error("error making parent: ", err)
	}

	child, err := MakeTreeKeyChain(parent)

	if err != nil {
		t.Error("error making child: ", err)
	}

	if !bytes.Equal(child.Key(), child.Self.Key()) {
		t.Error("keychain Key is not passing through correctly")
	}

	if bytes.Equal(child.Key(), parent.Key()) {
		t.Error("child key equal to parent key!")
	}

	if child.KeyString() != child.Self.KeyString() {

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