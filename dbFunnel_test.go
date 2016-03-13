package levTree

import (
	"fmt"
	"testing"
	"github.com/syndtr/goleveldb/leveldb"
	"os"
	"time"
)

func clearDb () error {
	err := os.RemoveAll(dbPath)

	if err != nil {
		fmt.Println("error clearing DB files")
		return err
	}

	waitBetweenWrites = 10 * time.Millisecond

	return nil
}

func initForSynchronousTests () error {
	dbPath = "./data/db"

	clearDb()

	err := initializeRoot()

	if err != nil {
		fmt.Println("database could not be initialized: ", err)
		return err
	}

	return nil
}

//I say sync, but what I actually mean is that this should not be used
//concurrently in general. (technically it ought to be fine for create 
//operations, but you could get some wackyness going on with concurrent
//updates)
func syncPut (n Node) error {
	db, err := leveldb.OpenFile(dbPath, nil)
	defer db.Close()

	if err != nil {
		fmt.Println("error opening database", err)
		return err
	}

	nSerial, err := serialize(n)
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

func TestRootManipulation (t *testing.T) {
	err := initForSynchronousTests()

	if err != nil {
		t.Error("db failed to initialize", err)
	}

	n, err := getNodeAt(rootRecord)

	if err != nil {
		t.Error("error getting root Node", err)
	}

	if !n.Record.Loc.equals(rootRecord.Loc) {
		t.Error("rootRecord was not saved with correct location")
	}

	n.Data = convertNumToUpdater(1)

	err = syncPut(n)

	if !n.Record.Loc.equals(rootRecord.Loc) {
		t.Error("error putting root Node back in db", err)
	}

	n, err = getNodeAt(rootRecord)

	if err != nil {
		t.Error("error getting root Node the second time", err)
	}

	storedData, ok := n.Data.(mockUpdateable)

	if !ok || storedData != 1 {
		t.Error("storedData was ", storedData, " not 1")
	}

	err = initializeRoot()

	if err != nil {
		t.Error("error attempting to reinitialize root: ", err)
	}

	n, err = getNodeAt(rootRecord)

	if err != nil {
		t.Error("error getting reinitialized root", err)
	}

	storedData, ok = n.Data.(mockUpdateable)

	if !ok || storedData != 1 {
		t.Error("storedData was ", storedData, " not 1")
	}

}

func TestFunnel (t *testing.T) {
	err := initForSynchronousTests()

	if err != nil {
		t.Error("error initializing db", err)
	}

	n1, err := getNodeAt(rootRecord)

	if err != nil {
		t.Error("error getting node")
	}

	if n1.Data == convertNumToUpdater(1) {
		t.Error("some how the root's data got initialized.  This Is A Bad Thing: ", err)
	}

	n1.Data = convertNumToUpdater(1)

	n2, err := getNodeIntoFunnel(n1)

	if err != nil {
		t.Error("error putting node into funnel")
	}

	if n2.Data == n1.Data {
		t.Error("since n1 isn't in the funnel n2 should have the original data for n1, but instead has: ", n2.Data)
	}

	n2.Data = n1.Data

	if funnel.nodes[n2.KeyString()].Data == n2.Data {
		//at some point I might have the funnel use pointers, but for the time
		//being It's important to some of the logic that changes to the node 
		//returned by getNodeIntoFunnel not change the copy of the node that is
		//in the funnel.
		t.Error("Changes to node should not change the copy in the funnel")
	}

	funnel.nodes[n2.KeyString()] = n2

	clearFunnel()

	n3, err := getNodeAt(n2)

	if err != nil {
		t.Error("error getting saved node: ", err)
	}

	if n3.Data != n2.Data {
		t.Error("funnel did not save data!")
	}
}


