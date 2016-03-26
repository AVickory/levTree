package keyChain

import (
	"fmt"
)
//A keyChain keeps track of any number of nested buckets and an Id which can be
//translated into a byte slice key.
type keyChain struct {
	Self loc
	Parent loc
	ChildBucket loc //note that using getId on this location doesn't really make sense
}

func MakeBranchKeyChain (parent keyChain) (keyChain, error) {
	self, err := makeBranchLoc(parent.ChildBucket, parent.Self)

	if err != nil {
		fmt.Println("Error making branch key:", err)
		return keyChain{}, err
	}

	b := keyChain{
		Self: self,
		Parent: parent.Self,
		ChildBucket: parent.ChildBucket,
	}

	return b, nil
}

func MakeTreeKeyChain (parent keyChain) (keyChain, error) {
	self, err := makeTreeLoc(parent.ChildBucket)

	if err != nil {
		fmt.Println("Error making branch key:", err)
		return keyChain{}, err
	}

	t := keyChain{
		Self: self,
		Parent: parent.Self,
		ChildBucket: self,
	}

	return t, nil
}

func (k keyChain) IsTree () bool {
	return k.ChildBucket.Equal(k.Self)
}

//Converts the keyChain into a single byte slice
func (k keyChain) Key() []byte {
	return k.Self.Key()
}

//produces the key as a string.  This is primarily so that locations can be
//converted to the keys of maps.
func (k keyChain) KeyString() string {
	return k.Self.KeyString()
}

func (k1 keyChain) Equal(k2 keyChain) bool {
	return k1.Self.Equal(k2.Self)
}
