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
// memory.  The basic building block is the record which consists of a location
// and some data (location basically being a key in the database).
/**/
// Nodes contain a record with thier own location and the data that's meant
// to be stored on that node.  They also contain a parent record which contains
// thier parent's location and some data describing thier parent.  Additionally
// they have a record for each of that node's children.
/**/
// Within the provided API, note that a locateable can be either a record or a 
// node.  This simplifies the api so that calling .UpdateNodeData on a node's
// parent record is the same thing as calling it on the parent.
/**/
// Within the database there are four kinds of nodes.  They are not different
// types, again to simplify the API, but rather designate how thier children will
// be bucketed relative to themselves.
/**/
// Branch - A "normal" node.  Children of this kind of node will be in the same bucket
// that this node is.  Most nodes should be Branches since other types will cause
// the keys in the database to get progressively longer.  a single branch should
// generally hold all of the data that you'll need on a given db access (within
// reason; leveldb will run slowly if your entries get to long).
/**/
// Tree - The root node of a single tree.  A tree is a node who's children will use the
// tree's key as thier bucket.  They allow for grouping of related entries in the
// db and allow for sequential access of all of thier descendants at once.  Trees
// are generally meant to be the children of either other trees or forests.
/**/
// Forest - A tree node that is attached to the root node of the db.  If you're
// coming from a SQL background, then you might think about forests as your 
// database's tables and a forest's child trees as sub tables.
/**/
// Root - the bottom namespace of the database.  This node is just a place to hold
// Meta data about your forests.  Since you can't directly name your forests, you
// can instead put identifying data in the root's child metadata records.
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
// have to compete for access.  When an update is called for an node that is in
// the funnel that update will be applied to that copy of the node in the funnel.
/*location.go*/
// The Location module provides a bucketing system for namespacing keys.  id 
// generation defaults to guuid V4, which is sufficient for my usecase, but other
// options may be added later.
package levTree

import (
	"fmt"
	"time"
)

//Should be run (and finish running) before any other operations on the db.
//Don't forget to register any types that you're storing using non-primitive
//types with gob.
func InitDb (path string, writeInterval time.Duration) error {
	dbPath = path

	waitBetweenWrites = writeInterval

	err := initializeRoot()
	
	if err != nil {
		fmt.Println("database could not be initialized: ", err)
		return err
	}

	go startFunnel()

	return nil
}

//a forest is a tree attached to the root node whose key is the namespace for
//all of it's children.  Modifications to the returned forest cannot be
//persisted.
func NewForest (metaData updater, data updater) error {
	funnel.mutex.Lock()
	defer funnel.mutex.Unlock()

	root, err := getNodeIntoFunnel(rootRecord)

	if err != nil {
		fmt.Println("error getting root into funnel: ", err)
		return err
	}

	newForest, err := makeForest(root, metaData, data)

	if err != nil {
		fmt.Println("error making forest node: ", err)
		return err
	}

	funnel.nodes[newForest.KeyString()] = newForest

	return nil
}

//Creates a child of the calling tree or forest in that tree's namespace, whose
//key is the namespace for all of it's children.  Modifications to the returned
//tree cannot be persisted.
func NewTree (parent locateable, metaData updater, data updater) error {
	funnel.mutex.Lock()
	defer funnel.mutex.Unlock()

	updatedParent, err := getNodeIntoFunnel(parent)

	if err != nil {
		fmt.Println("error getting root into funnel: ", err)
		return err
	}

	newTree, err := makeTree(updatedParent, metaData, data)

	if err != nil {
		fmt.Println("error making tree: ", err)
		return err
	}

	funnel.nodes[newTree.Loc.KeyString()] = newTree

	return err
}

//Makes and persists (eventually) a child node of the calling node and updates
//the calling node's children.  It doesn't return anything, because you won't
//be able to access it on the db until the funnel flushes. Modifications to the
//returned forest cannot be persisted.
func NewBranch (parent locateable, metaData updater, data updater) error {
	funnel.mutex.Lock()
	defer funnel.mutex.Unlock()

	updatedNode, err := getNodeIntoFunnel(parent)

	if err != nil {
		fmt.Println("error getting root into funnel: ", err)
		return err
	}

	newNode, err := makeBranch(updatedNode, metaData, data)

	if err != nil {
		fmt.Println("error making node: ", err)
		return err
	}

	funnel.nodes[newNode.Loc.KeyString()] = newNode

	return nil
}

//Updates the node's internal data by calling .Update on the node's Data 
//property and persisting the change to the server (eventually).
//it doesn't return the updated node because that node will not reflect the
//version of the node that is on the db
func updateNodeData (l locateable, u updateData) error {
	funnel.mutex.Lock()
	defer funnel.mutex.Unlock()

	updatedNode, err := getNodeIntoFunnel(l)

	if err != nil {
		fmt.Println("error getting updateable node", err)
		return err
	}

	err = updatedNode.Update(u)

	if err != nil {
		fmt.Println("error updating node", err)
		return err
	}

	return nil
}

//Updates the node's parent metaRecord and the node's parent's child metaRecord
//for this node by calling .Update on both of these metaRecords and persisting
//the change to the server (eventually)
//it doesn't return the updated record because that node will not reflect the
//version of the node that is on the db
func updateParentMeta (child locateable, u updateData) error {
	funnel.mutex.Lock()
	defer funnel.mutex.Unlock()

	childNode, err := getNodeIntoFunnel(child)

	if err != nil {
		fmt.Println("error getting updateable child", err)
		return err
	}

	parentNode, err := getNodeIntoFunnel(&childNode.Parent)

	if err != nil {
		fmt.Println("error getting updateable parent", err)
		return err
	}

	err = updateConnectionMeta(parentNode, childNode, u)

	if err != nil {
		fmt.Println("Error updating meta data")
		return err
	}

	return nil
}

//Updates the node's child metaRecord and the node's child's parent metaRecord
//for this node by calling .Update on both of these metaRecords and persisting
//the change to the server (eventually).
//it doesn't return the updated record because that node will not reflect the
//version of the node that is on the db
func updateChildMeta (parent locateable, child locateable, u updateData) error {
	funnel.mutex.Lock()
	defer funnel.mutex.Unlock()

	parentNode, err := getNodeIntoFunnel(parent)

	if err != nil {
		fmt.Println("error getting updateable parent", err)
		return err
	}

	childNode, err := getNodeIntoFunnel(child)

	if err != nil {
		fmt.Println("error getting updateable child", err)
		return err
	}

	err = updateConnectionMeta(parentNode, childNode, u)

	if err != nil {
		fmt.Println("Error updating meta data")
		return err
	}

	return nil
}

//Gets the node which the locateable describes (for instance, if called on a
//childmetaRecord, gets the actual child node).
//Note that if you pass in a node it will return the node unchanged without
//looking it up in the database.  This is meant to ensure that within a thread,
//it's harder to have two copies of any given entry in the db.  if you want to
//bypass this behavior then you can pass in the node's record field instead of
//the node.  This workaround should be used sparingly so that you don't run
//into consistency errors and avoid making more database queries than you need.
func Get (l locateable) (*node, error) {
	n, ok := l.(*node)

	if !ok {
		return getNodeAt(l)
	}

	return n, nil
}

//Gets the calling node's parent metaRecord.  Modifications to the returned forests cannot be 
//persisted.
func GetParentMeta (child locateable) (*record, error) {

 	n, err := Get(child)
	
	if err != nil {
		fmt.Println("error loading node: ", err)
		return nil, err
	}
	
	return &n.Parent, nil
}

//Gets the calling node's parent.  Modifications to the returned node cannot be
//persisted.
func GetParent (child locateable) (*node, error) {
	parentMeta, err := GetParentMeta(child)

	if err != nil {
		fmt.Println("error getting parent metadata: ", err)
		return nil, err
	}

	parent, err := Get(parentMeta)

	if err != nil {
		fmt.Println("error getting parent node", err)
		return nil, err
	}

	return parent, nil
}

//Gets the calling node's children's metaRecords.  Modifications to the
//returned records cannot be persisted.
func GetChildrenMeta (parent locateable) (map[string]record, error) {
	n, err := Get(parent)

	if err != nil {
		fmt.Println("error loading node: ", err)
		return nil, err
	}

	return n.Children, nil
}

//Gets all of the calling node's children.  Generally it's better to use
//the meta version And load a subset of children based on the meta data stored in
//the node.  Modifications to the returned nodes cannot be persisted.
func GetChildren (parent locateable) ([]*node, error) {
	childrenMeta, err := GetChildrenMeta(parent)

	if err != nil {
		fmt.Println("error loading node: ", err)
		return nil, err
	}

	children := make([]*node, 0, len(childrenMeta))
	for key, childRecord := range childrenMeta {
		n, err := Get(&childRecord)

		if err != nil {
			fmt.Println("error getting node: ", err)
			fmt.Println("key = ", key)
		}

		children = append(children, n)
	}

	return children, err

}

//Gets all metaRecords for forests in the db.  Modifications to the returned
//records cannot be persisted.
func GetForestsMeta () (map[string]record, error) {
	forests, err := GetChildrenMeta(rootRecord)

	if err != nil {
		fmt.Println("error getting forests: ", err)
		return nil, err
	}

	return forests, nil
}

//Gets all forests in the db.  Generally it's better to use the meta version
//And load a subset of forests based on the meta data stored in the root. 
//Modifications to the returned forests cannot be persisted.
func GetForests () ([]*node, error) {
	forests, err := GetChildren(rootRecord)

	if err != nil {
		fmt.Println("error getting forests: ", err)
		return nil, err
	}

	return forests, nil
}
