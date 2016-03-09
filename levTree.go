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
func NewForest (metaData updater, data updater) (*node, error) {
	funnel.mutex.Lock()
	defer funnel.mutex.Unlock()

	root, err := getNodeIntoFunnel(rootRecord)

	if err != nil {
		fmt.Println("error getting root into funnel: ", err)
		return nil, err
	}

	forest, err := makeForest(root, metaData, data)
	if err != nil {
		fmt.Println("error making forest node: ", err)
		return nil, err
	}

	funnel.nodes[forest.KeyString()] = forest

	return forest, nil
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

//Creates a child of the calling tree or forest in that tree's namespace, whose
//key is the namespace for all of it's children.  Modifications to the returned
//tree cannot be persisted.
func NewTree (parent locateable, metaData updater, data updater) (*node, error) {
	funnel.mutex.Lock()
	defer funnel.mutex.Unlock()

	updatedParent, err := getNodeIntoFunnel(parent)

	if err != nil {
		fmt.Println("error getting root into funnel: ", err)
		return nil, err
	}

	newTree, err := makeTree(updatedParent, metaData, data)
	if err != nil {
		fmt.Println("error making tree: ", err)
		return nil, err
	}

	funnel.nodes[newTree.Loc.KeyString()] = newTree

	return newTree, err
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

	return children, nil

}

//Gets the node which the locateable describes (for instance if called on a child
//metaRecord, gets the actual child node, if called on a node just returns the node).
func Get (l locateable) (*node, error) {
	n, ok := l.(*node)
	if !ok {
		return getNode(l)
	}
	return n, nil
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
func NewNode (parent locateable, metaData updater, data updater) (*node, error) {
	funnel.mutex.Lock()
	defer funnel.mutex.Unlock()

	updatedNode, err := getNodeIntoFunnel(parent)

	if err != nil {
		fmt.Println("error getting root into funnel: ", err)
		return nil, err
	}

	newNode, err := makeBranch(updatedNode, metaData, data)
	if err != nil {
		fmt.Println("error making node: ", err)
		return nil, err
	}

	funnel.nodes[newNode.Loc.KeyString()] = newNode

	return newNode, nil
}