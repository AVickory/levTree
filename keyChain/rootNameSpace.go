package keyChain

var root keyChain

var rootLoc loc

var rootId Id

func init () {
	rootId = Id{
		Identifier: []byte{},
		Height: 0,
	}

	rootLoc = loc{
		rootId,
	}

	root = keyChain{
		Self: rootLoc,
		Parent: rootLoc,
		ChildBucket: rootLoc,
	}
}