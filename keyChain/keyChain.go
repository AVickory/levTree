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
	IsTree bool
	NameSpace Loc //the sequence of ids to which this keychain's id must be added to get it's location
	GrandParentId Id
	ParentId Id
	Id
}

/*
rules:
	1.  Trees cannot be the descendant's of branches.
	2.  the prefix of a tree's descendants' namespace is the sequence of that tree's ancestors' ids plus that tree's id
	3.  the prefix of a branch's descendant's namespace is inherited directly from rule 2 with no modification 
	4.  The key of a tree or branch is computed by appending the keychain's immediate parent's id and it's id to the namespace generated by rule 2 and 3

results:
	1.  If the immediate parent's id is the same as the last element of the namespace then the parent is a tree
	2.  If the parent is a tree then it's location will be the child's namespace with the second to last elment of the child's namespace repeated twice
	3.  If the parent is a branch then it's location will be the child's grand parent and parent ids appended to the child's namespace
	4.  The prefix of a keychain's siblings is the keychain's namespace plus the parent's id
	5.  The prefix of a tree's descendants will be the tree's namespace plus it's own id
	6.  The prefix of a tree's immediate children will be the tree's namespace plus it's own id repeated twice
	7.  The prefix of a branch's immediate children will be the tree's namespace plus it's own id

(rules for root cases to be added soon)

I'll try to put a more full explaination of why this all works out and why it's important at a later date.
For now, I'll need you to accept that it will give us optimized look up performance for the immediate children of
All nodes and all descendants of trees.
Oh and that the results do in fact follow from the rules I set out.
That's important.
*/

func (k KeyChain) GetLoc () Loc {
	if k.Id.Equal(rootId) {
		return rootLoc
	}
	return k.NameSpace.copyAndAppend(k.ParentId, k.Id) //rule 4
}

func (k KeyChain) GetParentLoc() Loc {
	if k.Equal(Root) || k.ParentId.Equal(rootId) {
		return rootLoc
	}
	if k.ParentIsTree() {
		l := k.NameSpace.copyAndAppend(k.ParentId) //parent id is now repeated twice
		l[len(l) - 2] = k.GrandParentId //first parentId is replaced with Grandparent's Id, so that GP is repeated twice
		return l //result 2
	} else {
		return k.NameSpace.copyAndAppend(k.GrandParentId, k.ParentId) //result 3
	}

}

func (k KeyChain) GetSiblingBucket() Loc {
	return k.NameSpace.copyAndAppend(k.ParentId) //result 4
}

func (k KeyChain) GetDescendantBucket() Loc {
		return k.NameSpace.copyAndAppend(k.Id) //result 5
}

func (k KeyChain) GetChildBucket() Loc {
	if k.IsTree {
		return k.NameSpace.copyAndAppend(k.Id, k.Id) //result 6
	}
	return k.NameSpace.copyAndAppend(k.Id) //result 7
}

func (k KeyChain) childNameSpace() Loc {
	if k.IsTree {
		return k.NameSpace.copyAndAppend(k.Id) //rule 2
	} else {
		return k.NameSpace.copyAndAppend() //rule 3
	}
}

func (parent KeyChain) makeChild() (KeyChain, error) {
	childId, err := parent.Id.makeChildId()

	if err != nil {
		fmt.Println("error making child id", err)
		return KeyChain{}, err
	}

	return KeyChain{
		NameSpace: parent.childNameSpace(),
		ParentId: parent.Id,
		GrandParentId: parent.ParentId,
		Id: childId,
	}, nil	
}

func (parent KeyChain) MakeChildBranch() (KeyChain, error) {
	child, err := parent.makeChild()
	child.IsTree = false
	return child, err
}

func (parent KeyChain) MakeChildTree () (KeyChain, error) {
	//should return error if parent is a branch, but this technically shouldn't break it.
	//I just haven't explored the ramifications very thoroughly.
	child, err := parent.makeChild()
	child.IsTree = true
	return child, err
}

func (k KeyChain) MakeSibling () (KeyChain, error) {
	var err error
	k.Id, err = k.ParentId.makeChildId()
	return k, err
}

func (k KeyChain) ParentIsTree () bool {
	if k.Equal(Root) {
		return true
	}
	return k.NameSpace[len(k.NameSpace) - 1].Equal(k.ParentId)
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

func (k KeyChain) SiblingKeyPrefix() []byte {
	return k.GetSiblingBucket().Key()
}

//produces the key as a string.  This is primarily so that locations can be
//converted to the keys of maps.
func (k KeyChain) KeyString() string {
	return k.GetLoc().KeyString()
}

func (k1 KeyChain) Equal(k2 KeyChain) bool {
	return k1.GetLoc().Equal(k2.GetLoc())
}

