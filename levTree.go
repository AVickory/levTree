/*
LevTree
*/
// An eventually consistent tree database implemented on top of leveldb.
/**/
// There's a lot of work left to go, between some API concerns (the update
// interface is really cumbersome) and a lot of areas for performance
// improvements (I really don't take advantage of leveldb's sequential seek
// capabilities).  Plans are in the works for both of these, in addition to a
// test suite (which is the next thing on my backlog)
/**/
// Oh, and a proper README.  For now, if you'd like to check out the api you'll
// have to use the comments or use godoc (it'll eventually be hosted, but
// getting it working comes first).
/*
tree.go
*/
// The Tree module sets up the data structures required to represent tree nodes in
// memory.  The basic building block is the Record which consists of a location
// and some data (location basically being a key in the database).
/**/
// nodes contain a Record with thier own location and the data that's meant
// to be stored on that Node.  They also contain a parent Record which contains
// thier parent's location and some data describing thier parent.  Additionally
// they have a Record for each of that Node's children.
/**/
// Within the provided API, note that a locateable can be either a Record or a
// Node.  This simplifies the api so that calling .UpdateNodeData on a Node's
// parent Record is the same thing as calling it on the parent.
/**/
// Within the database there are four kinds of nodes.  They are not different
// types, again to simplify the API, but rather designate how thier children will
// be bucketed relative to themselves.
/**/
// Branch - A "normal" Node.  Children of this kind of Node will be in the same bucket
// that this Node is.  Most nodes should be Branches since other types will cause
// the keys in the database to get progressively longer.  a single branch should
// generally hold all of the data that you'll need on a given db access (within
// reason; leveldb will run slowly if your entries get to long).
/**/
// Tree - The root Node of a single tree.  A tree is a Node who's children will use the
// tree's key as thier bucket.  They allow for grouping of related entries in the
// db and allow for sequential access of all of thier descendants at once.  Trees
// are generally meant to be the children of either other trees or forests.
/**/
// Forest - A tree Node that is attached to the root Node of the db.  If you're
// coming from a SQL background, then you might think about forests as your
// database's tables and a forest's child trees as sub tables.
/**/
// Root - the bottom namespace of the database.  This Node is just a place to hold
// Meta data about your forests.  Since you can't directly name your forests, you
// can instead put identifying data in the root's child metadata Records.
/*
dbFunnel.go
*/
// The DbFunnel module is meant to make writes take up less total time and to
// ensure that consecutive updates of any document are
/**/
// One of the issues of goleveldb, is that if one transaction is open and you try
// to open another then the you get an error instead of causing the thread to
// block.  To manage this problem I use a funnel.  This also allows for writes to
// be periodically batch written to the db so that less total time is spent
// writing and hence blocking reads.
/**/
// One of the consequences of how this is implemented is that you should never
// assume that an update that you just ran is actually available to you
// through the provided read methods or is on the db.  All reads from the api go
// to the database itself and bypass the funnel so that reads and writes don't
// have to compete for access.  When an update is called for an Node that is in
// the funnel that update will be applied to that copy of the Node in the funnel.
/*location.go*/
// The Location module provides a bucketing system for namespacing keys.  id
// generation defaults to guuid V4, which is sufficient for my usecase, but other
// options may be added later.
package levTree

import (
	"fmt"
	"time"
	"github.com/AVickory/levTree/keyChain"
)

// Should be run (and finish running) before any other operations on the db.
// Don't forget to register any types that you're storing using non-primitive
// types with gob.

// The only role of root record is to allow the look up of the root Node.
var rootNode Node = Node{
	KeyChain: keyChain.Root,
}

func InitDb(path string, writeInterval time.Duration) {
	dbPath = path

	waitBetweenWrites = writeInterval

	go startFunnel()

}

// a forest is a tree attached to the root Node whose key is the namespace for
// all of it's children.  Modifications to the returned forest cannot be
// persisted.
func NewForest(data []byte) (locateable, error) {

	newForest, err := makeForest(data)

	if err != nil {
		fmt.Println("error making forest Node: ", err)
		return nil, err
	}

	err = createNode(newForest)

	if err != nil {
		fmt.Println("error putting forest in db: ", err)
		return nil, err
	}

	// return newForest.GetLoc(), nil
	return newForest.KeyChain, nil
}

// Creates a child of the calling tree or forest in that tree's namespace, whose
// key is the namespace for all of it's children.  Modifications to the returned
// tree cannot be persisted.
func NewTree(parent locateable, data []byte) (locateable, error) {
	newTree, err := makeTree(parent, data)

	if err != nil {
		fmt.Println("error making tree Node: ", err)
		return nil, err
	}

	err = createNode(newTree)

	if err != nil {
		fmt.Println("error putting tree in db: ", err)
		return nil, err
	}

	// return newTree.GetLoc(), err
	return newTree.KeyChain, err
}

// Makes and persists (eventually) a child Node of the calling Node and updates
// the calling Node's children.  It doesn't return anything, because you won't
// be able to access it on the db until the funnel flushes. Modifications to the
// returned forest cannot be persisted.
func NewBranch(parent locateable, data []byte) (locateable, error) {
	newBranch, err := makeBranch(parent, data)

	if err != nil {
		fmt.Println("error making branch Node: ", err)
		return nil, err
	}

	err = createNode(newBranch)

	if err != nil {
		fmt.Println("error putting branch in db: ", err)
		return nil, err
	}

	// return newBranch.GetLoc(), err
	return newBranch.KeyChain, err
}

// func NewSiblingBranch(n locateable, data []byte) (keyChain.Loc, error) {

// }

// func NewSiblingTree(n locateable, data []byte) (keyChain.Loc, error) {

// }

// Gets the Node which the locateable describes (for instance, if called on a
// childmetaRecord, gets the actual child Node).
// Note that if you pass in a Node it will return the Node unchanged without
// looking it up in the database.  This is meant to ensure that within a thread,
// it's harder to have two copies of any given entry in the db.  if you want to
// bypass this behavior then you can pass in the Node's Record field instead of
// the Node.  This workaround should be used sparingly so that you don't run
// into consistency errors and avoid making more database queries than you need.
func Get(kc locateable) (Node, error) {
	n, err := getNode(kc.GetLoc())
	if err != nil {
		fmt.Println("error getting location's node", err)
		return n, err
	}

	return n, nil
}

// Gets the calling Node's parent.  Modifications to the returned Node cannot be
// persisted.
func GetParent(child locateable) (Node, error) {
	parent, err := getNode(child.GetParentLoc())

	if err != nil {
		fmt.Println("error getting parent Node", err)
		return parent, err
	}

	return parent, nil
}

// Gets all of the calling Node's children.  Generally it's better to use
// the meta version And load a subset of children based on the meta data stored in
// the Node.  Modifications to the returned nodes cannot be persisted.
func GetChildren(parent locateable) ([]Node, error) {
	children, err := getNodesFromBucket(parent.GetChildBucket())

	if err != nil {
		fmt.Println("error getting children nodes: ", err)
		return children, err
	}

	return children, err

}

//only works for trees right now.  using it on branch's is not reccomended
//since it requires a breadth first search and will be a lot slower.
func GetDescendants(parent locateable) ([]Node, error) {
	descendants, err := getNodesFromBucket(parent.GetDescendantBucket())

	if err != nil {
		fmt.Println("error getting descendants nodes: ", err)
		return descendants, err
	}

	return descendants, err
}

func GetSiblings(l locateable) ([]Node, error) {
	siblings, err := getNodesFromBucket(l.GetSiblingBucket())

	if err != nil {
		fmt.Println("error getting sibling nodes: ", err)
		return siblings, err
	}

	return siblings, nil
}


// Gets all forests in the db.
// Modifications to the returned forests cannot be persisted.
// DOES NOT CURRENTLY WORK.  Right now it loads all nodes in the
//db.  the GetImmediateChildren function is next on the feature list
//and will fix this
func GetForests() ([]Node, error) {
	forests, err := GetChildren(rootNode)

	if err != nil {
		fmt.Println("error getting forests: ", err)
		return forests, err
	}

	return forests, nil
}

// Returns the most up-to-date version of the node at the locations ls indicates.
// These nodes can be updated, but must be passed into Close Update for thos
// updates to take place (or for any updates to ever take place again).  It is
// intended for updates only.  If you want to do insertions use the new
// functions and if you only need to read, then use the get functions.
// note that changing the child and parent meta data on one node does not
// automatically change the corresponding data on the parent or child node.
// DO NOT MODIFY LOCATIONS.  if you do, you may end up with duplicates on the
// db.
// Eventually I'll set it up to only lock individual nodes and only put a read
// lock on the funnel, but for now this sets up the api and general
// functionality.
func OpenUpdate(kcs ...locateable) ([]Node, error) {
	funnel.mutex.Lock()

	updateableNodes := make([]Node, len(kcs))

	for i, kc := range kcs {
		updateableNode, err := getNodeUpdateable(kc.GetLoc())
		if err != nil {
			fmt.Println("error getting updateable Node", err)
			return updateableNodes, err
		}
		updateableNodes[i] = updateableNode
	}

	return updateableNodes, nil
}

func CloseUpdate(updatedNodes ...Node) {
	for _, n := range updatedNodes {
		funnel.nodes[n.KeyString()] = n
	}

	funnel.mutex.Unlock()
}
