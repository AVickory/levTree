package levTree

import (
	"github.com/syndtr/goleveldb/leveldb"
	"fmt"
	"sync"
	"time"
)

//goleveldb transactions throw errors when a transaction is already open
//instead of blocking. The funnel is a blocking update cache so that writes can
//be batch written to the db.  This allows for less total time spent blocking
//reads during writing and
var funnel struct {
	mutex sync.RWMutex
	nodes map[string]*node
}

//time between write batches.  Will make this settable later
var waitBetweenWrites time.Duration = 1 * time.Second

//string of the filepath to leveldb
var dbPath string


func init () {
	funnel.nodes = make(map[string]*node)
}

func startFunnel () {
	for {
		time.Sleep(waitBetweenWrites)
		err := clearFunnel()
		if err != nil {
			fmt.Println("Error clearing funnel:", err)
		}
	}
}

func writeFunnelToBatch () *leveldb.Batch {
	batch := new(leveldb.Batch)
	for _, n := range funnel.nodes {
		nSerial, err := n.serialize()
		if err != nil {
			fmt.Println("error serializing node: ", err)
			fmt.Println(n.Loc.Key(), n.Data)
		} else {
			batch.Put(n.Key(), nSerial)
		}
	}
	funnel.nodes = make(map[string]*node)
	return batch
}

func transactionalBatch(batch *leveldb.Batch) error {
	db, err := leveldb.OpenFile(dbPath, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("error opening db: ", err)
		return err
	}

	t, err := db.OpenTransaction()
	if err != nil {
		fmt.Println("error creating transaction: ", err)
		t.Discard()
		return err
	}

	err = t.Write(batch, nil)
	if err != nil {
		fmt.Println("error writing to transaction: ", err)
		t.Discard()
		return err
	}

	err = t.Commit()
	if err != nil {
		fmt.Println("error comitting to transaction: ", err)
		return err
	}

	return nil
}

func clearFunnel() error {
	funnel.mutex.Lock()
	defer funnel.mutex.Unlock()

	if(len(funnel.nodes) != 0) {

		batch := writeFunnelToBatch()

		err := transactionalBatch(batch)
		if err != nil {
			fmt.Println("transaction err: ", err)
			return err
		}

	}

	return nil
}

func makeRoot () error {
	root := &node{
		Loc: NoNameSpace,
		Children: make(map[string]record),
	}
	db, err := leveldb.OpenFile(dbPath, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return err
	}
	t, err := db.OpenTransaction()
	if err != nil {
		fmt.Println("Error opening transaction: ", err)
		t.Discard()	
		return err	
	}
	rootInitialized, err := t.Has(root.Loc.Key(), nil)
	if err != nil {
		fmt.Println("Error checking for root: ", err)
		t.Discard()
		return err
	}
	if rootInitialized {
		t.Discard()
		return nil
	}
	rootSerial, err := root.serialize()
	err = t.Put(root.Loc.Key(), rootSerial, nil)
	if err != nil {
		fmt.Println("Error writing root to transaction: ", err)
		t.Discard()
		return err		
	}
	err = t.Commit()
	if err != nil {
		fmt.Println("Error commiting transaction: ", err)
		return err
	}
	return nil

}

//gets from the db.  Note that this will not necesarily be up to date if the
//funnle has not cleared updates into the db.
func getNode (l location) (*node, error) {
	var n *node
	db, err := leveldb.OpenFile(dbPath, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("error opening db", err)
		return n, nil
	}
	nSerial, err := db.Get(l.Key(), nil)
	if err != nil {
		fmt.Println("Error getting node from db: ", err)
		return nil, err
	}
	err = n.deserialize(nSerial)
	if err != nil {
		fmt.Println("Error deserializing node: ", err)
		return nil, err
	}
	return n, nil
}

//Lock and unlock funnel outside of this function.
//This allows update functions to behave atomically, without requiring 
//rewriting all of the boilerplate of figuring out whether or not the node is
//already in the funnel.  It should not be used outside of this context.
func getNodeIntoFunnel (l location) (*node, error) {
	n, isInFunnel := funnel.nodes[l.KeyString()]

	if !isInFunnel {
		n, err := getNode(l)
		if err != nil {
			fmt.Println("Error getting node: ", err)
			return nil, err
		}

		funnel.nodes[l.KeyString()] = n
	}
	return n, nil
}