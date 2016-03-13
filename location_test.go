package levTree

import (
	"bytes"
	"testing"
)

func testLoc(idNums []byte) location {
	idBytes := make([][]byte, len(idNums), len(idNums))
	ids := make([]id, len(idNums), len(idNums))
	for i, v := range idNums {
		idBytes[i] = []byte{v}
		ids[i] = &idBytes[i]
	}
	return location{
		Buckets: ids[:len(idNums)-1],
		Id:      ids[len(idNums)-1],
	}
}

func TestLocationKey(t *testing.T) {
	loc1 := testLoc([]byte{1, 2, 3, 4, 5})
	loc2 := testLoc([]byte{5, 4, 3, 2, 1})
	correctKey := []byte{1, 2, 3, 4, 5}

	if !bytes.Equal(loc1.Key(), correctKey) {
		t.Error("KEY WRONG")
	}
	if !(loc1.KeyString() == string(correctKey)) {
		t.Error("KEYSTRING WRONG")
	}
	if !(loc1.equals(loc1)) {
		t.Error("EQUALS EXPECTED TO BE TRUE")
	}
	if loc1.equals(loc2) {
		t.Error("EQUALS EXPECTED TO BE FALSE")
	}
}

func TestGetBucketLocation(t *testing.T) {
	loc := testLoc([]byte{1, 2, 3, 4, 5})
	correctBucket := testLoc([]byte{1, 2, 3, 4})
	bucket := loc.getBucketLocation()
	if !(correctBucket.equals(bucket)) {
		t.Error("GET BUCKET ERROR")
	}
}

func TestGetNewLocWithId(t *testing.T) {
	correctLoc := testLoc([]byte{1, 2, 3, 4, 5})
	bucket := correctLoc.getBucketLocation()
	loc := bucket.getNewLocWithId(correctLoc.Id)
	if !(correctLoc.equals(loc)) {
		t.Error("BASIC GETLOCWITHID ERROR")
	}

	Id := &[]byte{}
	correctLoc = noNameSpace
	loc = noNameSpace.getNewLocWithId(Id)
	if !(correctLoc.equals(loc)) {
		t.Error("EMPTYID GETLOCWITHID ERROR")
	}

	Id1 := &[]byte{1}
	correctLoc = location{Id: Id1}
	loc = noNameSpace.getNewLocWithId(Id1)
	if !(correctLoc.equals(loc)) {
		t.Error("noNameSpace GETLOCWITHID ERROR")
	}

	Id2 := &[]byte{2}
	correctLoc = location{
		Id:      Id2,
		Buckets: []id{Id1},
	}
	loc = loc.getNewLocWithId(Id2)
	if !(correctLoc.equals(loc)) {
		t.Error("NESTED GETLOCWITHID ERROR")
	}
}

func TestGetNewLoc(t *testing.T) {
	loc := testLoc([]byte{1, 2})
	locLen := len(loc.Key())
	loc, err := loc.getNewLoc()

	if len(loc.Key()) != locLen+16 || err != nil {
		t.Error("GUUID NOT APPENDED TO NEW LOC")
	}

}
