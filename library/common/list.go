package common

import "errors"

func Pop(myList []*interface{}) (interface{}, error) {
	if len(myList) == 0 {
		return nil, errors.New("empty")
	}
	return myList[0], nil
}
