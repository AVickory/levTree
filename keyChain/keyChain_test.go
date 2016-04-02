package keyChain

import (
	"bytes"
	"testing"
	"fmt"
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

var dbm bool = false

func debug(stuff ...interface{}) {
	if dbm {
		fmt.Println(stuff...)
	}
}

func TestMakeTreeKeyChainRoot (t *testing.T) {
	tree, err := Root.MakeChildTree()

	if err != nil {
		t.Error("error making tree: ", err)
	}

	tId := tree.Id

	if tId.Height != 1 {
		t.Error("something went wrong making the tree's Id: ", tId)
	}

	expectedLoc := Loc{rootId, rootId, tId,}

	if !tree.GetLoc().Equal(expectedLoc) {
		t.Error("tree has wrong location",
			"\nexpected: ", expectedLoc,
			"\ncomputed: ", tree.GetLoc())
	}

	expectedChildBucket := Loc{rootId, tId, tId}

	if !tree.GetChildBucket().Equal(expectedChildBucket) {
		t.Error("tree has wrong child prefix",
			"\nexpected: ", expectedChildBucket,
			"\ncomputed: ", tree.GetChildBucket())	
	}

	expectedDescendantBucket := expectedChildBucket[:2]

	if !tree.GetDescendantBucket().Equal(expectedDescendantBucket) {
		t.Error("tree has wrong descendant prefix",
			"\nexpected: ", expectedDescendantBucket,
			"\ncomputed: ", tree.GetDescendantBucket())	
	}

	expectedSiblingBucket := expectedLoc[:2]

	if !tree.GetSiblingBucket().Equal(expectedSiblingBucket) {
		t.Error("tree has wrong sibling prefix",
			"\nexpected: ", expectedSiblingBucket,
			"\ncomputed: ", tree.GetSiblingBucket())		
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

	pId := parent.Id

	child, err := parent.MakeChildTree()

	if err != nil {
		t.Error("error making tree: ", err)
	}

	cId := child.Id

	if cId.Height != 2 {
		t.Error("something went wrong making the tree's Id: ", cId)
	}

	expectedLoc := Loc{rootId, pId, pId, cId,}

	if !child.GetLoc().Equal(expectedLoc) {
		t.Error("tree's self is equal to root location", 
			"\nexpected: ", expectedLoc,
			"\ncomputed: ", child.GetLoc())
	}

	expectedChildBucket := Loc{rootId, pId, cId, cId}

	if !child.GetChildBucket().Equal(expectedChildBucket) {
		t.Error("child has wrong child prefix",
			"\nexpected: ", expectedChildBucket,
			"\ncomputed: ", child.GetChildBucket())	
	}

	expectedDescendantBucket := expectedChildBucket[:3]

	if !child.GetDescendantBucket().Equal(expectedDescendantBucket) {
		t.Error("child has wrong descendant prefix",
			"\nexpected: ", expectedDescendantBucket,
			"\ncomputed: ", child.GetDescendantBucket())	
	}

	expectedSiblingBucket := expectedLoc[:3]

	if !child.GetSiblingBucket().Equal(expectedSiblingBucket) {
		t.Error("child has wrong sibling prefix",
			"\nexpected: ", expectedSiblingBucket,
			"\ncomputed: ", child.GetSiblingBucket())		
	}

	if !child.GetParentLoc().Equal(parent.GetLoc()) {
		t.Error("child's parent should have been equal to parent Location", 
			"\nchild: ", child.GetParentLoc(), 
			"\nparent: ", parent.GetLoc())
	}
}

func TestMakeBranchKeyChainRoot (t *testing.T) {
	branch, err := Root.MakeChildBranch()

	if err != nil {
		t.Error("error making branch: ", err)
	}

	cId := branch.Id

	if cId.Height != 1 {
		t.Error("something went wrong making the tree's Id: ", cId)
	}

	expectedLoc := Loc{rootId, rootId, cId,}

	if !branch.GetLoc().Equal(expectedLoc) {
		t.Error("tree's self is equal to root location", 
			"\nexpected: ", expectedLoc,
			"\ncomputed: ", branch.GetLoc())
	}

	expectedChildBucket := Loc{rootId, cId}

	if !branch.GetChildBucket().Equal(expectedChildBucket) {
		t.Error("branch has wrong branch prefix",
			"\nexpected: ", expectedChildBucket,
			"\ncomputed: ", branch.GetChildBucket())	
	}

	expectedDescendantBucket := expectedChildBucket //not all descendants of a branch are in the same bucket

	if !branch.GetDescendantBucket().Equal(expectedDescendantBucket) {
		t.Error("branch has wrong descendant prefix",
			"\nexpected: ", expectedDescendantBucket,
			"\ncomputed: ", branch.GetDescendantBucket())	
	}

	expectedSiblingBucket := expectedLoc[:2]

	if !branch.GetSiblingBucket().Equal(expectedSiblingBucket) {
		t.Error("branch has wrong sibling prefix",
			"\nexpected: ", expectedSiblingBucket,
			"\ncomputed: ", branch.GetSiblingBucket())		
	}

	if !branch.GetParentLoc().Equal(rootLoc) {
		t.Error("branch's parent should have been equal to parent Location", 
			"\nchild: ", branch.GetParentLoc(), 
			"\nparent: ", rootLoc)
	}

}

func TestMakeBranchKeyChainOnTree (t *testing.T) {
	parent, err := Root.MakeChildTree()

	if err != nil {
		t.Error("error making parent: ", err)
	}

	pId := parent.Id

	child, err := parent.MakeChildBranch()

	if err != nil {
		t.Error("error making tree: ", err)
	}

	cId := child.Id

	if cId.Height != 2 {
		t.Error("something went wrong making the tree's Id: ", cId)
	}

	expectedLoc := Loc{rootId, pId, pId, cId,}

	if !child.GetLoc().Equal(expectedLoc) {
		t.Error("tree's self is equal to root location", 
			"\nexpected: ", expectedLoc,
			"\ncomputed: ", child.GetLoc())
	}

	expectedChildBucket := Loc{rootId, pId, cId}

	if !child.GetChildBucket().Equal(expectedChildBucket) {
		t.Error("child has wrong child prefix",
			"\nexpected: ", expectedChildBucket,
			"\ncomputed: ", child.GetChildBucket())	
	}

	expectedDescendantBucket := expectedChildBucket

	if !child.GetDescendantBucket().Equal(expectedDescendantBucket) {
		t.Error("child has wrong descendant prefix",
			"\nexpected: ", expectedDescendantBucket,
			"\ncomputed: ", child.GetDescendantBucket())	
	}

	expectedSiblingBucket := expectedLoc[:3]

	if !child.GetSiblingBucket().Equal(expectedSiblingBucket) {
		t.Error("child has wrong sibling prefix",
			"\nexpected: ", expectedSiblingBucket,
			"\ncomputed: ", child.GetSiblingBucket())		
	}

	if !child.GetParentLoc().Equal(parent.GetLoc()) {
		t.Error("child's parent should have been equal to parent Location", 
			"\nchild: ", child.GetParentLoc(), 
			"\nparent: ", parent.GetLoc())
	}

}

func TestMakeBranchKeyChainOnBranch (t *testing.T) {
	parent, err := Root.MakeChildBranch()

	if err != nil {
		t.Error("error making branch")
	}

	pId := parent.Id

	child, err := parent.MakeChildBranch()

	if err != nil {
		t.Error("error making branch")
	}

	cId := child.Id

	if cId.Height != 2 {
		t.Error("something went wrong making the tree's Id: ", cId)
	}

	expectedLoc := Loc{rootId, pId, cId,}

	if !child.GetLoc().Equal(expectedLoc) {
		t.Error("tree's self is equal to root location", 
			"\nexpected: ", expectedLoc,
			"\ncomputed: ", child.GetLoc())
	}

	expectedChildBucket := Loc{rootId, cId,}

	if !child.GetChildBucket().Equal(expectedChildBucket) {
		t.Error("child has wrong child prefix",
			"\nexpected: ", expectedChildBucket,
			"\ncomputed: ", child.GetChildBucket())	
	}

	expectedDescendantBucket := expectedChildBucket

	if !child.GetDescendantBucket().Equal(expectedDescendantBucket) {
		t.Error("child has wrong descendant prefix",
			"\nexpected: ", expectedDescendantBucket,
			"\ncomputed: ", child.GetDescendantBucket())	
	}

	expectedSiblingBucket := expectedLoc[:2]

	if !child.GetSiblingBucket().Equal(expectedSiblingBucket) {
		t.Error("child has wrong sibling prefix",
			"\nexpected: ", expectedSiblingBucket,
			"\ncomputed: ", child.GetSiblingBucket())		
	}

	if !child.GetParentLoc().Equal(parent.GetLoc()) {
		t.Error("child's parent should have been equal to parent Location", 
			"\nchild: ", child.GetParentLoc(), 
			"\nparent: ", parent.GetLoc())
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

	pId := parent.Id

	child, err := parent.MakeChildBranch()

	if err != nil {
		t.Error("error making branch")
	}

	cId := child.Id

	if cId.Height != 3 {
		t.Error("something went wrong making the tree's Id: ", cId)
	}

	expectedLoc := Loc{rootId, pId, cId,}

	if !child.GetLoc().Equal(expectedLoc) {
		t.Error("tree's self is equal to root location", 
			"\nexpected: ", expectedLoc,
			"\ncomputed: ", child.GetLoc())
	}

	expectedChildBucket := Loc{rootId, cId,}

	if !child.GetChildBucket().Equal(expectedChildBucket) {
		t.Error("child has wrong child prefix",
			"\nexpected: ", expectedChildBucket,
			"\ncomputed: ", child.GetChildBucket())	
	}

	expectedDescendantBucket := expectedChildBucket

	if !child.GetDescendantBucket().Equal(expectedDescendantBucket) {
		t.Error("child has wrong descendant prefix",
			"\nexpected: ", expectedDescendantBucket,
			"\ncomputed: ", child.GetDescendantBucket())	
	}

	expectedSiblingBucket := expectedLoc[:2]

	if !child.GetSiblingBucket().Equal(expectedSiblingBucket) {
		t.Error("child has wrong sibling prefix",
			"\nexpected: ", expectedSiblingBucket,
			"\ncomputed: ", child.GetSiblingBucket())		
	}

	if !child.GetParentLoc().Equal(parent.GetLoc()) {
		t.Error("child's parent should have been equal to parent Location", 
			"\nchild: ", child.GetParentLoc(), 
			"\nparent: ", parent.GetLoc())
	}

}

func TestMakeSiblingTree (t *testing.T) {
	tree1, err := Root.MakeChildTree()

	if err != nil {
		t.Error("error making tree1: ", err)
	}

	tId := tree1.Id

	tree2, err := tree1.MakeSibling()
	
	if err != nil {
		t.Error("error making sibling tree: ", err)
	}

	sId := tree2.Id

	if sId.Equal(tId) {
		t.Error("sibling has same id as original")
	}

	if sId.Height != 1 {
		t.Error("something went wrong making the tree's Id: ", sId)
	}

	expectedLoc := Loc{rootId, rootId, sId}

	if !tree2.GetLoc().Equal(expectedLoc) {
		t.Error("tree's self is equal to root location", 
			"\nexpected: ", expectedLoc,
			"\ncomputed: ", tree2.GetLoc())
	}

	expectedChildBucket := Loc{rootId, sId, sId}

	if !tree2.GetChildBucket().Equal(expectedChildBucket) {
		t.Error("child has wrong child prefix",
			"\nexpected: ", expectedChildBucket,
			"\ncomputed: ", tree2.GetChildBucket())	
	}

	expectedDescendantBucket := expectedChildBucket[:2]

	if !tree2.GetDescendantBucket().Equal(expectedDescendantBucket) {
		t.Error("child has wrong descendant prefix",
			"\nexpected: ", expectedDescendantBucket,
			"\ncomputed: ", tree2.GetDescendantBucket())	
	}

	expectedSiblingBucket := expectedLoc[:2]

	if !tree2.GetSiblingBucket().Equal(expectedSiblingBucket) {
		t.Error("child has wrong sibling prefix",
			"\nexpected: ", expectedSiblingBucket,
			"\ncomputed: ", tree2.GetSiblingBucket())		
	}

	if !tree2.GetParentLoc().Equal(Root.GetLoc()) {
		t.Error("child's parent should have been equal to parent Location", 
			"\nchild: ", tree2.GetParentLoc(),
			"\nparent: ", Root.GetLoc())
	}


}

func TestMakeSiblingBranch (t *testing.T) {
	branch1, err := Root.MakeChildBranch()

	if err != nil {
		t.Error("error making branch: ", err)
	}

	branch2, err := branch1.MakeChildBranch()

	if err != nil {
		t.Error("error making branch: ", err)
	}

	branch3, err := branch2.MakeChildBranch()

	if err != nil {
		t.Error("error making branch: ", err)
	}

	branch1S, err := branch1.MakeSibling()

	if err != nil {
		t.Error("error making sibling: ", err)
	}

	branch2S, err := branch2.MakeSibling()

	if err != nil {
		t.Error("error making sibling: ", err)
	}

	branch3S, err := branch3.MakeSibling()

	if err != nil {
		t.Error("error making sibling: ", err)
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

	treeIsTree := tree.IsTree

	if !treeIsTree {
		t.Error("tree is not a tree: ", tree)
	}

	branch, err := Root.MakeChildBranch()

	if err != nil {
		t.Error("error making branch: ", err)
	}

	branchIsTree := branch.IsTree

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

	if !bytes.Equal(child.SiblingKeyPrefix(), parent.ChildKeyPrefix()) {
		t.Error("child's sibling prefix is not parent's ChildKeyPrefix")
	}

	if !bytes.Equal(parent.SiblingKeyPrefix(), Root.ChildKeyPrefix()) {
		t.Error("parent's sibling prefix is not root's ChildKeyPrefix")
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