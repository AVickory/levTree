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

//noNameSpace is a blank keyChain to be used as the zeroth tier bucket.
var noNameSpace keyChain = keyChain{
	Self: location{}
	ChildBucket: location{}
}

func MakeBranchKeyChain (parentLoc keyChain) (keyChain, error) {
	self, err := parentLoc.ChildBucket.makeBranchLocation()

	if err != nil {
		fmt.Println("Error making branch key:", err)
		return keyChain{}, err
	}

	b := keyChain{
		Self: self,
		Parent: parentLoc.Self,
		ChildBucket: parentLoc.ChildBucket,
	}

	return b, nil
}

func MakeTreeKeyChain (parent keyChain) (keyChain, error) {
	self, err := parent.ChildBucket.makeTreeLocation()

	if err != nil {
		fmt.Println("Error making branch key:", err)
		return keyChain{}, err
	}

	t := keyChain{
		Self: self,
		Parent: parent.Self,
		ChildBucket: append(parent.ChildBucket, self.getId()),
	}

	return t, nil
}

func (k keyChain) IsTree () bool {
	return k.ChildBucket.Equals(k.Self)
}

//Converts the keyChain into a single byte slice
func (k keyChain) Key() ([]byte, error) {
	return k.Self.Key()
}

//produces the key as a string.  This is primarily so that locations can be
//converted to the keys of maps.
func (k keyChain) KeyString() (string, error) {
	return k.Self.KeyString()
}
