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
	"github.com/AVickory/levTree/keyChain"
)

//a Record describes a location in the db.
type Node struct {
	keyChain.KeyChain
	Data []byte
}

//Creates a Node whose children will be in the same namespace as this branch.
func (parent Node) makeBranch(data []byte) (Node, error) {
	var newBranch Node

	kc, err := parent.MakeChildBranch()

	if err != nil {
		fmt.Println("error getting new location", err)
		return newBranch, err
	}

	newBranch = Node{
		KeyChain:  kc,
		Data: data,
	}

	return newBranch, err
}

//creates a branch of the parent Node who's children will be in a different
//namespace than the new branch
func (parent Node) makeTree(data []byte) (Node, error) {
	var newTree Node

	kc, err := parent.MakeChildTree()

	if err != nil {
		fmt.Println("error getting new location", err)
		return newTree, err
	}

	newTree = Node{
		KeyChain:  kc,
		Data: data,
	}


	return newTree, nil
}

//branch nodes put their children in the same bucket that they are in while
//trees put their children in a different bucket (currently tree children
//have their namespace set to the id of the tree Node, but this may change
//in the future when I start optimizing for sequential reads through trees)
// func (n Node) IsTree() bool {
// 	return !n.Record.Loc.getBucketLocation().equals(n.ChildBucket)
// }

//creates a tree at height 0 attached to the root.  root should be the root of the db.  You could,
//but shouldn't, pass any other Node to this function
func makeForest(data []byte) (Node, error) {
	newForest, err := rootNode.makeTree(data)

	if err != nil {
		fmt.Println("Error creating template tree: ", err)
		return newForest, err
	}

	return newForest, nil
}

//the only special characteristic of a forest is that it's Height is 0.
//it's worth noting that this will return true for the root as well as
//forests.
// func (n Node) IsForest() bool {
// 	return n.Height == 0 && n.IsTree()
// }

func makeRoot() Node {
	return Node{
		keyChain.Root,
	}
}

func (n *Node) serialize() ([]byte, error) {
	var gobble bytes.Buffer
	enc := gob.NewEncoder(&gobble)
	err := enc.Encode(*n)

	if err != nil {
		fmt.Println("SERIALIZATION ERROR: ", err)
		return []byte{}, err
	}

	return gobble.Bytes(), nil
}

//fills the Record with deserialized data from the passed in gob
func (n *Node) deserialize(value []byte) (error) {
	// fmt.Println("value passed in: ", value)
	gobble := bytes.NewBuffer(value)
	// fmt.Println("gobble: ", gobble)
	dec := gob.NewDecoder(gobble)
	err := dec.Decode(n)

	if err != nil {
		fmt.Println("DESERIALIZATION ERROR: ", err)
		return err
	}

	return nil
}
