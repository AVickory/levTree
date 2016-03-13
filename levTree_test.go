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

	for _, n := range forests {
		if n.Parent.Data != forestsMeta[n.KeyString()].Data {
			t.Error("Metadata was not stored differently in root and forest")
		}
		expectedData, _ := n.Parent.Data.(mockUpdateable)
		expectedData += 1
		if n.Data != convertNumToUpdater(int(expectedData)) {
			t.Error("Wrong data stored in record with parentMeta: ", n.Parent.Data)
		}
	}

}