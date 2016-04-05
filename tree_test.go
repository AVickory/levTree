package levTree

import (
	"fmt"
	"testing"
	"bytes"
)

func testParentChildRel(t *testing.T, parent Node, child Node) {
	if !parent.GetLoc().Equal(child.GetParentLoc()) {
		t.Error("CHILD'S PARENT LOC IS NOT THE PARENT'S LOC")
	}
	if !parent.KeyChain.GetChildBucket().Equal(child.KeyChain.GetLoc()[:len(parent.KeyChain.GetChildBucket())]) {
		t.Error("This Node was put in the wrong bucket!")
	}
}

func testNode(t *testing.T, parent Node, n Node, data []byte) {
	if d := n.Data; !bytes.Equal(d, data) {
		t.Error("NODE HAS WRONG DATA: ", d)
	}
	if n.IsTree {
		if len(n.Key()) != len(parent.GetChildBucket().Key())+24 {
			t.Error("NODE KEY SHOULD BE A GUUID: ", len(n.Key()))
		}
	} else if n.ParentIsTree() {
		if len(n.Key()) != len(parent.GetChildBucket().Key())+48 {
			t.Error("NODE KEY SHOULD BE A GUUID: ", len(n.Key()))
		}
	} else {
		if len(n.Key()) != len(parent.GetChildBucket().Key())+24 {
			t.Error("NODE KEY SHOULD BE A GUUID: ", len(n.Key()))
		}
	}

	testParentChildRel(t, parent, n)
}

func testBranch(t *testing.T, parent Node, branch Node, data []byte) {
	testNode(t, parent, branch, data)
	if branch.Height != parent.Height+1 {
		t.Error("branch should be at height parent + 1 ")
	}
	if !branch.GetParentLoc().Equal(parent.GetLoc()) {
		t.Error("branch will put children in wrong bucket!")
	}
}

func testTree(t *testing.T, parent Node, tree Node, data []byte) {
	testNode(t, parent, tree, data)
	if tree.Height != parent.Height+1 {
		t.Error("tree should be at height parent + 1: ")
	}
	if tree.GetChildBucket().Equal(parent.GetLoc()) {
		t.Error("tree will put children in wrong bucket!")
	}
}

func testForest(t *testing.T, root Node, forest Node, data []byte) {
	testNode(t, root, forest, data)
	if forest.Height != 1 {
		t.Error("forest should be at height = 1: ")
	}
	if forest.GetChildBucket().Equal(root.GetLoc()) {
		t.Error("forest will put children in wrong bucket!")
	}
}

func TestMakeForest(t *testing.T) {
	data := []byte{2}
	forest, err := makeForest(data)
	if err != nil {
		t.Error("ERROR MAKING FOREST: ", err)
	}
	testForest(t, rootNode, forest, data)
}

func TestMakeTwoForests(t *testing.T) {
	data1 := []byte{2}
	data2 := []byte{4}


	forest1, err := makeForest(data1)
	if err != nil {
		t.Error("ERROR MAKING FOREST: ", err)
	}
	forest2, err := makeForest(data2)
	if err != nil {
		t.Error("ERROR MAKING FOREST: ", err)
	}

	testForest(t, rootNode, forest1, data1)
	testForest(t, rootNode, forest2, data2)

}

func TestMakeTree(t *testing.T) {
	data := []byte{2}
	forest, err := makeForest(data)
	if err != nil {
		t.Error("ERROR MAKING FOREST: ", err)
	}
	dataTree := []byte{4}
	tree, err := makeTree(forest, dataTree)
	if err != nil {
		t.Error("ERROR MAKING TREE: ", err)
	}
	testTree(t, forest, tree, dataTree)
}

func TestMakeTwoTrees(t *testing.T) {
	data := []byte{2}
	forest, err := makeForest(data)
	if err != nil {
		t.Error("ERROR MAKING FOREST: ", err)
	}
	dataTree1 := []byte{4}
	tree1, err := makeTree(forest, dataTree1)
	if err != nil {
		t.Error("ERROR MAKING TREE: ", err)
	}
	dataTree2 := []byte{6}
	tree2, err := makeTree(forest, dataTree2)
	if err != nil {
		t.Error("ERROR MAKING TREE: ", err)
	}
	testTree(t, forest, tree1, dataTree1)
	testTree(t, forest, tree2, dataTree2)
}

func TestMakeTreeOnTree(t *testing.T) {
	data := []byte{2}

	forest, err := makeForest(data)
	if err != nil {
		t.Error("ERROR MAKING FOREST: ", err)
	}
	dataTree1 := []byte{4}
	tree1, err := makeTree(forest, dataTree1)
	if err != nil {
		t.Error("ERROR MAKING TREE: ", err)
	}
	dataTree2 := []byte{6}
	tree2, err := makeTree(tree1, dataTree2)
	if err != nil {
		t.Error("ERROR MAKING TREE: ", err)
	}
	testTree(t, tree1, tree2, dataTree2)
}

func TestMakeBranch(t *testing.T) {
	data := []byte{2}
	forest, err := makeForest(data)
	if err != nil {
		t.Error("ERROR MAKING FOREST: ", err)
	}
	dataBranch := []byte{4}
	branch, err := makeBranch(forest, dataBranch)
	if err != nil {
		t.Error("ERROR MAKING BRANCH: ", err)
	}
	testBranch(t, forest, branch, dataBranch)
}

func TestMakeTwoBranches(t *testing.T) {
	data := []byte{2}
	forest, err := makeForest(data)
	if err != nil {
		t.Error("ERROR MAKING FOREST: ", err)
	}
	dataBranch1 := []byte{4}
	branch1, err := makeBranch(forest, dataBranch1)
	if err != nil {
		t.Error("ERROR MAKING BRANCH: ", err)
	}
	dataBranch2 := []byte{6}
	branch2, err := makeBranch(forest, dataBranch2)
	if err != nil {
		t.Error("ERROR MAKING BRANCH: ", err)
	}
	testBranch(t, forest, branch2, dataBranch2)
	testBranch(t, forest, branch1, dataBranch1)
}

func TestMakeBranchOnTree(t *testing.T) {
	data := []byte{2}
	forest, err := makeForest(data)
	if err != nil {
		t.Error("ERROR MAKING FOREST: ", err)
	}
	dataTree := []byte{4}
	tree, err := makeTree(forest, dataTree)
	if err != nil {
		t.Error("ERROR MAKING TREE: ", err)
	}
	databranch := []byte{6}
	branch, err := makeBranch(tree, databranch)
	if err != nil {
		t.Error("ERROR MAKING BRANCH: ", err)
	}
	testBranch(t, tree, branch, databranch)
}

func TestMakeBranchOnBranch(t *testing.T) {
	data := []byte{2}
	forest, err := makeForest(data)
	if err != nil {
		t.Error("ERROR MAKING FOREST: ", err)
	}
	dataBranch1 := []byte{4}
	branch1, err := makeBranch(forest, dataBranch1)
	if err != nil {
		t.Error("ERROR MAKING BRANCH: ", err)
	}
	dataBranch2 := []byte{4}
	branch2, err := makeBranch(branch1, dataBranch2)
	if err != nil {
		t.Error("ERROR MAKING BRANCH: ", err)
	}
	testBranch(t, branch1, branch2, dataBranch2)
}

func testNodeEquality(n1 Node, n2 Node) bool {
	tests := []bool{
		bytes.Equal(n1.Data, n2.Data),
		n1.Equal(n2.KeyChain),
		n1.Height == n2.Height,
		// testMapEquality(n1.ChildrenMap, n2.ChildrenMap),
	}
	for i, v := range tests {
		if !v {
			fmt.Println("FAILED NODE EQUALITY TEST AT INDEX ", i)
			return false
		}
	}
	return true
}

func serializeDeserializeTest(t *testing.T, n Node) {
	gobble, err := n.serialize()
	if err != nil {
		t.Error("SERIALIZE ERROR")
	}
	var newNode Node
	err = newNode.deserialize(gobble)
	if err != nil {
		t.Error("DESERIALIZE ERROR", err)
	}
	if !testNodeEquality(n, newNode) {
		fmt.Println("ORIGINAL: ", n)
		fmt.Println("NEW: ", newNode)
		t.Error("DESERIALIZE DID NOT RETURN SERIALIZED NODE")
	}
}

func TestSerializeDeSerialize(t *testing.T) {
	forest, err := makeForest([]byte{0})
	if err != nil {
		t.Error("ERROR MAKING FOREST: ", err)
	}
	tree, err := makeTree(forest, []byte{2})
	if err != nil {
		t.Error("GUUID ERROR", err)
	}
	n1, err := makeBranch(tree, []byte{4})
	if err != nil {
		t.Error("GUUID ERROR", err)
	}
	n2, err := makeBranch(tree, []byte{6})
	if err != nil {
		t.Error("GUUID ERROR", err)
	}
	n11, err := makeBranch(n1, []byte{8})
	if err != nil {
		t.Error("GUUID ERROR", err)
	}
	n12, err := makeBranch(n1, []byte{1})
	if err != nil {
		t.Error("GUUID ERROR", err)
	}
	if !testNodeEquality(n1, n1) {
		t.Error("NODE EQUALITY FUNCTION IS BROKEN")
	}
	// fmt.Println("derp", *tree)
	serializeDeserializeTest(t, tree)
	serializeDeserializeTest(t, n1)
	serializeDeserializeTest(t, n2)
	serializeDeserializeTest(t, n11)
	serializeDeserializeTest(t, n12)
}
