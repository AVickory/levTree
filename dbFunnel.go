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
have to compete for access.  When an update is called for an Node that is in
the funnel that update will be applied to that copy of the Node in the funnel.
*/
import (
	"encoding/gob"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"sync"
	"time"
	"github.com/AVickory/levTree/keyChain"
)

//goleveldb transactions throw errors when a transaction is already open
//instead of blocking. The funnel is a blocking update cache so that writes can
//be batch written to the db.  This allows for less total time spent blocking
//reads during writing and
var funnel struct {
	mutex sync.Mutex
	nodes map[string]Node
}

//time between write batches.  Will make this settable later
var waitBetweenWrites time.Duration = 1 * time.Second

//string of the filepath to leveldb
var dbPath string

//initializes the funnel and registers package types with gob.  Any named types contained in a
//Record's data property must also be registered before serializing or
//deserializing to or from the db.
func init() {
	funnel.nodes = make(map[string]Node)
	gob.Register(keyChain.KeyChain{})
	gob.Register(keyChain.Id{})
	gob.Register(keyChain.Loc{})
	gob.Register(Node{})
}

//starts the funnel.  This will periodically write all entries from the funnel
//to disk and then clear the entries from the funnel.
func startFunnel() {
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
func writeFunnelToBatch() *leveldb.Batch {
	batch := new(leveldb.Batch)

	for _, n := range funnel.nodes {
		nSerial, err := n.serialize()

		if err != nil {
			fmt.Println("error serializing Node: ", err)
			fmt.Println(n.Key(), n.Data)

		} else {
			batch.Put(n.Key(), nSerial)
		}
	}

	funnel.nodes = make(map[string]Node)

	return batch
}

func writeBatch(batch *leveldb.Batch) error {
	db, err := leveldb.OpenFile(dbPath, nil)
	defer db.Close()

	if err != nil {
		fmt.Println("error opening db: ", err)
		return err
	}

	err = db.Write(batch, nil)

	if err != nil {
		fmt.Println("error writing batch", err)
		return err
	}

	return nil
} 

//blocks funnel access, Writes all entries in the funnel to disk and then
//resets the funnel.
func clearFunnel() error {
	funnel.mutex.Lock()
	defer funnel.mutex.Unlock()

	if len(funnel.nodes) != 0 {

		batch := writeFunnelToBatch()

		// err := transactionalBatch(batch)

		err := writeBatch(batch)

		if err != nil {
			fmt.Println("error clearing funnel ", err)
			return err
		}

	}

	return nil
}

//At somepoint the return from here and the funnel will be put into a trie, but
//for now I'm sticking with the basics.  Also this function is too long.
func getNodesFromBucket(bucket Keyor) ([]Node, error) { 
	db, err := leveldb.OpenFile(dbPath, nil)
	defer db.Close()

	if err != nil {
		fmt.Println("Error opening file: ", err)
		return nil, err
	}

	nodes := make([]Node, 0, 10)

	iter := db.NewIterator(util.BytesPrefix(bucket.Key()), nil)

	for iter.Next() {
		// nodes = append(nodes, Node{})

		nSerial := iter.Value()
		
		n := new(Node)
		err := n.deserialize(nSerial)

		if err != nil {
			fmt.Println("error deserializing record", 
				"\n\tkey: ", iter.Key(),
				"\n\tpayload: ", nSerial,
				"\n\terror: ", err) //should this be returned?
		} else {
			nodes = append(nodes, *n) //this is super inefficient.  I'll fix the resizing behavior later.
		}
	}

	iter.Release()

	err = iter.Error()

	if err != nil {
		fmt.Println("error in iterator: ", err)
		return nodes, err
	}

	return nodes, err
}

func getNodesFromBucketUpdateable(bucket Keyor) ([]Node, error) {
	dbNodes, err := getNodesFromBucket(bucket)
	if err != nil {
		fmt.Println("error getting nodes from bucket")
		return nil, err
	}

	for idx, node := range dbNodes {
		upToDateNode, isInFunnel := funnel.nodes[node.KeyString()]
		if isInFunnel {
			dbNodes[idx] = upToDateNode
		} else {
			funnel.nodes[node.KeyString()] = node
		}
	}

	return dbNodes, nil
}

//gets from the db.  Note that this will not necesarily be up to date if the
//funnle has not cleared updates into the db.
func getNode(l Keyor) (Node, error) {
	var n Node

	db, err := leveldb.OpenFile(dbPath, nil)
	defer db.Close()

	if err != nil {
		fmt.Println("error opening db", err)
		return n, nil
	}

	nSerial, err := db.Get(l.Key(), nil)

	if err != nil {
		fmt.Println("Error getting Node from db: ", err, 
			"\nnode Key: ", l.Key())
		return n, err
	}

	err = n.deserialize(nSerial)

	if err != nil {
		fmt.Println("Error deserializing Node: ", err)
		return n, err
	}

	return n, nil
}

//Lock and unlock funnel outside of this function if used in concurrent context.
//This allows update functions to behave atomically, without requiring
//rewriting all of the boilerplate of figuring out whether or not the Node is
//already in the funnel.  It should not be used outside of this context.
func getNodeUpdateable(l Keyor) (Node, error) {

	n, isInFunnel := funnel.nodes[l.KeyString()]

	if !isInFunnel {
		var err error
		n, err = getNode(l)

		if err != nil {
			fmt.Println("Error getting Node: ", err)
			return n, err
		}

		funnel.nodes[l.KeyString()] = n
	}
	return n, nil
}

func createNode(n Node) error {
	db, err := leveldb.OpenFile(dbPath, nil)
	defer db.Close()

	if err != nil {
		fmt.Println("error opening db: ", err)
		return err
	}

	nSerial, err := n.serialize()

	if err != nil {
		fmt.Println("error putting node: ", err)
		return err
	}

	err = db.Put(n.Key(), nSerial, nil)

	if err != nil {
		fmt.Println("error writing node to db: ", err)
		return err
	}

	return nil
}

func bulkPut(nodes ...Node) {
	for _, v := range nodes {
		funnel.nodes[v.KeyString()] = v
	}
}
