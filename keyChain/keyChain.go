package keyChain

import (
	"fmt"
)
//A keyChain keeps track of any number of nested buckets and an Id which can be
//translated into a byte slice key.
type KeyChain struct {
	Self Loc
	Parent Loc
	ChildBucket Loc //note that using getId on this location doesn't really make sense
}

func MakeBranchKeyChain (parent KeyChain) (KeyChain, error) {
	self, err := makeBranchLoc(parent.ChildBucket, parent.Self)

	if err != nil {
		fmt.Println("Error making branch key:", err)
		return KeyChain{}, err
	}

	b := KeyChain{
		Self: self,
		Parent: parent.Self,
		ChildBucket: parent.ChildBucket,
	}

	return b, nil
}

func MakeTreeKeyChain (parent KeyChain) (KeyChain, error) {
	self, err := makeTreeLoc(parent.ChildBucket)

	if err != nil {
		fmt.Println("Error making branch key:", err)
		return KeyChain{}, err
	}

	t := KeyChain{
		Self: self,
		Parent: parent.Self,
		ChildBucket: self,
	}

	return t, nil
}

func (k KeyChain) IsTree () bool {
	return k.ChildBucket.Equal(k.Self)
}

//Converts the KeyChain into a single byte slice
func (k KeyChain) Key() []byte {
	return k.Self.Key()
}

//produces the key as a string.  This is primarily so that locations can be
//converted to the keys of maps.
func (k KeyChain) KeyString() string {
	return k.Self.KeyString()
}

func (k1 KeyChain) Equal(k2 KeyChain) bool {
	return k1.Self.Equal(k2.Self)
}

func (k KeyChain) Height() uint64 {
	return k.Self.Height()
}
