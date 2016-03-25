package keyChain

import (
	"fmt"
	"github.com/nu7hatch/gouuid"
	"encoding/binary"
	"bytes"
)

//A namespace Component.
type Id struct {
	Identifier []byte
	Height uint64
}

func (i Id) heightToByteSlice () ([]byte, error) {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.BigEndian, i.Height)

	if err != nil {
		fmt.Println("error writing uint to buffer: ", err)
		return nil, err
	}

	return buf.Bytes(), nil
}

func (i Id) Key () ([]byte, error) {
	hSlice, err := i.heightToByteSlice()

	if err != nil {
		fmt.Println("error converting height to byteSlice: ", err)
		return nil, err
	}

	return append(hSlice, i.Identifier...), err
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
