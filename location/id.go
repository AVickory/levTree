package keyChain

import (
	"fmt"
	"github.com/nu7hatch/gouuid"
	"encoding/binary"
)

//A namespace Component.
type Id struct {
	Identifier []byte
	Height uint64
}

func (i Id) heightTobyteArr () []byte {
	byteArr := make([]byte, 8)

	_ = binary.PutUvarint(byteArr, i.Height)

	return byteArr
}

func (i Id) key () []byte {
	return append(i.heightTobyteArr(), i.Identifier...)
}

func makeId (h uint64) (Id, error) {
	identifier, err := uuid.NewV4()

	if err != nil {
		fmt.Println("UUID GENERATOR ERROR: ", err)
		return nil, err
	}

	i := Id{
		Identifier: identifier[:],
		Height: h,
	}

	return i
}

func (k location) getId () Id {
	return k[len(k) - 1]