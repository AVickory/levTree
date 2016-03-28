package keyChain

import (
	"testing"
	"bytes"
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
	i, err := makeId(bitwiseConverters[0] | bitwiseConverters[5] | bitwiseConverters[7])
	
	if err != nil {
		t.Error("error makeing Id: ", err)
	}

	var byteSlice []byte

	byteSlice = i.heightToByteSlice()

	if len(byteSlice) != 8 {
		t.Error("byte slice of 64 bit uint should have been 8 bits long, but was: ", len(byteSlice))
	}

	expected := []byte{
		255, 0, 0, 0, 0, 255, 0, 255,
	}

	if !bytes.Equal(byteSlice, expected) {
		t.Error("byteSlice is holding the wrong number: ", byteSlice)
	}

}

func TestIdKey (t *testing.T) {
	i1, err := makeId(257)

	if err != nil {
		t.Error("error making Id: ", err)
	}

	i2, err := makeId(257)

	if err != nil {
		t.Error("error making Id: ", err)
	}

	h1 := i1.heightToByteSlice()

	k1 := i1.Key()

	k2 := i2.Key()

	if !bytes.Equal(h1, k1[:8]) {
		t.Error("height was not stored correctly in key: ", k1[:8], "\n\tshould have been: ", h1)
	}

	if bytes.Equal(k1[8:], k2[8:]) {
		t.Error("different ids' were the same!")
	}

}


func TestMakeChildId (t *testing.T) {
	i1, err := makeId(0)

	if err != nil {
		t.Error("error making first id: ", err)
	}

	i2, err := i1.makeChildId()

	if err != nil {
		t.Error("error making child id: ", err)
	}

	if i2.Height != 1 {
		t.Error("child height was not set correctly.  Height was: ", i2.Height)
	}

	if bytes.Equal(i2.Identifier, i1.Identifier) {
		t.Error("child had the same identifier as parent")
	}
}

func TestMakeSiblingId (t *testing.T) {
	i1, err := makeId(0)

	if err != nil {
		t.Error("error making first id: ", err)
	}

	i2, err := i1.makeSiblingId()

	if err != nil {
		t.Error("error making child id: ", err)
	}

	if i2.Height != 0 {
		t.Error("child height was not set correctly.  Height was: ", i2.Height)
	}

	if bytes.Equal(i2.Identifier, i1.Identifier) {
		t.Error("child had the same identifier as parent")
	}	
}

func TestRootId (t *testing.T) {
	i, err := rootId.makeChildId()

	if err != nil {
		t.Error("error making id: ", err)
	}

	if i.Height != 1 {
		t.Error("childHeight is not correct: ", i.Height)
	}

	if len(i.Key()) != 24 {
		t.Error("there were not the right number of bytes in the new id's key: ", len(i.Key()))
	}

	if bytes.Equal(i.Identifier, rootId.Identifier) {
		t.Error("Id had the same identifier as root")
	}
}