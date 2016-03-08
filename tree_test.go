package levTree

import (
	"testing"
	"fmt"
)

func TestNewTree (t *testing.T) {
	tree, err := NewTree(1, 2)
	if err != nil {
		t.Error("GUUID ERROR", err)
	}
	if tree.Data != 2 {
		t.Error("TREE DATA WRONG")
	}
	if forest.Loc.equals(tree.Loc) {
		t.Error("TREE LOC WRONG")
	}
	if tree.Parent.Data != 1 {
		t.Error("TREE PARENT DATA WRONG")
	}
	if !(forest.Loc.equals(tree.Parent.Loc)) {
		t.Error("TREE PARENT LOC ERROR")
	}
	if tree.Height != 0 {
		t.Error("TREE HEIGHT WRONG")
	}
	if len(tree.Children) != 0 {
		t.Error("NEW TREE HAS CHILDREN")
	}
	if len(forest.Children) != 1 {
		t.Error("FOREST HAS WRONG NUMBER OF CHILDREN")
	}
	if forest.Children[0].Data != tree.Parent.Data {
		t.Error("FOREST CHILD DATA WRONG")
	}
	if !(forest.Children[0].Loc.equals(tree.Loc)) {
		t.Error("FOREST CHILD LOC WRONG")
	}
}

func TestNewChild (t *testing.T) {
	tree, err := NewTree(1,2)
	if err != nil {
		t.Error("GUUID ERROR", err)
	}

	n1, err := tree.NewChild(3, 4)
	n2, err := tree.NewChild(5, 6)
	if err != nil {
		t.Error("GUUID ERROR", err)
	}
	if n1.Data != 4 {
		t.Error("N1 DATA WRONG")
	}
	if tree.Loc.equals(n1.Loc) {
		t.Error("N1 LOC WRONG")
	}
	if n1.Parent.Data != 3 {
		t.Error("N1 PARENT DATA WRONG")
	}
	if !(tree.Loc.equals(n1.Parent.Loc)) {
		t.Error("N1 PARENT LOC ERROR")
	}
	if n1.Height != 1 {
		t.Error("N1 HEIGHT WRONG")
	}
	if len(n1.Children) != 0 {
		t.Error("N1 HAS CHILDREN")
	}
	if len(tree.Children) != 2 {
		t.Error("Tree HAS WRONG NUMBER OF CHILDREN")
	}
	if tree.Children[0].Data != n1.Parent.Data {
		t.Error("TREE CHILD DATA WRONG")
	}
	if !(tree.Children[0].Loc.equals(n1.Loc)) {
		t.Error("TREE CHILD LOC WRONG")
	}

	if err != nil {
		t.Error("GUUID ERROR", err)
	}
	if n2.Data != 6 {
		t.Error("N2 DATA WRONG")
	}
	if tree.Loc.equals(n2.Loc) {
		t.Error("N2 LOC WRONG")
	}
	if n2.Parent.Data != 5 {
		t.Error("N2 PARENT DATA WRONG")
	}
	if !(tree.Loc.equals(n2.Parent.Loc)) {
		t.Error("N2 PARENT LOC ERROR")
	}
	if n2.Height != 1 {
		t.Error("N2 HEIGHT WRONG")
	}
	if len(n2.Children) != 0 {
		t.Error("N2 HAS CHILDREN")
	}
	if tree.Children[1].Data != n2.Parent.Data {
		t.Error("TREE CHILD DATA WRONG")
	}
	if !(tree.Children[1].Loc.equals(n2.Loc)) {
		t.Error("TREE CHILD LOC WRONG")
	}

	n11, err := n1.NewChild(7, 8)
	if err != nil {
		t.Error("GUUID ERROR", err)
	}
	if n11.Data != 8 {
		t.Error("N11 DATA WRONG")
	}
	if n1.Loc.equals(n11.Loc) {
		t.Error("N11 LOC WRONG")
	}
	if n11.Parent.Data != 7 {
		t.Error("N11 PARENT DATA WRONG")
	}
	if !(n1.Loc.equals(n11.Parent.Loc)) {
		t.Error("N11 PARENT LOC ERROR")
	}
	if n11.Height != 2 {
		t.Error("N11 HEIGHT WRONG")
	}
	if len(n11.Children) != 0 {
		t.Error("N11 HAS CHILDREN")
	}
	if len(n1.Children) != 1 {
		t.Error("N1 HAS WRONG NUMBER OF CHILDREN")
	}
	if n1.Children[0].Data != n11.Parent.Data {
		t.Error("N1 CHILD DATA WRONG")
	}
	if !(n1.Children[0].Loc.equals(n11.Loc)) {
		t.Error("N1 CHILD LOC WRONG")
	}
}

//tests if two nodes contain the same data.  If you're checking if the nodes
//describe the same place in memory, then compare their locations with
//location.equals(location).
//data must be ==able

func testNodeEquality (n1, n2 *node) bool {
	tests := []bool{
		n1.Data == n2.Data,
		n1.Loc.equals(n2.Loc),
		n1.Height == n2.Height,
		testRecordEquality(n1.Parent, n2.Parent),
		testRecordListEquality(n1.Children, n2.Children),
		testMapEquality(n1.ChildrenMap, n2.ChildrenMap),
	}
	for _, v := range tests {
		if !v {
			return false
		}
	}
	return true
}
func testRecordEquality(r1, r2 record) bool {
	return r1.Data == r2.Data && r1.Loc.equals(r2.Loc)
}
func testRecordListEquality(rs1, rs2 []record) bool {
	if len(rs1) != len(rs2) {
		return false
	}
	for i, v := range rs1 {
		if !testRecordEquality(v, rs2[i]) {
			return false
		}
	}
	return true
}
func testMapEquality(m1, m2 map[string]int) bool {
	if len(m1) != len(m2) {
		return false
	}
	for i, v := range m1 {
		if v != m2[i] {
			return false
		}
	}
	return true
}

func serializeDeserializeTest (t *testing.T, n *node) {
	gobble, err := n.serialize()
	if err != nil {
		fmt.Println("SERIALIZE ERROR")
	}
	var newNode *node = &node{}
	err = newNode.deserialize(gobble)
	if err != nil {
		fmt.Println("DESERIALIZE ERROR")
	}
	if !testNodeEquality(n, newNode) {
		t.Error("DESERIALIZE DID NOT RETURN SERIALIZED NODE")
	}
}

func TestSerializeDeSerialize (t *testing.T) {
	tree, err := NewTree(1, 2)
	if err != nil {
		t.Error("GUUID ERROR", err)
	}
	n1, err := tree.NewChild(3,4)
	if err != nil {
		t.Error("GUUID ERROR", err)
	}
	n2, err := tree.NewChild(5,6)
	if err != nil {
		t.Error("GUUID ERROR", err)
	}
	n11, err := n1.NewChild(7,8)
	if err != nil {
		t.Error("GUUID ERROR", err)
	}
	n12, err := n1.NewChild(9,10)
	if err != nil {
		t.Error("GUUID ERROR", err)
	}
	if !testNodeEquality(n1, n1) {
		t.Error("NODE EQUALITY FUNCTION IS BROKEN")
	}
	serializeDeserializeTest(t, tree)
	serializeDeserializeTest(t, n1)
	serializeDeserializeTest(t, n2)
	serializeDeserializeTest(t, n11)
	serializeDeserializeTest(t, n12)
}

//this tests the above functions but for a specific namespace placed in
//forest.
func TestInitForest (t *testing.T) {
	fmt.Println("InitForest Tests Basic")
	InitForest([]byte("derp"), 1)
	if forest.Data != 1 {
		t.Error("FOREST DATA NOT SET")
	}
	if forest.Loc.KeyString() != "derp" {
		t.Error("FOREST NOT NAMESPACED")
	}
	TestNewTree(t)
	TestNewChild(t)

	fmt.Println("InitForest Tests Empty Name Space")
	InitForest([]byte{}, 2)
	if forest.Data != 2 {
		t.Error("FOREST DATA NOT SET")
	}
	if forest.Loc.KeyString() != "" {
		t.Error("FOREST NOT NAMESPACED")
	}
	TestNewTree(t)
	TestNewChild(t)
}