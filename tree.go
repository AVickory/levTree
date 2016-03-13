package levTree

/*
The Tree module sets up the data structures required to represent tree nodes in
memory.  The basic building block is the Record which consists of a location
and some data (location basically being a key in the database).

Nodes contain a Record with thier own location and the data that's meant
to be stored on that Node.  They also contain a parent Record which contains
thier parent's location and some data describing thier parent.  Additionally
they have a Record for each of that Node's children.

Within the provided API, note that a locateable can be either a Record or a
Node.  This simplifies the api so that calling .UpdateNodeData on a Node's
parent Record is the same thing as calling it on the parent.

Within the database there are four kinds of nodes.  They are not different
types, again to simplify the API, but rather designate how thier children will
be bucketed relative to themselves.

Branch - A "normal" Node.  Children of this kind of Node will be in the same bucket
that this Node is.  Most nodes should be Branches since other types will cause
the keys in the database to get progressively longer.  a single branch should
generally hold all of the data that you'll need on a given db access (within
reason; leveldb will run slowly if your entries get to long).

Tree - The root Node of a single tree.  A tree is a Node who's children will use the
tree's key as thier bucket.  They allow for grouping of related entries in the
db and allow for sequential access of all of thier descendants at once.  Trees
are generally meant to be the children of either other trees or forests.

Forest - A tree Node that is attached to the root Node of the db.  If you're
coming from a SQL background, then you might think about forests as your
database's tables and a forest's child trees as sub tables.

Root - the bottom namespace of the database.  This Node is just a place to hold
Meta data about your forests.  Since you can't directly name your forests, you
can instead put identifying data in the root's child metadata records.
*/

import (
	"fmt"
)

//this type is just a lexical device to make it clear when an argument
//represents something that will be passed to an updater's .Update method.
type updateData interface{}

//Updates are assumed to be commutative.  Note that because of the interface
//type of the argument, this method will have to coerce it's argument.
type updater interface {
	Update(updateData) (updater, error)
}

//allows records and nodes to be handled by the same methods.  This ends up
//simplifying the API in levTree.go a lot since I only have to type assert once
//instead of needing lots of internal type conversions or copies of identical
//functions
type locateable interface {
	Key() []byte
	KeyString() string
	Update(updateData) (Record, error)
	GetData() updater
}

//a Record describes a location in the db.
type Record struct {
	Loc  location
	Data updater
}

func (r Record) Key() []byte {
	return r.Loc.Key()
}
func (r Record) KeyString() string {
	return r.Loc.KeyString()
}
func (r Record) Update(u updateData) (Record, error) {
	Data, err := r.Data.Update(u)
	r.Data = Data
	if err != nil {
		fmt.Println("error updating record")
		return r, err
	}
	return r, nil
}
func (r Record) GetData() updater {
	return r.Data
}


//one tree Node.  It is itself a Record and contains records, but it's
//important to note that the Record objects in the parent and children fields
//do not necesarily contain the same data that is at the corresponding
//locations in the db, but rather data describing the connection between those
//locations and this one.  this is meant to allow a user to traverse the tree
//without loading more nodes from the db than necesary.
type Node struct {
	Record

	Height int

	Parent Record

	ChildBucket location
	Children    map[string]Record //maps child locations to indices in children slice
}

func (n Node) Key() []byte {
	return n.Record.Key()
}
func (n Node) KeyString() string {
	return n.Record.KeyString()
}
func (n Node) Update(u updateData) (Record, error) {
	r, err := n.Record.Update(u)
	n.Record = r
	if err != nil {
		return n.Record, err
	}
	return n.Record, nil
}
func (n Node) GetData() updater {
	return n.Record.GetData()
}

func joinParentAndChild(parent Node, child Node, metaData updater) (Node, Node) {
	if parent.Children == nil {
		parent.Children = make(map[string]Record)
	}

	parent.Children[child.Loc.KeyString()] = Record{
		Loc:  child.Loc,
		Data: metaData,
	}
	child.Parent = Record{
		Loc:  parent.Loc,
		Data: metaData,
	}
	return parent, child
}

//Updates the parent metaData for the child Node and the child metaData for the parent Node.
func updateConnectionMeta(parent Node, child Node, u updateData) error {

	parentMeta, err := child.Parent.Update(u)
	child.Parent = parentMeta
	if err != nil {
		fmt.Println("error getting updating child Node's parent", err)
		return err
	}

	childMeta := parent.Children[child.KeyString()]

	childMeta.Data = parentMeta.Data

	parent.Children[child.KeyString()] = childMeta

	if err != nil {
		fmt.Println("error getting updating child Node's parent", err)
		return err
	}

	return err
}

//Creates a Node whose children will be in the same namespace as this branch.
func makeBranch(parent Node, metaData updater, data updater) (Node, error) {
	var newBranch Node
	bucket := parent.ChildBucket

	loc, err := bucket.getNewLoc()

	if err != nil {
		fmt.Println("error getting new location", err)
		return newBranch, err
	}

	newBranch = Node{
		Record: Record{
			Loc:  loc,
			Data: data,
		},
		Height:      parent.Height + 1,
		ChildBucket: bucket,
		Children:    make(map[string]Record),
	}

	parent, newBranch = joinParentAndChild(parent, newBranch, metaData)

	return newBranch, err
}

//creates a branch of the parent Node who's children will be in a different
//namespace than the new branch
func makeTree(parent Node, metaData updater, data updater) (Node, error) {
	newTree, err := makeBranch(parent, metaData, data)

	if err != nil {
		fmt.Println("Error creating template branch: ", err)
		return newTree, err
	}

	newTree.ChildBucket = newTree.Record.Loc

	return newTree, nil
}

func (n Node) IsTree() bool {
	//branch nodes put their children in the same bucket that they are in while
	//trees put their children in a different bucket (currently tree children
	//have their namespace set to the id of the tree Node, but this may change
	//in the future when I start optimizing for sequential reads through trees)
	return !n.Record.Loc.getBucketLocation().equals(n.ChildBucket)
}

//creates a tree at height 0 attached to the root.  root should be the root of the db.  You could,
//but shouldn't, pass any other Node to this function
func makeForest(root Node, metaData updater, data updater) (Node, error) {
	newForest, err := makeTree(root, metaData, data)

	if err != nil {
		fmt.Println("Error creating template tree: ", err)
		return newForest, err
	}

	newForest.Height = 0

	return newForest, nil
}

func (n Node) IsForest() bool {
	//the only special characteristic of a forest is that it's Height is 0.
	//it's worth noting that this will return true for the root as well as
	//forests.
	return n.Height == 0 && n.IsTree()
}

func makeRoot() Node {
	return Node{
		Record: Record{
			Loc: noNameSpace,
		},
		ChildBucket: noNameSpace,
		Children:    make(map[string]Record),
	}
}
