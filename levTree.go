package levTree
//An eventually consistent persistant tree database implemented on top of leveldb.

import (
	"fmt"
)

var rootRecord *record = &record{
	Loc: NoNameSpace,
}

//Should be run (and finish running) before any other executing any other operations on the db.
func InitDb (path string) error {
	dbPath = path
	err := makeRoot()
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
func NewForest (nameSpace []byte, metaData updater, data updater) (*forest, error) {
	funnel.mutex.Lock()
	defer funnel.mutex.Unlock()

	root, err := getNodeIntoFunnel(rootRecord.Loc)

	if err != nil {
		fmt.Println("error getting root into funnel: ", err)
		return nil, err
	}

	f := makeForest(root, nameSpace, metaData, data)
	n := node(*f)
	funnel.nodes[f.Loc.KeyString()] = &n

	return f, nil
}

//Gets all metaRecords for forests in the db.  Modifications to the returned
//records cannot be persisted.
func GetForestsMeta () (map[string]record, error) {
	root, err := rootRecord.Get()
	if err != nil {
		fmt.Println("error getting root: ", err)
		return nil, err
	}
	return root.Children, nil
}

//Gets all forests in the db.  Generally it's better to use the meta version
//And load a subset of forests based on the meta data stored in the root. 
//Modifications to the returned forests cannot be persisted.
func GetForests () ([]*forest, error) {
	forestMap, err := GetForestsMeta()
	if err != nil {
		fmt.Println("error loading forests: ", err)
	}

	forests := make([]*forest, 0, len(forestMap))
	for key, forestRecord := range forestMap {
		f, err := forestRecord.GetForest()
		if err != nil {
			fmt.Println("error getting forest: ", err)
			fmt.Println("key = ", key)
		}
		forests = append(forests, f)
	}

	return forests, nil
}

//Creates a child of the calling tree or forest in that tree's namespace, whose
//key is the namespace for all of it's children.  Modifications to the returned
//tree cannot be persisted.
func (parent *tree) NewTree (metaData updater, data updater) (*tree, error) {
	funnel.mutex.Lock()
	defer funnel.mutex.Unlock()

	n, err := getNodeIntoFunnel(parent.Loc)
	t := tree(*n)
	updateableTree := &t

	if err != nil {
		fmt.Println("error getting root into funnel: ", err)
		return nil, err
	}

	newTree, err := makeTree(updateableTree, metaData, data)
	if err != nil {
		fmt.Println("error making tree: ", err)
		return nil, err
	}

	funnel.nodes[newTree.Loc.KeyString()] = n

	return newTree, err
}

//Gets the calling node's parent metaRecord.  Modifications to the returned forests cannot be 
//persisted.
func (r *record) GetParentMeta () (*record, error) {
 	n, err := r.Get()
	if err != nil {
		fmt.Println("error loading node: ", err)
		return nil, err
	}
	return &n.Parent, nil
}

//Gets the calling node's parent.  Modifications to the returned node cannot be
//persisted.
func (r *record) GetParent () (*node, error) {
	parentMeta, err := r.GetParentMeta()
	if err != nil {
		fmt.Println("error getting parent metadata: ", err)
		return nil, err
	}
	parent, err := parentMeta.Get()
	if err != nil {
		fmt.Println("error getting parent node", err)
		return nil, err
	}
	return parent, nil
}

//Gets the calling node's children's metaRecords.  Modifications to the
//returned records cannot be persisted.
func (r *record) GetChildrenMeta () (map[string]record, error) {
	n, err := r.Get()
	if err != nil {
		fmt.Println("error loading node: ", err)
		return nil, err
	}
	return n.Children, nil
}

//Gets all of the calling node's children.  Generally it's better to use
//the meta version And load a subset of children based on the meta data stored in
//the node.  Modifications to the returned nodes cannot be persisted.
func (r *record) GetChildren () ([]*node, error) {
	childrenMeta, err := r.GetChildrenMeta()
	if err != nil {
		fmt.Println("error loading node: ", err)
		return nil, err
	}

	children := make([]*node, 0, len(childrenMeta))
	for key, childRecord := range children {
		n, err := childRecord.Get()
		if err != nil {
			fmt.Println("error getting node: ", err)
			fmt.Println("key = ", key)
		}
		children = append(children, n)
	}

	return children, nil

}

//Gets the node which the record describes (for instance if called on a child
//metaRecord, gets the actual child node).
func (r *record) Get () (*node, error) {
	// n, ok := &node(*r)
	// if !ok {
		return getNode(r.Loc)
	// }
	// return n, nil
}

//Gets the node this record describes and converts it to a forest.  
//Modifications to the returned forest cannot be persisted.
func (r *record) GetForest () (*forest, error) {
	n, err := r.Get()
	f := forest(*n)
	return &f, err
}

//Gets the node this record describes and converts it to a forest.  
//Modifications to the returned forest cannot be persisted.
func (r *record) GetTree () (*tree, error) {
	n, err := r.Get()
	t := tree(*n)
	return &t, err
}

//Updates the node's internal data by calling .Update on the node's Data 
//property and persisting the change to the server (eventually)
// func (r *record) updateNodeData (updateData interface{}) err {

// }

//Updates the node's parent metaRecord and the node's parent's child metaRecord
//for this node by calling .Update on both of these metaRecords and persisting
//the change to the server (eventually)
// func (r *record) updateParentMeta (updateData interface{}) err {

// }

//Updates the node's child metaRecord and the node's child's parent metaRecord
//for this node by calling .Update on both of these metaRecords and persisting
//the change to the server (eventually).
// func (r *record) updateChildMeta (child *record, updateData interface{}) err {

// }

//Makes and persists (eventually) a child node of the calling node and updates
//the calling node's children.  It doesn't return anything, because you won't
//be able to access it on the db until the funnel flushes. Modifications to the
//returned forest cannot be persisted.
func (r *record) NewNode (metaData updater, data updater) (*node, error) {
	funnel.mutex.Lock()
	defer funnel.mutex.Unlock()

	updateableNode, err := getNodeIntoFunnel(r.Loc)

	if err != nil {
		fmt.Println("error getting root into funnel: ", err)
		return nil, err
	}

	newNode, err := makeNode(updateableNode, metaData, data)
	if err != nil {
		fmt.Println("error making node: ", err)
		return nil, err
	}

	funnel.nodes[newNode.Loc.KeyString()] = newNode

	return newNode, nil
}