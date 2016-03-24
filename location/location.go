package levTree

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

//A namespace Component.
type Id struct {
	Identifier []byte
	Height uint64
}

func (i Id) heightTobyteArr () []byte {
	byteArr := make([]byte, 8)

	_ = binary.PutUvarint(byteArr, i.Height)

	return byteArr
}

func (i Id) key () []byte {
	return append(i.heightTobyteArr(), i.Identifier...)
}

func makeId (h uint64) (Id, error) {
	identifier, err := uuid.NewV4()

	if err != nil {
		fmt.Println("UUID GENERATOR ERROR: ", err)
		return nil, err
	}

	i := Id{
		Identifier: identifier[:],
		Height: h,
	}

	return i
}

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


//A keyChain keeps track of any number of nested buckets and an Id which can be
//translated into a byte slice key.
type keyChain struct {
	Self location
	Parent location
	ChildBucket location //note that using getId on this location doesn't really make sense
}

func makeBranchKeyChain (parentLoc keyChain) keyChain, error {
	self, err := parent.ChildBucket.makeBranchLocation()

	if err != nil {
		fmt.Println("Error making branch key:", err)
		return nil, err
	}

	b := keyChain{
		Self: self,
		Parent: parent.Self,
		ChildBucket: parent.ChildBucket,
	}

	return b
}

func makeTreeKeyChain (parent keyChain) keyChain {
	self, err := parent.ChildBucket.makeTreeLocation()

	if err != nil {
		fmt.Println("Error making branch key:", err)
		return nil, err
	}

	t := keyChain{
		Self: self,
		Parent: parent.Self,
		ChildBucket: append(parent.ChildBucket, self),
	}
}

func (l keyChain) isTree () bool {
	return bytes.Equal(l.ChildBucket.Key(), l.Self.Key())
}

//noNameSpace is a blank keyChain to be used as the zeroth tier bucket.
//This module will break if anything is put inside noNameSpace.
var noNameSpace keyChain = keyChain{}

//Converts the keyChain into a single byte slice
func (l keyChain) Key() []byte {
	l.Self.Key()
}

//produces the key as a string.  This is primarily so that locations can be
//converted to the keys of maps.
func (l keyChain) KeyString() string {
	return string(l.Key())
}
