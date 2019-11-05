package sqlBuild

import (
	"errors"
	"strconv"
)

var (
	options = [] string {
		"=", "<", ">", "<=", ">=", "<>", "!=", "<=>",
		"ike", "ike binary", "ot like", "like",
		"&", "|", "^", "<<", ">>",
		"like", "ot rlike", "egexp", "ot regexp",
		"~", "~*", "!~", "!~*", "imilar to",
		"ot similar to", "ot ilike", "~~*", "!~~*",}
)

type whereInfo struct {
	boolean bool			// false AND     true OR
	whereType string
	conditionType string
	field string
	value interface{}
	build *SqlBuild
}

func newWhere(whereType, conditionType, field string, value interface{}, boolean bool) (*whereInfo, error) {
	return &whereInfo{
		boolean: boolean,
		whereType: whereType,
		conditionType: conditionType,
		field: field,
		value: value,
	}, nil
}

func newBuild(build *SqlBuild, whereType, conditionType string, boolean bool) (*whereInfo, error) {
	return &whereInfo{
		boolean: boolean,
		whereType: whereType,
		conditionType: conditionType,
		build: build,
	}, nil
}

func (this *whereInfo) Where(args... string) (*whereInfo, error) {
	boolean := args[len(args)-1]
	args = args[:len(args) - 1]
	if len(args) == 0 {
		panic("Where param error")
	}

	var field, conditionType, value string = "", "",""
	switch len(args) {
		case 1:
			value = args[0]
		case 2:
			field = args[0]
			value = args[1]
			conditionType = "="
		case 3:
			field = args[0]
			value = args[2]
			conditionType = args[1]
		default:
			panic("Where func params num up to 3")
	}

	if !InvalidOperator(conditionType) {
		return nil, errors.New("condition option error")
	}

	return newWhere("Basic", conditionType, field, value, boolean == "0")
}

func (this *whereInfo) WhereFunc(callback func(build *SqlBuild) *SqlBuild, build *SqlBuild, boolean bool) (*whereInfo, error) {
	return newBuild(callback(build), "Nested", "", boolean)
}

func (this *whereInfo) WhereExists(conditionType string, callback func(build *SqlBuild) *SqlBuild, boolean bool) (*whereInfo, error) {
	return newBuild(callback(NewBuild()), "Exists", conditionType, boolean)
}

func (this *whereInfo) WhereArray(arrayWhere [][]string, newBuild *SqlBuild, boolean bool) (*whereInfo, error) {
	return this.WhereFunc(func(build *SqlBuild) *SqlBuild {
		for _, value := range arrayWhere {
			if boolean {
				_ = build.OrWhereString(value...)
			} else {
				_ = build.WhereString(value...)
			}
		}
		return build
	}, newBuild, boolean)
}

func (this *whereInfo) WhereMap(mapWhere map[string] interface{}, newBuild *SqlBuild, boolean bool) (*whereInfo, error) {
	return this.WhereFunc(func(build *SqlBuild) *SqlBuild {
		for key, value := range mapWhere {
			_ = build.Where(key, value)
		}
		return build
	}, newBuild, boolean)
}

func (this *whereInfo) whereInFactory(isNot bool, field string, listValue [] interface{}, boolean bool) (*whereInfo, error) {
	var myList []string
	for _, item := range listValue {
		myList = append(myList, "'" + TransferString(item) + "'")
	}

	conditionType := "IN"
	if isNot {
		conditionType = "NOT IN"
	}
	return newWhere("In", conditionType, field, myList, boolean)
}


func (this *whereInfo) WhereIn(field string, listValue [] interface{}, boolean bool) (*whereInfo, error) {
	return this.whereInFactory(false ,field, listValue, boolean)
}

func (this *whereInfo) WhereNotIn(field string, listValue [] interface{}, boolean bool) (*whereInfo, error) {
	return this.whereInFactory(true ,field, listValue, boolean)
}

func (this *whereInfo) whereBetweenFatory(isNot bool, field string, interval [] interface{}, boolean bool) (*whereInfo, error) {
	var myList []string
	for i := 0; i < 2; i++ {
		switch item := interval[i].(type) {
			case string:
				myList = append(myList, "'" + item + "'")
			case int:
				myList = append(myList, strconv.Itoa(item))
			default:
				panic("WhereBetween param error")
		}
	}

	conditionType := "Between"
	if isNot {
		conditionType = "NOT Between"
	}
	return newWhere("Between", conditionType, field, myList, boolean)
}

func (this *whereInfo) WhereBetween(field string, interval [] interface{}, boolean bool) (*whereInfo, error) {
	if len(field) == 0 || len(interval) != 2 {
		panic("WhereBetween param error")
	}
	return this.whereBetweenFatory(false, field, interval, boolean)
}

func (this *whereInfo) WhereNotBetween(field string, interval [] interface{}, boolean bool) (*whereInfo, error) {
	return this.whereBetweenFatory(true, field, interval, boolean)
}

func (this *whereInfo) WhereNull(field string, boolean bool) (*whereInfo, error) {
	return newWhere("Null","IS", field, "NULL", boolean)
}

func (this *whereInfo) WhereNotNull(field string, boolean bool) (*whereInfo, error) {
	return newWhere("Null","GIS NOT", field, "NULL", boolean)
}

func (this *whereInfo) WhereDate(field, date string, boolean bool) (*whereInfo, error) {
	char := GetDateStringJoiner(date)
	format := SetMysqlDateFormatByChar(char)
	func_field := "DATE_FORMAT(" + field + ", '" + format + "')"
	return newWhere("Func", "=", func_field, AddSingleSymbol(date), boolean)
}

func (this *whereInfo) WhereMonth(field, month string, boolean bool) (*whereInfo, error) {
	func_field := "'DATE_FORMAT(" + field + ", '%m')'"
	return newWhere("Func","=", func_field, AddSingleSymbol(month), boolean)
}

func (this *whereInfo) WhereDay(field, day string, boolean bool) (*whereInfo, error) {
	func_field := "DATE_FORMAT(" + field + ", '%d')"
	return newWhere("Func","=", func_field, AddSingleSymbol(day), boolean)
}

func (this *whereInfo) whereYear(field, year string, boolean bool) (*whereInfo, error) {
	func_field := "DATE_FORMAT(" + field + ", '%Y')"
	return newWhere("Func","=", func_field, AddSingleSymbol(year), boolean)
}

func (this *whereInfo) WhereTime(field, condition, timestamp string, boolean bool) (*whereInfo, error) {
	func_field := "DATE_FORMAT(" + field + ", '%H:%i:%s')"
	return newWhere("Func", condition, func_field, AddSingleSymbol(timestamp), boolean)
}

func (this *whereInfo) WhereRaw(sql string, boolean bool) (*whereInfo, error) {
	return newWhere("Raw", "", "", sql, boolean)
}
