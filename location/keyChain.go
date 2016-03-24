package keyChain

import (
	"fmt"
)
//A keyChain keeps track of any number of nested buckets and an Id which can be
//translated into a byte slice key.
type keyChain struct {
	Self location
	Parent location
	ChildBucket location //note that using getId on this location doesn't really make sense
}

func MakeBranchKeyChain (parentLoc keyChain) keyChain, error {
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

func MakeTreeKeyChain (parent keyChain) keyChain {
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

func (l keyChain) IsTree () bool {
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
