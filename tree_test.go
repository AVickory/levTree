package levTree

import (
	"testing"
	"fmt"
	"encoding/gob"
)

type mockUpdateable int

func init() {
	var emptyMock mockUpdateable
	gob.Register(emptyMock)
	// gob.Register(location{})
	// gob.Register(Record{})
	// gob.Register(Node{})
}

func (number mockUpdateable) Update(u updateData) (updater, error) {
	number += u.(mockUpdateable)
	return number, nil
}

func testParentChildRel(t *testing.T, parent Node, child Node) {
	if !parent.Record.Loc.equals(child.Parent.Loc) {
		t.Error("CHILD'S PARENT LOC IS NOT THE PARENT'S LOC")
	}
	if !child.Record.Loc.equals(parent.Children[child.KeyString()].Loc) {
		t.Error("PARENT'S CHILD LOC IS NOT CHILD'S LOC")
	}
	if parent.Children[child.KeyString()].Data != child.Parent.Data {
		t.Error("PARENT'S AND CHILD'S METADATA ARE NOT EQUAL")
	}
	if !parent.ChildBucket.equals(child.Loc.getBucketLocation()) {
		t.Error("This Node was put in the wrong bucket!")
	}
}

func testNode(t *testing.T, parent Node, n Node, metaData updater, data updater) {
	if d := n.Data; d != data {
		t.Error("NODE HAS WRONG DATA: ", d)
	}
	if len(n.Children) != 0 {
		t.Error("NODE SHOULD NOT HAVE CHILDREN: ", len(n.Children))
	}
	if len(n.Key()) != len(parent.ChildBucket.Key())+16 {
		t.Error("NODE KEY SHOULD BE A GUUID: ", len(n.Key()))
	}
	if d := n.Parent.Data; d != metaData {
		t.Error("NODE PARENT HAS WRONG META DATA: ", d)
	}

	testParentChildRel(t, parent, n)
}

func testBranch(t *testing.T, parent Node, branch Node, metaData updater, data updater) {
	testNode(t, parent, branch, metaData, data)
	if branch.Height != parent.Height+1 {
		t.Error("branch should be at height + 1")
	}
	if !branch.ChildBucket.equals(parent.ChildBucket) {
		t.Error("branch will put children in wrong bucket!")
	}
}

func testTree(t *testing.T, parent Node, tree Node, metaData updater, data updater) {
	testNode(t, parent, tree, metaData, data)
	if tree.Height != parent.Height+1 {
		t.Error("tree should be at height + 1: ")
	}
	if tree.ChildBucket.equals(parent.Loc) {
		t.Error("tree will put children in wrong bucket!")
	}
}

func testForest(t *testing.T, root Node, forest Node, metaData updater, data updater) {
	testNode(t, root, forest, metaData, data)
	if forest.Height != 0 {
		t.Error("forest should be at height + 1: ")
	}
	if forest.ChildBucket.equals(root.Loc) {
		t.Error("forest will put children in wrong bucket!")
	}
}

func TestMakeForest(t *testing.T) {
	metaData := mockUpdateable(1)
	data := mockUpdateable(2)
	root := makeRoot()
	forest, err := makeForest(root, metaData, data)
	if err != nil {
		t.Error("ERROR MAKING FOREST: ", err)
	}
	// fmt.Println(t, root, forest, metaData, data)
	testForest(t, root, forest, metaData, data)
}

func TestMakeTwoForests(t *testing.T) {
	metaData1 := mockUpdateable(1)
	data1 := mockUpdateable(2)
	metaData2 := mockUpdateable(3)
	data2 := mockUpdateable(4)

	root := makeRoot()

	forest1, err := makeForest(root, metaData1, data1)
	if err != nil {
		t.Error("ERROR MAKING FOREST: ", err)
	}
	forest2, err := makeForest(root, metaData2, data2)
	if err != nil {
		t.Error("ERROR MAKING FOREST: ", err)
	}

	testForest(t, root, forest1, metaData1, data1)
	testForest(t, root, forest2, metaData2, data2)

}

func TestMakeTree(t *testing.T) {
	metaData := mockUpdateable(1)
	data := mockUpdateable(2)
	root := makeRoot()
	forest, err := makeForest(root, metaData, data)
	if err != nil {
		t.Error("ERROR MAKING FOREST: ", err)
	}
	metaDataTree := mockUpdateable(3)
	dataTree := mockUpdateable(4)
	tree, err := makeTree(forest, metaDataTree, dataTree)
	if err != nil {
		t.Error("ERROR MAKING TREE: ", err)
	}
	testTree(t, forest, tree, metaDataTree, dataTree)
}

func TestMakeTwoTrees(t *testing.T) {
	metaData := mockUpdateable(1)
	data := mockUpdateable(2)
	root := makeRoot()
	forest, err := makeForest(root, metaData, data)
	if err != nil {
		t.Error("ERROR MAKING FOREST: ", err)
	}
	metaDataTree1 := mockUpdateable(3)
	dataTree1 := mockUpdateable(4)
	tree1, err := makeTree(forest, metaDataTree1, dataTree1)
	if err != nil {
		t.Error("ERROR MAKING TREE: ", err)
	}
	metaDataTree2 := mockUpdateable(5)
	dataTree2 := mockUpdateable(6)
	tree2, err := makeTree(forest, metaDataTree2, dataTree2)
	if err != nil {
		t.Error("ERROR MAKING TREE: ", err)
	}
	testTree(t, forest, tree1, metaDataTree1, dataTree1)
	testTree(t, forest, tree2, metaDataTree2, dataTree2)
}

func TestMakeTreeOnTree(t *testing.T) {
	metaData := mockUpdateable(1)
	data := mockUpdateable(2)
	root := makeRoot()
	forest, err := makeForest(root, metaData, data)
	if err != nil {
		t.Error("ERROR MAKING FOREST: ", err)
	}
	metaDataTree1 := mockUpdateable(3)
	dataTree1 := mockUpdateable(4)
	tree1, err := makeTree(forest, metaDataTree1, dataTree1)
	if err != nil {
		t.Error("ERROR MAKING TREE: ", err)
	}
	metaDataTree2 := mockUpdateable(5)
	dataTree2 := mockUpdateable(6)
	tree2, err := makeTree(tree1, metaDataTree2, dataTree2)
	if err != nil {
		t.Error("ERROR MAKING TREE: ", err)
	}
	testTree(t, tree1, tree2, metaDataTree2, dataTree2)
}

func TestMakeBranch(t *testing.T) {
	metaData := mockUpdateable(1)
	data := mockUpdateable(2)
	root := makeRoot()
	forest, err := makeForest(root, metaData, data)
	if err != nil {
		t.Error("ERROR MAKING FOREST: ", err)
	}
	metaDataBranch := mockUpdateable(3)
	dataBranch := mockUpdateable(4)
	branch, err := makeBranch(forest, metaDataBranch, dataBranch)
	if err != nil {
		t.Error("ERROR MAKING BRANCH: ", err)
	}
	testBranch(t, forest, branch, metaDataBranch, dataBranch)
}

func TestMakeTwoBranches(t *testing.T) {
	metaData := mockUpdateable(1)
	data := mockUpdateable(2)
	root := makeRoot()
	forest, err := makeForest(root, metaData, data)
	if err != nil {
		t.Error("ERROR MAKING FOREST: ", err)
	}
	metaDataBranch1 := mockUpdateable(3)
	dataBranch1 := mockUpdateable(4)
	branch1, err := makeBranch(forest, metaDataBranch1, dataBranch1)
	if err != nil {
		t.Error("ERROR MAKING BRANCH: ", err)
	}
	metaDataBranch2 := mockUpdateable(5)
	dataBranch2 := mockUpdateable(6)
	branch2, err := makeBranch(forest, metaDataBranch2, dataBranch2)
	if err != nil {
		t.Error("ERROR MAKING BRANCH: ", err)
	}
	testBranch(t, forest, branch2, metaDataBranch2, dataBranch2)
	testBranch(t, forest, branch1, metaDataBranch1, dataBranch1)
}

func TestMakeBranchOnTree(t *testing.T) {
	metaData := mockUpdateable(1)
	data := mockUpdateable(2)
	root := makeRoot()
	forest, err := makeForest(root, metaData, data)
	if err != nil {
		t.Error("ERROR MAKING FOREST: ", err)
	}
	metaDataTree := mockUpdateable(3)
	dataTree := mockUpdateable(4)
	tree, err := makeTree(forest, metaDataTree, dataTree)
	if err != nil {
		t.Error("ERROR MAKING TREE: ", err)
	}
	metaDatabranch := mockUpdateable(5)
	databranch := mockUpdateable(6)
	branch, err := makeBranch(tree, metaDatabranch, databranch)
	if err != nil {
		t.Error("ERROR MAKING BRANCH: ", err)
	}
	testBranch(t, tree, branch, metaDatabranch, databranch)
}

func TestMakeBranchOnBranch(t *testing.T) {
	metaData := mockUpdateable(1)
	data := mockUpdateable(2)
	root := makeRoot()
	forest, err := makeForest(root, metaData, data)
	if err != nil {
		t.Error("ERROR MAKING FOREST: ", err)
	}
	metaDataBranch1 := mockUpdateable(3)
	dataBranch1 := mockUpdateable(4)
	branch1, err := makeBranch(forest, metaDataBranch1, dataBranch1)
	if err != nil {
		t.Error("ERROR MAKING BRANCH: ", err)
	}
	metaDataBranch2 := mockUpdateable(3)
	dataBranch2 := mockUpdateable(4)
	branch2, err := makeBranch(branch1, metaDataBranch2, dataBranch2)
	if err != nil {
		t.Error("ERROR MAKING BRANCH: ", err)
	}
	testBranch(t, branch1, branch2, metaDataBranch2, dataBranch2)
}


func testNodeEquality (n1 Node, n2 Node) bool {
	tests := []bool{
		n1.Record.Data == n2.Record.Data,
		n1.Record.Loc.equals(n2.Record.Loc),
		n1.Height == n2.Height,
		testRecordEquality(n1.Parent, n2.Parent),
		testRecordListEquality(n1.Children, n2.Children),
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
func testRecordEquality(r1, r2 Record) bool {
	return r1.Data == r2.Data && r1.Loc.equals(r2.Loc)
}
func testRecordListEquality(rs1, rs2 map[string]Record) bool {
	if len(rs1) != len(rs2) {
		fmt.Println("FAILED LIST EQUALITY")
		return false
	}
	for i, v := range rs1 {
		if !testRecordEquality(v, rs2[i]) {
			fmt.Println("FAILED RECORD EQUALITY")
			fmt.Println("key:", i)
			return false
		}
	}
	return true
}

func serializeDeserializeTest (t *testing.T, n Node) {
	gobble, err := serialize(n)
	if err != nil {
		t.Error("SERIALIZE ERROR")
	}
	newNode, err := deserialize(gobble)
	if err != nil {
		t.Error("DESERIALIZE ERROR", err)
	}
	if !testNodeEquality(n, newNode) {
		fmt.Println("ORIGINAL: ", n)
		fmt.Println("NEW: ", newNode)
		t.Error("DESERIALIZE DID NOT RETURN SERIALIZED NODE")
	}
}

func convertNumToUpdater(x int) updater {
	u := mockUpdateable(x)
	return u
}

func TestSerializeDeSerialize (t *testing.T) {
	root := makeRoot()
	forest, err := makeForest(root, convertNumToUpdater(-1), convertNumToUpdater(0))
	if err != nil {
		t.Error("ERROR MAKING FOREST: ", err)
	}
	tree, err := makeTree(forest, convertNumToUpdater(1), convertNumToUpdater(2))
	if err != nil {
		t.Error("GUUID ERROR", err)
	}
	n1, err := makeBranch(tree, convertNumToUpdater(3), convertNumToUpdater(4))
	if err != nil {
		t.Error("GUUID ERROR", err)
	}
	n2, err := makeBranch(tree, convertNumToUpdater(5), convertNumToUpdater(6))
	if err != nil {
		t.Error("GUUID ERROR", err)
	}
	n11, err := makeBranch(n1, convertNumToUpdater(7), convertNumToUpdater(8))
	if err != nil {
		t.Error("GUUID ERROR", err)
	}
	n12, err := makeBranch(n1, convertNumToUpdater(9), convertNumToUpdater(10))
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