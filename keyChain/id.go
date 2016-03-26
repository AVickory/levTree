package keyChain

import (
	"fmt"
	"github.com/nu7hatch/gouuid"
	"bytes"
)

//A namespace Component.
type Id struct {
	Identifier []byte
	Height uint64
}

var bitwiseConverters []uint64 = []uint64{
	0Xff00000000000000,
	0Xff000000000000,
	0Xff0000000000,
	0Xff00000000,
	0Xff000000,
	0Xff0000,
	0Xff00,
	0Xff,
}

func (i Id) heightToByteSlice () []byte {
	if(i.Height == 0) {
		return nil
	}
	byteSlice := make([]byte, 8)

	for ind := range byteSlice {
		byteSlice[ind] = byte( (i.Height & bitwiseConverters[ind]) >> uint( (7 - ind) * 8 ) )
	}
	
	return byteSlice
}

func (i Id) Key () []byte {
	return append(i.heightToByteSlice(), i.Identifier...)
}

func makeId (h uint64) (Id, error) {
	identifier, err := uuid.NewV4()

	if err != nil {
		fmt.Println("UUID GENERATOR ERROR: ", err)
		return Id{}, err
	}

	i := Id{
		Identifier: identifier[:],
		Height: h,
	}

	return i, nil
}

func (i Id) makeChildId() (Id, error) {
	return makeId(i.Height + 1)
}

func (i1 Id) Equal(i2 Id) bool {
	return i1.Height == i2.Height && bytes.Equal(i1.Identifier, i2.Identifier)
}
