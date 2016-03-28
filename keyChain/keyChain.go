package keyChain

import (
	"fmt"
)
//A keyChain keeps track of any number of nested buckets and an Id which can be
//translated into a byte slice key.
// type KeyChain struct {
// 	Self Loc
// 	Parent Loc
// 	ChildBucket Loc
// }

type KeyChain struct {
	NameSpace Loc
	Ancestors Loc //this contains up to two entries.  Not actually a location, don't use .Key on it
	Id
}

func (child KeyChain) GetParentLoc () Loc {
	if child.IsTree() {
		return child.NameSpace[:len(child.NameSpace) - 1]		
	} else if child.ParentIsTree() {
		//parent is a tree and child is a branch
		return child.NameSpace
	} else {
		//parent and child are branches
		return child.NameSpace.copyAndAppend(child.Ancestors...)
	}
	//child is a tree
}

func (parent KeyChain) GetChildBucket() Loc {
	if !parent.IsTree() {
		return parent.NameSpace.copyAndAppend(parent.Id)
	} else {
		return parent.NameSpace
	}
}

func (k KeyChain) GetLoc() Loc {
	var self Loc
	if !k.IsTree() {
		self = k.NameSpace.copyAndAppend(k.Ancestors.GetId(), k.Id)
	} else {
		self = k.NameSpace
	}
	return self
}

func (parent KeyChain) MakeChildBranch() (KeyChain, error) {
	if parent.IsTree() || parent.ParentIsTree() {
		//fewer than two elements in ancestors, so we can just append the parent.
		parent.Ancestors = parent.Ancestors.copyAndAppend(parent.Id)
	} else {
		newAncestors := make([]Id, 2, 2)
		newAncestors[0] = parent.Ancestors[1]
		newAncestors[1] = parent.Id
		parent.Ancestors = newAncestors
	}

	var err error
	parent.Id, err = parent.Id.makeChildId()

	if err != nil {
		fmt.Println("error making child Id: ", err)
		return parent, err
	}

	return parent, nil
}

func (parent KeyChain) MakeChildTree () (KeyChain, error) {
	parent.Ancestors = make([]Id, 0)
	var err error
	parent.Id, err = parent.Id.makeChildId()
	parent.NameSpace = parent.NameSpace.copyAndAppend(parent.Id)
	if err != nil {
		fmt.Println("error making child Id: ", err)
		return parent, err
	}

	return parent, nil
}

func (k KeyChain) MakeSibling () (KeyChain, error) {
	var err error
	k.Id, err = k.Id.makeSiblingId()
	if err != nil {
		fmt.Println("error making sibling id", err)
		return k, err
	}
	return k, nil
}


func (k KeyChain) IsTree () bool {
	return k.Id.Equal(k.NameSpace.GetId()) && len(k.Ancestors) == 0
}

func (branch KeyChain) ParentIsTree () bool {
	return len(branch.Ancestors) == 1
}

//Converts the KeyChain into a single byte slice
func (k KeyChain) Key() []byte {
	return k.GetLoc().Key()
}

func (k KeyChain) ParentKey() []byte {
	return k.GetParentLoc().Key()
}

func (k KeyChain) ChildKeyPrefix() []byte {
	return k.GetChildBucket().Key()
}

//produces the key as a string.  This is primarily so that locations can be
//converted to the keys of maps.
func (k KeyChain) KeyString() string {
	return k.GetLoc().KeyString()
}

func (k1 KeyChain) Equal(k2 KeyChain) bool {
	return k1.GetLoc().Equal(k2.GetLoc())
}

