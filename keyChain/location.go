package keyChain

/*
The location module provides a bucketing system for namespacing keys.  Id
generation defaults to guuid V4, which is sufficient for my usecase, but other
options may be added later.
*/

import (
	"fmt"
	"bytes"
)

type location []Id

func (l location) getId () Id {
	return l[len(l) - 1]
}

func (l location) Key() ([]byte, error) {
	key := make([]byte, 0, len(l)*8)

	for _, id := range l {
		idKey, err := id.Key()

		if err != nil {
			fmt.Println("error converting id, ", id, " to key")
			return nil, err
		}

		key = append(key, idKey...)
	}

	return key, nil
}

func (l location) KeyString() (string, error) {
	k, err := l.Key()
	return string(k), err
}

func (bucket location) makeBranchLocation () (location, error) {
	length := len(bucket) + 2

	l := make([]Id, 0, length)
	copy(l, bucket)

	bucketId := bucket.getId()

	l[length - 2] = bucketId
	var err error
	l[length - 1], err = makeId(bucketId.Height + 1)
	
	if err != nil {
		fmt.Println("Error making Id", err)
		return l, err
	}

	return l, nil
}

func (bucket location) makeTreeLocation () (location, error) {
	bucketId := bucket.getId()

	id, err := makeId(bucketId.Height + 1)

	if err != nil {
		fmt.Println("Error making Id", err)
		return location{}, err
	}

	bucket = append(bucket, id)

	return bucket, nil
} 

func (l1 location) Equals (l2 location) bool {
	k1, err := l1.Key()

	if err != nil {
		fmt.Println("Error converting location 1 to key to check equality", err)
		return false
	}

	k2, err := l2.Key()

	if err != nil {
		fmt.Println("Error converting location 1 to key to check equality", err)
		return false
	}

	return bytes.Equal(k1, k2)
}
