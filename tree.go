//Sets up the data structures required to represent tree nodes in memory
package levTree

import (
	"encoding/gob"
	"bytes"
	"fmt"
	"sync"
)

//Updates are assumed to be commutative.  Note that because of the interface
//type of the argument, this method will have to coerce it's argument.
type updater interface {
	Update(interface{}) error
}

//a record is a location in the db and some data.
type record struct {
	Loc location
	Data updater
}

//one tree node.  It is itself a record and contains records, but it's
//important to note that the record objects in the parent and children fields
//do not necesarily contain the same data that is at the corresponding
//locations in the db, but rather data describing the connection between those
//locations and this one.  this is meant to allow a user to traverse the tree
//without loading more nodes from the db than necesary.
type node struct {
	record
	Loc location
	Data updater

	Height int

	Parent record
	Children map[string]record //maps child locations to indices in children slice
}

type tree node

type forest tree

func joinParentAndChild (parent *node, child *node, metaData updater) {
	parent.Children[child.Loc] = record{
		Loc: child.Loc,
		Data: metaData,
	}
	child.Parent = record{
		Loc: parent.Loc,
		Data: metaData,
	}
}

//Adds a child to the node and returns the new child node.
//It can, but is not intended to, be called on forest.  Use NewTree for that.
func (n *node) MakeNode (metaData updater, data updater) (*node, error) {

	var bucket location
	if n.Height == 0 { //n = 0 is the root of the tree
		bucket = n.Loc
	} else {
		bucket = n.Loc.getBucketLocation()
	}
	loc, err := bucket.getNewLoc()

	if err != nil {
		return nil, err
	}

	newNode := &node{
		Loc: loc,
		Data: data,
		Height: n.Height + 1,
		Children: make(map[string]int),
	}

	joinParentAndChild(n, newNode, metaData)

	return newNode, err
}



//Sets up the root nameSpace for all trees.  If this is not called then forest
//is set as a node at location NoNameSpace and has no data attached.
func makeForest (root *node, NameSpace []byte, metaData updater, data updater) *forest {

	loc := root.Loc.getNewLocWithId(&NameSpace)
	f := &node{
		Loc: loc,
		Data: data,
		Children: make(map[string]int),
		Height: -1,
	}
	joinParentAndChild(root, f)
	return f
}

//adds a tree to the forest namespace and returns the new root.
func MakeTree (parent *tree, metaData updater, data updater) (*node, error) {

	loc, err := forest.Loc.getNewLoc() //creates a new namespace inside forest for this tree

	if err != nil {
		return nil, err
	}

	rootNode := &node{
		Loc: loc,
		Data: data,
		Parent: record{
			Loc: forest.Loc,
			Data: metaData,
		},
		Height: parent.Height + 1,
		Children: make(map[string]int),
	}

	metaNode := record{
		Loc: loc,
		Data: metaData,
	}

	forest.Children[loc.KeyString()] = metaNode

	return rootNode, nil

}

//serializes the record into a gob and returns it as a byte slice
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

//A passthrough to abstract out the location module from modules importing this module
func (n *node) Key () []byte {
	return n.Loc.Key()
}

//A passthrough to abstract out the location module from modules importing this module
func (r *record) Key () []byte {
	return r.Loc.Key()
}

//registers required types with gob.  Any named types contained in a 
//record's data property must also be registered before serializing or 
//deserializing to or from the db.
func init () {
	gob.Register(location{})
	gob.Register(record{})
	gob.Register(node{})
	gob.Register(tree{})
	gob.Register(forest{})
}