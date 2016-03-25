package keyChain

import (
	"testing"
	"bytes"
	"fmt"
)

func TestMakeId (t *testing.T) {
	h1_1, err := makeId(1)

	if err != nil {
		t.Error("error makeing Id: ", err)
	}

	h1_2, err := makeId(1)

	if err != nil {
		t.Error("error makeing Id: ", err)
	}

	h2, err := makeId(2)

	if err != nil {
		t.Error("error makeing Id: ", err)
	}

	if h1_1.Height != 1 {
		t.Error("height should have been 1 but was: ", h1_1.Height)
	}

	if bytes.Equal(h1_1.Identifier, h1_2.Identifier) {
		t.Error("two keys with the same height also had the same identifier!")
	}

	if h2.Height != 2 {
		t.Error("height should have been 2 but was: ", h2.Height)
	}

	if bytes.Equal(h1_1.Identifier, h2.Identifier) {
		t.Error("two keys with different heights had the same identifier!")
	}
}

func TestHeightToByteSlice (t *testing.T) {
	i, err := makeId(257)
	
	if err != nil {
		t.Error("error makeing Id: ", err)
	}

	var byteSlice []byte

	byteSlice, err = i.heightToByteSlice()

	if err != nil {
		t.Error("error converting height to byte slice: ", err)
	}

	if len(byteSlice) != 8 {
		t.Error("byte slice of 64 bit uint should have been 8 bits long, but was: ", len(byteSlice))
	}

	if byteSlice[len(byteSlice) - 1] != 1 || byteSlice[len(byteSlice) - 2] != 1 {
		t.Error("byteSlice is holding the wrong number: ", byteSlice)
	}
}


func TestKey (t *testing.T) {
	i1, err := makeId(257)

	if err != nil {
		t.Error("error making Id: ", err)
	}

	i2, err := makeId(257)

	if err != nil {
		t.Error("error making Id: ", err)
	}

	h1, err := i1.heightToByteSlice()

	if err != nil {
		t.Error("error converting height to byte slice: ", err)
	}

	k1, err := i1.Key()

	if err != nil {
		t.Error("error converting id to key: ", err)
	}

	k2, err := i2.Key()

	if err != nil {
		t.Error("error converting id to key: ", err)
	}

	if !bytes.Equal(h1, k1[:8]) {
		t.Error("height was not stored correctly: ", k1[:8])
		fmt.Println("should have been: ", h1)
	}

	if bytes.Equal(k1[8:], k2[8:]) {
		t.Error("different ids' were the same!")
	}

}