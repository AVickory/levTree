package levTree

import (
	"testing"
	"fmt"
)

func TestForest (t *testing.T) {
	err := initForSynchronousTests()

	if err != nil {
		t.Error("error initializing db: ", err)
	}

	err = NewForest(convertNumToUpdater(1), convertNumToUpdater(2))

	if err != nil {
		t.Error("error writing forest (to funnel): ", err)
	}

	err = NewForest(convertNumToUpdater(3), convertNumToUpdater(4))

	if err != nil {
		t.Error("error writing forest (to funnel): ", err)
	}

	err = clearFunnel() //both forests are now in db

	if err != nil {
		t.Error("error writing forests (to disk): ", err)
	}

	forestsMeta, err := GetForestsMeta()

	if err != nil {
		t.Error("error getting forest meta from root: ", err)
	}

	forests, err := GetForests()

	if err != nil {
		t.Error("error getting forests", err)
	}

	if len(forestsMeta) != 2 {
		t.Error("too much meta data in root: ", len(forestsMeta))
	}

	if len(forests) != 2 {
		fmt.Println(forests)
		t.Error("too many forests: ", len(forests))
	}

	foundData := make(map[string]bool)

	for _, n := range forests {

		if n.Parent.Data != forestsMeta[n.KeyString()].Data {
			t.Error("Metadata was not stored differently in root and forest")
		}

		expectedData, _ := n.Parent.Data.(mockUpdateable)

		expectedData += 1

		if n.Data != convertNumToUpdater(int(expectedData)) {
			t.Error("Wrong data stored in record with parentMeta: ", n.Parent.Data)
		}

		foundData[n.KeyString()] = true

	}

	if len(foundData) != 2 {
		t.Error("Not all forests were saved!")
	}

	if forests[0].Data == forests[1].Data {
		t.Error("one forest was saved at both locations!")
	}

}

		// err = UpdateNodeData(n, 1)

		// if err != nil {
		// 	t.Error("error updating funnel node", err)
		// }

		// err = clearFunnel()

		// if err != nil {
		// 	t.Error("error writing node to disk", err)
		// }

		// updatedN, err := Get(n.Record)

		// if err != nil {
		// 	t.Error("error getting node from disk", err)
		// }

		// if updatedN.Data == n.Data {
		// 	t.Error("error updating forest data")
		// }

func TestNewTree (t *testing.T) {
	err := initForSynchronousTests()

	if err != nil {
		t.Error("error initializing db: ", err)
	}

	err = NewForest(convertNumToUpdater(1), convertNumToUpdater(2))

	if err != nil {
		t.Error("error putting forest in funnel: ", err)
	}

	err = clearFunnel()

	if err != nil {
		t.Error("error saving forest to disk", err)
	}

	forests, err := GetForests()
	forest := forests[0]
	if err != nil {
		t.Error("error getting forest from db")
	}

	err = NewTree(forest, convertNumToUpdater(3), convertNumToUpdater(4))

	if err != nil {
		t.Error("error putting tree in funnel: ", err)
	}

	err = NewTree(forest, convertNumToUpdater(5), convertNumToUpdater(6))

	if err != nil {
		t.Error("error putting tree in funnel: ", err)
	}

	err = clearFunnel()

	if err != nil {
		t.Error("error saving trees to disk", err)
	}

	firstTrees, err := GetChildren(forest.Record)
	//the get(node.Record) syntax guarantees the most up to date version of a
	//given node *that is on the database* since we're using the clearFunnel
	//function outside of the funnel, we need to use this to get at and check
	//the actual saved data.  clearFunnel is deliberately not part of the api
	//because opening a transaction when one is open causes an error instead of
	//blocking, so sequential transactions can't really be done concurrently.
	//  This method of getting updated records immediately to and then from the
	//db is only meant for testing and if initDb has been called then it can 
	//cause errors that are not easily caught both in your code and in the
	//funnel.
	//getForests and getForestsMeta always use this syntax (without
	//clearFunnel) which is why it didn't show up in the last test function.

	if err != nil {
		t.Error("error loading trees", err)
	}

	var parentTree Node

	for _, v := range firstTrees {
		if v.Data == convertNumToUpdater(6) {
			parentTree = v
		}
	}

	err = NewTree(parentTree, convertNumToUpdater(7), convertNumToUpdater(8))

	if err != nil {
		t.Error("error putting tree in funnel: ", err)
	}

	clearFunnel()

	trees, err := GetChildren(forest.Record)

	if err != nil {
		t.Error("error getting updated first level trees", err)
	}

	var updatedParentTree Node

	var notParentTree Node

	for _, v := range trees {
		if v.Data == convertNumToUpdater(6) {
			updatedParentTree = v
		} else {
			notParentTree = v
		}
	}

	tempTreeList, err := GetChildren(updatedParentTree.Record)

	if err != nil {
		t.Error("error getting second level tree", err)
	}

	childTree := tempTreeList[0]

	if updatedParentTree.Data != parentTree.Data {
		t.Error("adding a child changed parent's data!")
	}

	if !childTree.Parent.Loc.equals(updatedParentTree.Loc) {
		t.Error("child does not have the right data for it's parent")
	}

	if !updatedParentTree.Children[childTree.KeyString()].Loc.equals(childTree.Loc) {
		t.Error("parent does not have the right data for it's child")
	}

	if notParentTree.Data != convertNumToUpdater(4) {
		t.Error("the not parent tree has the wrong data!")
	}
}


