package keyChain

/*
The location module provides a bucketing system for namespacing keys.  Id
generation defaults to guuid V4, which is sufficient for my usecase, but other
options may be added later.
*/

import (
	"fmt"
)

type loc []Id

func makeBranchLoc (bucket loc, parent loc) (loc, error) {
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

func makeTreeLoc (bucket loc) (loc, error) {
	bucketId := bucket.GetId()

	if bucketId.Identifier == nil {
		bucketId.Height = 0
	}

	id, err := bucketId.makeChildId()

	if err != nil {
		fmt.Println("Error making Id", err)
		return loc{}, err
	}

	bucket = append(bucket, id)

	return bucket, nil
} 

func (loc loc) Key() []byte {
	key := make([]byte, 0, len(loc)*8)

	for _, id := range loc {
		key = append(key, id.Key()...)
	}

	return key
}

func (loc loc) KeyString() string {
	return string(loc.Key())
}

func (loc loc) GetId () Id {
	if(len(loc) != 0) {
		return loc[len(loc) - 1]
	} else {
		return Id{}
	}
}

func (loc1 loc) Equal (loc2 loc) bool {
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
