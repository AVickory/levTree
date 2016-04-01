package levTree

import (
	"fmt"
	"testing"
	"time"
	"bytes"
	// "github.com/AVickory/levTree/keyChain"
)

func nodeTest(t *testing.T, data []byte, kc locateable) Node {
	err := clearFunnel()

	if err != nil {
		t.Error("error clearing funnel: ", err)
	}

	n, err := Get(kc)

	if err != nil {
		t.Error("error getting node from db: ", err)
	}

	if !n.GetLoc().Equal(kc.GetLoc()) {
		t.Error("node was saved with wrong location", 
			"\n expected: ", kc,
			"\n found: ", n.KeyChain)
	}

	if !bytes.Equal(n.Data, data) {
		t.Error("node was saved with wrong data",
			"\nexpected: ", data,
			"\nfound: ", n.Data)
	}

	return n
}

func forestTest (t *testing.T, data []byte) Node {
	forestKc, err := NewForest(data)

	if err != nil {
		t.Error("error making new forest: ", err)
	}

	return nodeTest(t, data, forestKc)
}

func treeTest (t *testing.T, parent Node, data []byte) Node {
	treeKc, err := NewTree(parent, data)

	if err != nil {
		t.Error("error making new tree: ", err)
	}

	return nodeTest(t, data, treeKc)
}

func branchTest (t *testing.T, parent Node, data []byte) Node {
	branchKc, err := NewBranch(parent, data)

	if err != nil {
		t.Error("error making new branch: ", err)
	}

	return nodeTest(t, data, branchKc)
}

func TestNew (t *testing.T) {
	err := initForSynchronousTests()

	if err != nil {
		t.Error("error initializing db: ", err)
	}

	f0 := forestTest(t, []byte{0})

	_ = forestTest(t, []byte{1})

	t0 := treeTest(t, f0, []byte{2})

	_ = treeTest(t, f0, []byte{3})

	b0 := branchTest(t, f0, []byte{4})

	_ = branchTest(t, f0, []byte{5})

	_ = branchTest(t, t0, []byte{6})

	_ = branchTest(t, t0, []byte{7})

	_ = branchTest(t, b0, []byte{8})

	_ = branchTest(t, b0, []byte{9})

}

func getParentTest (t *testing.T, parent Node, child Node) {
	foundParent, err := GetParent(child)

	if err != nil {
		t.Error("error getting parent: ", err)
	}

	if !bytes.Equal(foundParent.Data, parent.Data) {
		t.Error("parent data was not on the node retrieved from database",
			"\nexpected: ", parent.Data,
			"\nfound: ", foundParent.Data)
	}
}

func rangeSearchTest(t *testing.T, parent Node, nodes []Node, numExpected int) {
	if len(nodes) != numExpected {
		nodesData := make([]byte, 0, len(nodes))
		for _, v := range nodes {
			nodesData = append(nodesData, v.Data[0])
		}
		t.Error("wrong number of nodes",
			"\nparent data: ", parent.Data[0],
			"\nexpected: ", numExpected,
			"\nfound: ", len(nodes),
			"\nfound data: ", nodesData)
	}
}

func getChildrenTest (t *testing.T, parent Node, numChildren int) {
	children, err := GetChildren(parent)

	if err != nil {
		t.Error("error getting children: ", err)
	}

	rangeSearchTest(t, parent, children, numChildren)
}

func getSiblingsTest (t *testing.T, n Node, numSiblings int) {
	siblings, err := GetSiblings(n)

	if err != nil {
		t.Error("error getting siblings: ", err)
	}

	rangeSearchTest(t, n, siblings, numSiblings)
}

func getForestsTest(t *testing.T, numForests int) {
	forests, err := GetForests()

	if err != nil {
		t.Error("error getting forests: ", err)
	}

	rangeSearchTest(t, rootNode, forests, numForests)
}

//This test will fail until I get the immediate children search set up
func TestGet (t *testing.T) {
	err := initForSynchronousTests()

	if err != nil {
		t.Error("error initializing db: ", err)
	}

	f0 := forestTest(t, []byte{0})

	_ = forestTest(t, []byte{1})

	t0 := treeTest(t, f0, []byte{2})

	_ = treeTest(t, f0, []byte{3})

	b0 := branchTest(t, f0, []byte{4})

	_ = branchTest(t, f0, []byte{5})

	_ = branchTest(t, t0, []byte{6})

	_ = branchTest(t, t0, []byte{7})

	b1 := branchTest(t, b0, []byte{8})

	_ = branchTest(t, b0, []byte{9})

	b2 := branchTest(t, b1, []byte{8})

	_ = branchTest(t, b1, []byte{9})

	getParentTest(t, f0, t0)

	getParentTest(t, f0, b0)

	getParentTest(t, b0, b1)

	getParentTest(t, b1, b2)

	//trees may be one to high after keyChain update
	getChildrenTest(t, rootNode, 2) //2 immediate

	getChildrenTest(t, f0, 4) //4 immediate

	getChildrenTest(t, t0, 2) //2 immediate

	getChildrenTest(t, b0, 2) //2 immediate

	getChildrenTest(t, b1, 2) //2 immediate

	getChildrenTest(t, b2, 0) //no immediate


	getSiblingsTest(t, f0, 1) // this really should be 1

	getSiblingsTest(t, t0, 4) // should be 4.  finds parent if parent is tree and finds descendants of all siblings

	getSiblingsTest(t, b0, 4) // should be 4.  

	getSiblingsTest(t, b1, 2) // self and sibling

	getSiblingsTest(t, b2, 2) // self and sibling


	getForestsTest(t, 12)

}

func TestUpdate (t *testing.T) {
	err := initForSynchronousTests()

	if err != nil {
		t.Error("error initializing db: ", err)
	}

	f0 := forestTest(t, []byte{0})

	_ = forestTest(t, []byte{1})

	t0 := treeTest(t, f0, []byte{2})

	_ = treeTest(t, f0, []byte{3})

	_ = branchTest(t, f0, []byte{4}) //this used to be b0

	_ = branchTest(t, f0, []byte{5})

	b0 := branchTest(t, t0, []byte{6}) //now this one is b0.  this is to ensure consistent ordering.

	_ = branchTest(t, t0, []byte{7})

	b1 := branchTest(t, b0, []byte{8})

	_ = branchTest(t, b0, []byte{9})

	b2 := branchTest(t, b1, []byte{8})

	_ = branchTest(t, b1, []byte{9})


	initialNodes := []locateable{f0, t0, b0, b1, b2}

	nodesToUpdate, err := OpenUpdate(initialNodes...)

	if err != nil {
		t.Error("error opening update: ", err)
	}

	if len(nodesToUpdate) != len(initialNodes) {
		t.Error("open update did not return the right number of nodes: ", len(nodesToUpdate))
	}

	for i, v := range initialNodes {
		initial, _ := v.(Node)
		toUpdate := nodesToUpdate[i]
		if toUpdate.Data[0] != initial.Data[0] {
			t.Error("something went wrong getting the node/putting it in the funnel",
				"\nexpected: ", initial.Data[0],
				"\nfound: ", toUpdate.Data[0])
		}
		toUpdate.Data[0]++
		if toUpdate.Data[0] == initial.Data[0] {
			t.Error("changing the updateable should not change the original (because the updateable should be either the version already in the funnel or from the db.)",
				"\nexpected: ", initial.Data[0] + 1,
				"\nfound: ", toUpdate.Data[0])
		}
	}

	CloseUpdate(nodesToUpdate...)

	InitDb("./data/db", 10 * time.Millisecond)

	time.Sleep(waitBetweenWrites * 2)

	if err != nil {
		t.Error("error clearing funnel: ", err)
	}

	for i, v := range initialNodes {
		updatedNode, err := Get(v)

		if err != nil {
			t.Error("error getting node.  original data was: ", v.(Node).Data)
		}

		if updatedNode.Data[0] != nodesToUpdate[i].Data[0] {
			t.Error("db was not updataed for node with original data: ", v.(Node).Data)
		}
	}

	fmt.Println("herpderp")


}
