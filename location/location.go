package keyChain

/*
The location module provides a bucketing system for namespacing keys.  Id
generation defaults to guuid V4, which is sufficient for my usecase, but other
options may be added later.
*/

import (
	"fmt"
	"github.com/nu7hatch/gouuid"
	"encoding/binary"
	"bytes"
)

type location []Id

func (k location) getId () Id {
	return k[len(k) - 1]
}

func (k location) Key() []byte {
	key := make([]byte, 0, len(k)*8)

	for _, v := range l.Bucket {
		key = append(key, v.key()...)
	}

	return key
}

func (bucket location) makeBranchLocation () (location, err) {
	length := len(bucket) + 2

	k := make([]Id, 0, length)
	copy(k, bucket)

	bucketId := bucket.getId()

	k[length - 2] = bucketId
	k[length - 1], err = makeId(bucketId.Height + 1)
	
	if err != nil {
		fmt.Println("Error making Id", err)
		return k, err
	}

	return k, nil
}

func (bucket location) makeTreeLocation () (location, error) {
	bucketId := bucket.getId()

	id, err := makeId(bucketId.Height + 1)

	if err != nil {
		fmt.Println("Error making Id", err)
		return k, err
	}

	bucket = append(bucket, id)

	return bucket
} 
