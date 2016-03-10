package levTree
/*
The DbFunnel module is meant to make writes take up less total time and to
ensure that consecutive updates of any document are 

One of the issues of goleveldb, is that if one transaction is open and you try
to open another then the you get an error instead of causing the thread to 
block.  To manage this problem I use a funnel.  This also allows for writes to
be periodically batch written to the db so that less total time is spent
writing and hence blocking reads.

One of the consequences of how this is implemented is that you should never
assume that an update that you just ran is actually available to you
through the provided read methods or is on the db.  All reads from the api go
to the database itself and bypass the funnel so that reads and writes don't
have to compete for access.  When an update is called for an node that is in
the funnel that update will be applied to that copy of the node in the funnel.
*/
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
	mutex sync.Mutex
	nodes map[string]*node
}

//time between write batches.  Will make this settable later
var waitBetweenWrites time.Duration = 1 * time.Second

//string of the filepath to leveldb
var dbPath string

//initializes the funnel.
func init () {
	funnel.nodes = make(map[string]*node)
}

//starts the funnel.  This will periodically write all entries from the funnel
//to disk and then clear the entries from the funnel.
func startFunnel () {
	for {
		time.Sleep(waitBetweenWrites)
		err := clearFunnel()
		if err != nil {
			//I need to check for db failure here and figure out how to 
			//handle it gracefully.
			fmt.Println("Error clearing funnel:", err)
		}
	}
}

//Takes all entries from the funnel and puts then in a batch object.
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

//Atomically writes a batch to disk it to disk.
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

//blocks funnel access, Writes all entries in the funnel to disk and then 
//resets the funnel.
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

//Sets up the rootNode.  
//This function is way to long.  I need to figure out how to break it apart.
func initializeRoot () error {
	root := makeRoot()

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

	if err != nil {
		fmt.Println("error serializing root: ", err)
		t.Discard()
		return nil
	}

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
func getNodeAt (l locateable) (*node, error) {
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

//Lock and unlock funnel outside of this function if used in concurrent context.
//This allows update functions to behave atomically, without requiring 
//rewriting all of the boilerplate of figuring out whether or not the node is
//already in the funnel.  It should not be used outside of this context.
func getNodeIntoFunnel (l locateable) (*node, error) {
	n, isInFunnel := funnel.nodes[l.KeyString()]

	if !isInFunnel {
		n, err := getNodeAt(l)

		if err != nil {
			fmt.Println("Error getting node: ", err)
			return nil, err
		}

		funnel.nodes[l.KeyString()] = n
	}
	return n, nil
}