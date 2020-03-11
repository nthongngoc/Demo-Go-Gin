package helpers

import "gopkg.in/mgo.v2/bson"

func CheckItemExists(arrayType []string, item string) (int, bool) {
	for i, v := range arrayType {
		if v == item {
			return i, true
		}
	}
	return -1, false
}

func CheckObjectIdExists(arrayType []bson.ObjectId, item string) (int, bool) {
	if item == "" {
		return -1, false
	}

	for i, v := range arrayType {
		if v == bson.ObjectIdHex(item) {
			return i, true
		}
	}
	return -1, false
}
