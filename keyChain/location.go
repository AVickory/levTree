package keyChain

/*
The location module provides a bucketing system for namespacing keys.  Id
generation defaults to guuid V4, which is sufficient for my usecase, but other
options may be added later.
*/

import (
	"fmt"
)

type Loc []Id

func makeBranchLoc (bucket Loc, parent Loc) (Loc, error) {
	length := len(bucket) + 2

	loc := make([]Id, length)
	copy(loc, bucket)

	parentId := parent.GetId()

	loc[length - 2] = parentId
	var err error
	loc[length - 1], err = parentId.makeChildId()
	
	if err != nil {
		fmt.Println("Error making Id", err)
		return loc, err
	}

	return loc, nil
}

func makeTreeLoc (bucket Loc) (Loc, error) {
	bucketId := bucket.GetId()

	if bucketId.Identifier == nil || len(bucketId.Identifier) == 0{
		bucketId.Height = 0
	}

	id, err := bucketId.makeChildId()

	if err != nil {
		fmt.Println("Error making Id", err)
		return Loc{}, err
	}

	bucket = append(bucket, id)

	return bucket, nil
} 

func (Loc Loc) Key() []byte {
	key := make([]byte, 0, len(Loc)*8)

	for _, id := range Loc {
		key = append(key, id.Key()...)
	}

	return key
}

func (loc Loc) KeyString() string {
	return string(loc.Key())
}

func (loc Loc) GetId () Id {
	if(len(loc) != 0) {
		return loc[len(loc) - 1]
	} else {
		return rootId
	}
}

func (loc1 Loc) Equal (loc2 Loc) bool {
	if len(loc1) != len(loc2) {
		return false
	}

	for ind, id := range loc1 {
		if !id.Equal(loc2[ind]) {
			return false
		}
	}

	return true
}

func (loc Loc) Height() uint64 {
	return loc.GetId().Height
}
