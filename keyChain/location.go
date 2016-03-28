package keyChain

/*
The location module provides a bucketing system for namespacing keys.  Id
generation defaults to guuid V4, which is sufficient for my usecase, but other
options may be added later.
*/

import (
	// "fmt"
)

type Loc []Id

func (l Loc) copyAndAppend (Ids ...Id) Loc {
	newL := make([]Id, len(l), len(l) + len(Ids))
	copy(newL, l)
	newL = append(newL, Ids...)
	return newL
}

func (Loc Loc) Key() []byte {
	key := make([]byte, 0, len(Loc)*8)

	for _, id := range Loc {
		key = append(key, id.Key()...)
	}

	return key
}

func (loc Loc) KeyString() string {
	return string(loc.Key())
}

func (loc Loc) GetId () Id {
	if(len(loc) != 0) {
		return loc[len(loc) - 1]
	} else {
		return rootId
	}
}

func (loc1 Loc) Equal (loc2 Loc) bool {
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
