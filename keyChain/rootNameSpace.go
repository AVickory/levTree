package keyChain

var Root KeyChain

var rootLoc Loc

var rootId Id

func init () {
	rootId = Id{
		Identifier: []byte{},
		Height: 0,
	}

	rootLoc = Loc{
		rootId,
	}

	Root = KeyChain{
		NameSpace: Loc{},
		Id: rootId,
		IsTree: true,
	}
}