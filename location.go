package levTree

/*
The Location module provides a bucketing system for namespacing keys.  id
generation defaults to guuid V4, which is sufficient for my usecase, but other
options may be added later.
*/

import (
	"fmt"
	"github.com/nu7hatch/gouuid"
)

//A namespace Component.
type id *[]byte

//A location keeps track of any number of nested buckets and an id which can be
//translated into a byte slice key.
type location struct {
	Buckets []id
	Id      id
}

//noNameSpace is a blank location to be used as the zeroth tier bucket.
//This module will break if anything is put inside noNameSpace.
var noNameSpace location = location{}

//Converts the location into a single byte slice
func (l location) Key() []byte {
	if l.Id != nil {
		var key []byte
		for _, v := range l.Buckets {
			key = append(key, *v...)
		}
		return append(key, *l.Id...)
	} else {
		return []byte{}
	}
}

//produces the key as a string.  This is primarily so that locations can be
//converted to the keys of maps.
func (l location) KeyString() string {
	return string(l.Key())
}

//Checks if both locations produce the same stringified key.  This is an
//imperfect method of checking equality, but will do for the initial
//implementation.
func (l1 location) equals(l2 location) bool {
	return l1.KeyString() == l2.KeyString()
}

//creates a location object for the bucket containing this location
func (l location) getBucketLocation() location {
	bucketIndex := len(l.Buckets) - 1
	if bucketIndex >= 0 {
		return location{
			Buckets: l.Buckets[:bucketIndex],
			Id:      l.Buckets[bucketIndex],
		}
	} else {
		return noNameSpace
	}
}

//creates a location inside this bucket with the given id
func (bucket location) getNewLocWithId(identifier id) location {
	if len(*identifier) == 0 {
		return bucket
	}
	numBuckets := len(bucket.Buckets)
	if !bucket.equals(noNameSpace) {
		newBucket := make([]id, numBuckets, numBuckets+1)
		copy(newBucket, bucket.Buckets)
		newBucket = append(newBucket, bucket.Id)
		return location{
			Buckets: newBucket,
			Id:      identifier,
		}
	} else {
		return location{
			Buckets: noNameSpace.Buckets,
			Id:      identifier,
		}
	}
}

//creates a location inside this bucket with an auto generated id
func (bucket location) getNewLoc() (location, error) {
	identifier, err := uuid.NewV4()

	if err != nil {
		fmt.Println("UUID GENERATOR ERROR: ", err)
		return bucket, err
	}
	Id := identifier[:]
	newLocation := bucket.getNewLocWithId(&Id)

	return newLocation, nil
}
