package levTree
//An eventually consistent persistant tree database implemented on top of leveldb.

import (
	"fmt"
)

var rootRecord *record = &record{
	Loc: NoNameSpace,
}

//Should be run (and finish running) before any other operations on the db.
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

//Gets the node which the locateable describes (for instance if called on a child
//metaRecord, gets the actual child node, if called on a node just returns the node).
func Get (l locateable) (*node, error) {
	n, ok := l.(*node)

	if !ok {
		return getNode(l)
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
