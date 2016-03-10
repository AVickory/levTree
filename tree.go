package levTree
/*
The Tree module sets up the data structures required to represent tree nodes in
memory.  The basic building block is the record which consists of a location
and some data (location basically being a key in the database).

Nodes contain a record with thier own location and the data that's meant
to be stored on that node.  They also contain a parent record which contains
thier parent's location and some data describing thier parent.  Additionally
they have a record for each of that node's children.

Within the provided API, note that a locateable can be either a record or a 
node.  This simplifies the api so that calling .UpdateNodeData on a node's
parent record is the same thing as calling it on the parent.

Within the database there are four kinds of nodes.  They are not different
types, again to simplify the API, but rather designate how thier children will
be bucketed relative to themselves.

Branch - A "normal" node.  Children of this kind of node will be in the same bucket
that this node is.  Most nodes should be Branches since other types will cause
the keys in the database to get progressively longer.  a single branch should
generally hold all of the data that you'll need on a given db access (within
reason; leveldb will run slowly if your entries get to long).

Tree - The root node of a single tree.  A tree is a node who's children will use the
tree's key as thier bucket.  They allow for grouping of related entries in the
db and allow for sequential access of all of thier descendants at once.  Trees
are generally meant to be the children of either other trees or forests.

Forest - A tree node that is attached to the root node of the db.  If you're
coming from a SQL background, then you might think about forests as your 
database's tables and a forest's child trees as sub tables.

Root - the bottom namespace of the database.  This node is just a place to hold
Meta data about your forests.  Since you can't directly name your forests, you
can instead put identifying data in the root's child metadata records.
*/

import (
	"encoding/gob"
	"bytes"
	"fmt"
)



//this type is just a lexical device to make it clear when an argument 
//represents something that will be passed to an updater's .Update method.
type updateData interface{}

//Updates are assumed to be commutative.  Note that because of the interface
//type of the argument, this method will have to coerce it's argument.
type updater interface {
	Update(updateData) error
}

//allows records and nodes to be handled by the same methods.  This ends up
//simplifying the API in levTree.go a lot since I only have to type assert once
//instead of needing lots of internal type conversions or copies of identical
//functions
type locateable interface {
	Key() []byte
	KeyString() string
	Update(updateData) error
	GetData() updater
}

//a record describes a location in the db.
type record struct {
	Loc location
	Data updater
}

func (r *record) Key() []byte {
	return r.Loc.Key()
}
func (r *record) KeyString() string {
	return r.Loc.KeyString()
}
func (r *record) Update(u updateData) error {
	return r.Data.Update(u)
}
func (r *record) GetData() updater {
	return r.Data
}

//one tree node.  It is itself a record and contains records, but it's
//important to note that the record objects in the parent and children fields
//do not necesarily contain the same data that is at the corresponding
//locations in the db, but rather data describing the connection between those
//locations and this one.  this is meant to allow a user to traverse the tree
//without loading more nodes from the db than necesary.
type node struct {
	*record

	Height int

	Parent record

	ChildBucket location
	Children map[string]record //maps child locations to indices in children slice
}

var rootRecord *record = &record{
	Loc: noNameSpace,
}

func (n *node) Key() []byte {
	return n.record.Key()
}
func (n *node) KeyString() string {
	return n.record.KeyString()
}
func (n *node) Update(u updateData) error {
	return n.record.Update(u)
}
func (n *node) GetData() updater {
	return n.record.GetData()
}

func joinParentAndChild (parent *node, child *node, metaData updater) {
	parent.Children[child.Loc.KeyString()] = record{
		Loc: child.Loc,
		Data: metaData,
	}
	child.Parent = record{
		Loc: parent.Loc,
		Data: metaData,
	}
}

//Updates the parent metaData for the child Node and the child metaData for the parent Node.  
func updateConnectionMeta (parent *node, child *node, u updateData) error {

	err := child.Parent.Update(u)

	if err != nil {
		fmt.Println("error getting updating child node's parent", err)
		return err
	}

	parent.Children[child.KeyString()] = child.Parent

	if err != nil {
		fmt.Println("error getting updating child node's parent", err)
		return err
	}

	return err
}

//Creates a node whose children will be in the same namespace as this branch.
func makeBranch (parent *node, metaData updater, data updater) (*node, error) {

	bucket := parent.ChildBucket

	loc, err := bucket.getNewLoc()

	if err != nil {
		fmt.Println("error getting new location", err)
		return nil, err
	}

	newBranch := &node{
		record: &record{
			Loc: loc,
			Data: data,
		},
		Height: parent.Height + 1,
		ChildBucket: bucket,
		Children: make(map[string]record),
	}

	joinParentAndChild(parent, newBranch, metaData)

	return newBranch, err
}

//creates a branch of the parent node who's children will be in a different
//namespace than the new branch
func makeTree (parent *node, metaData updater, data updater) (*node, error) {
	newTree, err := makeBranch(parent, metaData, data)

	if err != nil {
		fmt.Println("Error creating template branch: ", err)
		return nil, err
	}

	newTree.ChildBucket = newTree.record.Loc

	return newTree, nil
}

func (n *node) IsTree() bool {
	//branch nodes put their children in the same bucket that they are in while
	//trees put their children in a different bucket (currently tree children
	//have their namespace set to the id of the tree node, but this may change
	//in the future when I start optimizing for sequential reads through trees)
	return !n.record.Loc.getBucketLocation().equals(n.ChildBucket)
}

//creates a tree at height 0
func makeForest (root *node, metaData updater, data updater) (*node, error) {
	newForest, err := makeTree(root, metaData, data)
	
	if err != nil {
		fmt.Println("Error creating template tree: ", err)
		return nil, err
	}

	newForest.Height = 0
	
	return newForest, nil
}


func makeRoot () *node {
	return &node{
		record: &record{
			Loc: noNameSpace,
		},
		ChildBucket: noNameSpace,
		Children: make(map[string]record),
	}
}

func (n *node) IsForest() bool {
	//the only special characteristic of a forest is that it's Height is 0.
	//it's worth noting that this will return true for the root as well as
	//forests.
	return n.Height == 0 && n.IsTree()
}

//serializes the node into a gob and returns it as a byte slice
func (n *node) serialize () ([]byte, error) {
	gobble := new(bytes.Buffer)
	enc := gob.NewEncoder(gobble)
	err := enc.Encode(n)

	if err != nil {
		fmt.Println("SERIALIZATION ERROR: ", err)
		return []byte{}, err
	}

	return gobble.Bytes(), nil
}

//fills the record with deserialized data from the passed in gob
func (n *node) deserialize(value []byte) error {
	gobble := bytes.NewBuffer(value)
	dec := gob.NewDecoder(gobble)
	err := dec.Decode(n)
	
	if err != nil {
		fmt.Println("DESERIALIZATION ERROR: ", err)
		return err
	}

	return nil
}

//registers required types with gob.  Any named types contained in a 
//record's data property must also be registered before serializing or 
//deserializing to or from the db.
func init () {
	gob.Register(location{})
	gob.Register(record{})
	gob.Register(node{})
}