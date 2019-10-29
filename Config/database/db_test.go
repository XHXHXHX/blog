package database

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGetMysqlConfig(t *testing.T) {
	config, err := GetMysqlConfig()

	if err != nil {
		t.Errorf("%s", err)
	}

	fmt.Println(config)
	fmt.Println(reflect.TypeOf(config))
}
