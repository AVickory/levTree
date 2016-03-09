//Sets up the data structures required to represent tree nodes in memory
package levTree

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
//simplifying the API in levTree.go a lot since it gives me a lot freedom in
//when to use type conversion and
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