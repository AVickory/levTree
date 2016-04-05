package levTree

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"os"
	"testing"
	"time"
	"bytes"
)

func clearDb() error {
	err := os.RemoveAll(dbPath)

	if err != nil {
		fmt.Println("error clearing DB files")
		return err
	}

	return nil
}

func initForSynchronousTests() error {
	dbPath = "./data/db"
	waitBetweenWrites = 10 * time.Millisecond

	err := clearDb()

	if err != nil {
		fmt.Println("error clearing db: ", err)
		return err
	}

	return nil
}

//I say sync, but what I actually mean is that this should not be used
//concurrently in general. (technically it ought to be fine for create
//operations, but you could get some wackyness going on with concurrent
//updates)
func syncPut(n Node) error {
	db, err := leveldb.OpenFile(dbPath, nil)
	defer db.Close()

	if err != nil {
		fmt.Println("error opening database", err)
		return err
	}

	nSerial, err := n.serialize()
	if err != nil {
		fmt.Println("error serializing node: ", err)
		return err
	}

	err = db.Put(n.Key(), nSerial, nil)

	if err != nil {
		fmt.Println("error putting root into db", err)
		return err
	}

	return nil
}

func TestFunnel(t *testing.T) {
	err := initForSynchronousTests()

	if err != nil {
		t.Error("error initializing db", err)
	}

	n1, err := makeForest([]byte{1})

	if err != nil {
		t.Error("error making node")
	}

	err = createNode(n1)

	if err != nil {
		t.Error("error entering node into db")
	}

	n1.Data = []byte{2}

	n2, err := getNodeUpdateable(n1)

	if err != nil {
		t.Error("error putting node into funnel")
	}

	if bytes.Equal(n2.Data, n1.Data) {
		t.Error("since n1 isn't in the funnel n2 should have the original data for n1, but instead has: ", n2.Data)
	}

	n2.Data = n1.Data

	if bytes.Equal(funnel.nodes[n2.KeyString()].Data, n2.Data) {
		//at some point I might have the funnel use pointers, but for the time
		//being It's important to some of the logic that changes to the node
		//returned by getNodeIntoFunnel not change the copy of the node that is
		//in the funnel.
		t.Error("Changes to node should not change the copy in the funnel")
	}

	funnel.nodes[n2.KeyString()] = n2

	err = clearFunnel()

	if err != nil {
		t.Error("error clearing funnel", err)
	}

	n3, err := getNode(n2)

	if err != nil {
		t.Error("error getting saved node: ", err)
	}

	if !bytes.Equal(n3.Data, n2.Data) {
		t.Error("funnel did not save data!", 
			"\nsent to funnel: ", n2.Data,
			"\nfrom db: ", n3.Data)
	}
}

func setUpChildSearch (t *testing.T) (map[string]Node) {
	nodes := make(map[string]Node)
	var err error

	nodes["forest"], err = makeForest([]byte{1})

	if err != nil {
		t.Error("error making forest: ", err)
	}

	nodes["tree1"], err = makeTree(nodes["forest"], []byte{2})

	if err != nil {
		t.Error("error making tree: ", err)
	}

	nodes["tree2"], err = makeTree(nodes["forest"], []byte{3})

	if err != nil {
		t.Error("error making tree: ", err)
	}

	nodes["tree11"], err = makeTree(nodes["tree1"], []byte{4})

	if err != nil {
		t.Error("error making tree: ", err)
	}

	nodes["branch1"], err = makeBranch(nodes["forest"], []byte{5})

	if err != nil {
		t.Error("error making tree: ", err)
	}

	nodes["branch2"], err = makeBranch(nodes["forest"], []byte{6})

	if err != nil {
		t.Error("error making tree: ", err)
	}

	nodes["branch11"], err = makeBranch(nodes["branch1"], []byte{7})

	if err != nil {
		t.Error("error making tree: ", err)
	}

	nodes["branch12"], err = makeBranch(nodes["branch1"], []byte{8})

	if err != nil {
		t.Error("error making tree: ", err)
	}

	nodes["branch111"], err = makeBranch(nodes["branch11"], []byte{9})

	for name, val := range nodes {
		// fmt.Println(name)
		err = createNode(val)
		if err != nil {
			t.Error("error saving node with name: ", name)
		}

		_, err = getNode(val.GetLoc())
		if err != nil {
			t.Error(name, " didn't make it to the db")
		}
	}

	return nodes
}

func copyMap (m map[string]Node) map[string]Node {
	newMap := make(map[string]Node, len(m))
	for k, v := range m {
		newMap[k] = v
	}
	return newMap

}

func MapValsNotInSlice (m map[string]Node, s []Node) map[string]Node {
	absentVals := copyMap(m)

	for name, val := range m {
		for _, v := range s {
			if val.KeyString() == v.KeyString() {
				delete(absentVals, name)
			}
		}
	}
	return absentVals
}

func checkNumChildrenAbsentFromSearch(t *testing.T, nodes map[string]Node, n Node, expected int) bool {
	children, err := getNodesFromBucket(n.GetChildBucket())

	if err != nil {
		t.Error("error getting node's children", err)
	}

	// printLocs(children)

	// absent := MapValsNotInSlice(nodes, children)

	if len(children) != expected {
		t.Error("getting the children of the root should have returned", expected, "entries, but returned:\n", len(children))
		return true
	}
	return false
}

func printLocs(nodes []Node) {
	for _, v := range nodes {
		fmt.Println(v.Id)
	}
}

func TestChildSearch (t *testing.T) {
	err := initForSynchronousTests()
	if err != nil {
		t.Error("error initializing db")
	}
	nodes := setUpChildSearch(t)

	// fmt.Println("rootNode: ")
	// t.Error("rootNode: ")
	_ = checkNumChildrenAbsentFromSearch(t, nodes, rootNode, len(nodes)) // all nodes plus
	// fmt.Println("\n\nforest:")
	// t.Error("forest:")
	_ = checkNumChildrenAbsentFromSearch(t, nodes, nodes["forest"], 4) //all but itself

	// fmt.Println("\n\ntree1:")
	// t.Error("tree1:")
	_ = checkNumChildrenAbsentFromSearch(t, nodes, nodes["tree1"], 1) //tree1 and tree11

	// fmt.Println("\n\ntree2:")
	// t.Error("tree2:")
	_ = checkNumChildrenAbsentFromSearch(t, nodes, nodes["tree2"], 0) //tree2

	// fmt.Println("\n\ntree11:")
	// t.Error("tree11:")
	_ = checkNumChildrenAbsentFromSearch(t, nodes, nodes["tree11"], 0) //tree11

	// fmt.Println("\n\nbranch2:")
	// t.Error("branch2:")
	_ = checkNumChildrenAbsentFromSearch(t, nodes, nodes["branch2"], 0)

	// fmt.Println("\n\nbranch1:")
	// t.Error("branch1:")
	_ = checkNumChildrenAbsentFromSearch(t, nodes, nodes["branch1"], 2) //branch11 and branch12

	// fmt.Println("\n\nbranch11:")
	// t.Error("branch11:")
	_ = checkNumChildrenAbsentFromSearch(t, nodes, nodes["branch11"], 1) //branch111

	// fmt.Println("\n\nbranch12:")
	// t.Error("branch12:")
	_ = checkNumChildrenAbsentFromSearch(t, nodes, nodes["branch12"], 0)

	_ = checkNumChildrenAbsentFromSearch(t, nodes, nodes["branch111"], 0)

	forestChildren, err := getNodesFromBucketUpdateable(nodes["forest"].GetChildBucket())

	if err != nil {
		t.Error("error getting updateable children", err)
	}

	for _, v := range forestChildren {
		child := v
		child.Data = append(child.Data, child.Data[0]+1)
		if bytes.Equal(funnel.nodes[child.KeyString()].Data, child.Data) {
			t.Error("funnel's data should not have changed when updateable node changed")
		}
		funnel.nodes[child.KeyString()] = child
	}

	err = clearFunnel()

	if err != nil {
		t.Error("error clearing funnel")
	}

	forestChildren2, err := getNodesFromBucket(nodes["forest"].GetChildBucket())

	for i, v := range forestChildren2 {
		forestChild := forestChildren[i].Data
		if !(v.Data[0] == forestChild[0] && len(v.Data) == 2) {
			t.Error("update didn't work correctly", 
				"\n original: ", forestChildren[i].Data,
				"\n updated", v.Data)
		}
	}
	
}

func TestBulkPut (t *testing.T) {
	err := initForSynchronousTests()

	if err != nil {
		t.Error("error initializing db: ", err)
	}

	nodes := make([]Node, 3, 3)

	nodes[0], err = makeForest([]byte{0})

	if err != nil {
		t.Error("error making forest")
	}

	nodes[1], err = makeTree(nodes[0], []byte{1})

	if err != nil {
		t.Error("error making forest")
	}

	nodes[2], err = makeTree(nodes[1], []byte{2})

	if err != nil {
		t.Error("error making forest")
	}

	bulkPut(nodes...)

	err = clearFunnel()

	if err != nil {
		t.Error("error clearing funnel", err)
	}

	savedNodes, err := getNodesFromBucket(nodes[0].GetChildBucket())

	for i, v := range savedNodes {
		if !bytes.Equal(v.Data, []byte{byte(i+1)}) {
			t.Error("node not saved correctly", 
				"\nindex: ", i,
				"\nnode: ", v.Data)
		}
	}


}