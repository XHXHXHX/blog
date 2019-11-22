package result

import (
	"blog/library/log"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

func MakeResult(rows *sql.Rows) ([]map[string] string, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	l := len(columns)
	result := make([]map[string] string, 0)
	values := make([]interface{}, l)
	valueCollect := make([]interface{}, l)

	for i := 0; i < l; i++ {
		valueCollect[i] = &values[i]
	}

	for rows.Next() {
		tmp := make(map[string] string )
		_ = rows.Scan(valueCollect...)
		for i, name := range columns {
			val := values[i]
			switch v := val.(type) {
				case []byte:
					tmp[name] = string(v)
				case string:
					tmp[name] = v
				case int:
					tmp[name] = strconv.Itoa(v)
				case int64:
					tmp[name] = strconv.FormatInt(int64(v), 10)
				case nil:

				default:
					panic(fmt.Sprintf("column type error %v", v))
			}
		}

		result = append(result, tmp)
	}

	return result, nil
}



func ModelResult(model interface{}, resArr []map[string]string) error {
	typ := reflect.TypeOf(model)
	typElem := typ.Elem()
	switch typElem.Kind() {
	case reflect.Struct:
		AutoReflectField(model, resArr[0])
	case reflect.Slice:
		val := reflect.ValueOf(model)
		ve := val.Elem()
		typElemElem := typElem.Elem()
		for i := 0; i < len(resArr); i++ {
			slice := reflect.New(typElemElem)
			AutoReflectField(slice.Interface(), resArr[i])
			ve.Set(reflect.Append(ve, slice.Elem()))
		}
	default:
		log.New().Error("ModelResult param type error")
		return errors.New("param type error")
	}

	return nil
}

func AutoReflectField(model interface{}, resMap map[string]string) {
	model_type := reflect.TypeOf(model).Elem()
	model_val := reflect.ValueOf(model)
	model_val_elem := model_val.Elem()
	for i := 0; i < model_type.NumField(); i++ {
		field := model_type.Field(i)
		tag := field.Tag.Get("from")
		if val, ok := resMap[tag]; ok {
			model_val_elem.Field(i).SetString(val)
		}
	}
}
