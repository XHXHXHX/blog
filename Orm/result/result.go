package result

import (
	"database/sql"
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
				default:
					panic("column type error")
			}
		}

		result = append(result, tmp)
	}

	return result, nil
}
