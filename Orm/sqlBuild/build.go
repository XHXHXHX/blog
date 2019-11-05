package sqlBuild

import (
	"blog/Config/database"
	"strconv"
	"strings"
	"time"
)

type SqlBuild struct {
	ExeSql string
	ExeParam [] interface{}
	exeType string
	runtime time.Duration
	sqlPart *sqlInfo
	ShowSql string
	joinTmp *joinInfo
}


var whereFactory *whereInfo
var prefix string

func init () {
	config, _ := database.GetMysqlConfig()
	prefix = config.Prifex
	whereFactory = &whereInfo{}
}

func DB() *SqlBuild {
	return NewBuild()
}

func (this *SqlBuild) build() *SqlBuild {
	var sql string
	var params []interface{}
	switch this.exeType {
	case "SELECT":
		sql, params = this.sqlPart.BuildQuery()
	case "INSERT":
		sql, params = this.sqlPart.BuildInsert()
	case "UPDATE":
		sql, params = this.sqlPart.BuildUpdate()
	case "DELETE":
		sql, params = this.sqlPart.BuildDelete()
	case "ALTER" :
	default:

	}

	this.ExeSql = sql
	this.setShowSql(params)

	return this
}

func (this *SqlBuild) newBuild() *SqlBuild {
	build := NewBuild()

	return build.Table(this.sqlPart.table, this.sqlPart.alias)
}

func (this *SqlBuild) setShowSql(params []interface{}) {
	this.ShowSql = this.ExeSql

	for _, value := range params {
		switch val := value.(type) {
			case string:
				this.ShowSql = strings.Replace(this.ShowSql, "?", AddSingleSymbol(val), 1)
			case int:
				this.ShowSql = strings.Replace(this.ShowSql, "?", strconv.Itoa(val), 1)
			default:
				panic("param error")
		}
		this.ExeParam = append(this.ExeParam, value)
	}
}

func (this *SqlBuild) whereResult(result *whereInfo, err error) *SqlBuild {
	if err != nil {
		panic(err)
	}

	this.sqlPart.setWhere(result)
	return this
}

func (this *SqlBuild) Table(args... string) *SqlBuild {
	if len(args) == 0 {
		panic("Table param error")
	}
	table, alias := args[0], ""
	if len(args) > 1 && len(args[1]) > 0 {
		alias = args[1]
	}

	this.sqlPart.setTable(prefix + table, alias)
	return this
}

func (this *SqlBuild) JoinTable(table string) {
	if len(table) == 0 {
		panic("JoinTable param error")
	}
	var alias string = ""
	if strings.Count(table, " as ") == 1 {
		tmp := strings.Split(table, " as ")
		table = tmp[0]
		alias = tmp[1]
	}
	_ = this.Table(table, alias)
}

func (this *SqlBuild) Select(args... string) *SqlBuild {
	if len(args) > 0 {
		this.sqlPart.selectData = strings.Join(args, ",")
	} else {
		this.sqlPart.selectData = "*"
	}

	return this
}

func (this *SqlBuild) SelectRaw(sql string) *SqlBuild {
	if len(sql) == 0 {
		panic("SelectRaw param error")
	}
	this.sqlPart.selectData = sql
	return this
}

/***************************************  JOIN  **********************************************************/

func (this *SqlBuild) JoinFactory(joinType, table, thatRelationField, relationCondition, thisRelationField string) *SqlBuild {
	build := NewBuild()
	build.JoinTable(table)
	this.sqlPart.setJoin(build, joinType, thatRelationField, relationCondition, thisRelationField)

	return this
}

func (this *SqlBuild) Join(table, thatRelationField, relationCondition, thisRelationField string) *SqlBuild {
	return this.JoinFactory("Inner Join", table, thatRelationField, relationCondition, thisRelationField)
}

func (this *SqlBuild) LeftJoin(table, thatRelationField, relationCondition, thisRelationField string) *SqlBuild {
	return this.JoinFactory("Left Join", table, thatRelationField, relationCondition, thisRelationField)
}

func (this *SqlBuild) RightJoin(table, thatRelationField, relationCondition, thisRelationField string) *SqlBuild {
	return this.JoinFactory("Right Join", table, thatRelationField, relationCondition, thisRelationField)
}

func (this *SqlBuild) InnerJoin(table, thatRelationField, relationCondition, thisRelationField string) *SqlBuild {
	return this.JoinFactory("Inner Join", table, thatRelationField, relationCondition, thisRelationField)
}

func (this *SqlBuild) On(thatRelationField, relationCondition, thisRelationField string) *SqlBuild {
	this.joinTmp = &joinInfo{}
	this.joinTmp.JoinOn(thatRelationField, relationCondition, thisRelationField)
	return this
}

func (this *SqlBuild) JoinFuncFactory(joinType, table string, callback func(build *SqlBuild) *SqlBuild) *SqlBuild {
	build :=  NewBuild()
	build.JoinTable(table)
	_ = callback(build)
	this.joinTmp.buildInfo = build
	this.joinTmp.joinType = joinType
	this.sqlPart.joinData = append(this.sqlPart.joinData, this.joinTmp)
	return this
}

func (this *SqlBuild) LeftJoinFunc(table string, callback func(build *SqlBuild) *SqlBuild) *SqlBuild {
	return this.JoinFuncFactory("Left Join", table, callback)
}

func (this *SqlBuild) RightJoinFunc(table string, callback func(build *SqlBuild) *SqlBuild) *SqlBuild {
	return this.JoinFuncFactory("right Join", table, callback)
}

func (this *SqlBuild) InnerJoinFunc(table string, callback func(build *SqlBuild) *SqlBuild) *SqlBuild {
	return this.JoinFuncFactory("Inner Join", table, callback)
}

func (this *SqlBuild) JoinFunc(table string, callback func(build *SqlBuild) *SqlBuild) *SqlBuild {
	return this.JoinFuncFactory("Inner Join", table, callback)
}
/*************************************  JOIN  END  *******************************************************/

/***************************************  WHERE  **********************************************************/

func (this *SqlBuild) Where(args... interface{}) *SqlBuild {
	if len(args) == 0 {
		panic("Where param error")
	}
	args = append(args, 1)
	return this.whereResult( whereFactory.Where(FormatWhereParam(args...)...))
}

func (this *SqlBuild) OrWhere(args... interface{}) *SqlBuild {
	if len(args) == 0 {
		panic("Where param error")
	}
	args = append(args, 0)
	return this.whereResult( whereFactory.Where(FormatWhereParam(args...)...))
}

func (this *SqlBuild) WhereString(args... string) *SqlBuild {
	if len(args) == 0 {
		panic("Where param error")
	}
	args = append(args, "1")
	return this.whereResult( whereFactory.Where(args...))
}

func (this *SqlBuild) OrWhereString(args... string) *SqlBuild {
	if len(args) == 0 {
		panic("Where param error")
	}
	args = append(args, "0")
	return this.whereResult( whereFactory.Where(args...))
}

func (this *SqlBuild) WhereArray(arrayWhere [][] string) *SqlBuild {
	if len(arrayWhere) == 0 {
		panic("WhereArray param error")
	}
	return this.whereResult( whereFactory.WhereArray(arrayWhere, this.newBuild(), false))
}

func (this *SqlBuild) OrWhereArray(arrayWhere [][] string) *SqlBuild {
	if len(arrayWhere) == 0 {
		panic("WhereArray param error")
	}
	return this.whereResult( whereFactory.WhereArray(arrayWhere, this.newBuild(), true))
}

func (this *SqlBuild) WhereMap(mapWhere map[string] interface{}) *SqlBuild {
	if len(mapWhere) == 0 {
		panic("WhereMap param error")
	}
	return this.whereResult( whereFactory.WhereMap(mapWhere, this.newBuild(), false))
}

func (this *SqlBuild) OrWhereMap(mapWhere map[string] interface{}) *SqlBuild {
	if len(mapWhere) == 0 {
		panic("WhereMap param error")
	}
	return this.whereResult( whereFactory.WhereMap(mapWhere, this.newBuild(), true))
}

func (this *SqlBuild) WhereIn(field string, listValue [] interface{}) *SqlBuild {
	if len(field) == 0 || len(listValue) == 0 {
		panic("WhereIn param error")
	}
	return this.whereResult( whereFactory.WhereIn(field, listValue, false))
}

func (this *SqlBuild) OrWhereIn(field string, listValue [] interface{}) *SqlBuild {
	if len(field) == 0 || len(listValue) == 0 {
		panic("WhereIn param error")
	}
	return this.whereResult( whereFactory.WhereIn(field, listValue, true))
}

func (this *SqlBuild) WhereNotIn(field string, listValue [] interface{}) *SqlBuild {
	if len(field) == 0 || len(listValue) == 0 {
		panic("WhereNotIn param error")
	}
	return this.whereResult( whereFactory.WhereNotIn(field, listValue, false))
}

func (this *SqlBuild) OrWhereNotIn(field string, listValue [] interface{}) *SqlBuild {
	if len(field) == 0 || len(listValue) == 0 {
		panic("WhereNotIn param error")
	}
	return this.whereResult( whereFactory.WhereNotIn(field, listValue, true))
}

func (this *SqlBuild) WhereBetween(field string, interval [] interface{}) *SqlBuild {
	if len(field) == 0 || len(interval) != 2 {
		panic("WhereBetween param error")
	}
	return this.whereResult( whereFactory.WhereBetween(field, interval, false))
}

func (this *SqlBuild) OrWhereBetween(field string, interval [] interface{}) *SqlBuild {
	if len(field) == 0 || len(interval) != 2 {
		panic("WhereBetween param error")
	}
	return this.whereResult( whereFactory.WhereBetween(field, interval, true))
}

func (this *SqlBuild) WhereNotBetween(field string, interval [] interface{}) *SqlBuild {
	if len(field) == 0 || len(interval) != 2 {
		panic("WhereNotBetween param error")
	}
	return this.whereResult( whereFactory.WhereNotBetween(field, interval, false))
}

func (this *SqlBuild) OrWhereNotBetween(field string, interval [] interface{}) *SqlBuild {
	if len(field) == 0 || len(interval) != 2 {
		panic("WhereNotBetween param error")
	}
	return this.whereResult( whereFactory.WhereNotBetween(field, interval, true))
}

func (this *SqlBuild) WhereNull(field string) *SqlBuild {
	if len(field) == 0 {
		panic("WhereNull param error")
	}
	return this.whereResult( whereFactory.WhereNull(field, false))
}

func (this *SqlBuild) OrWhereNull(field string) *SqlBuild {
	if len(field) == 0 {
		panic("WhereNull param error")
	}
	return this.whereResult( whereFactory.WhereNull(field, true))
}

func (this *SqlBuild) WhereNotNull(field string) *SqlBuild {
	if len(field) == 0 {
		panic("WhereNotNull param error")
	}
	return this.whereResult( whereFactory.WhereNotNull(field, false))
}

func (this *SqlBuild) OrWhereNotNull(field string) *SqlBuild {
	if len(field) == 0 {
		panic("WhereNotNull param error")
	}
	return this.whereResult( whereFactory.WhereNotNull(field, true))
}

func (this *SqlBuild) WhereDate(field, date string) *SqlBuild {
	if len(field) == 0 || len(date) == 0 {
		panic("WhereDate param error")
	}
	return this.whereResult( whereFactory.WhereDate(field, date, false))
}

func (this *SqlBuild) OrWhereDate(field, date string) *SqlBuild {
	if len(field) == 0 || len(date) == 0 {
		panic("WhereDate param error")
	}
	return this.whereResult( whereFactory.WhereDate(field, date, true))
}

func (this *SqlBuild) WhereMonth(field, month string) *SqlBuild {
	if len(field) == 0 || len(month) == 0 {
		panic("WhereMonth param error")
	}
	return this.whereResult( whereFactory.WhereMonth(field, month, false))
}

func (this *SqlBuild) OrWhereMonth(field, month string) *SqlBuild {
	if len(field) == 0 || len(month) == 0 {
		panic("WhereMonth param error")
	}
	return this.whereResult( whereFactory.WhereMonth(field, month, true))
}

func (this *SqlBuild) WhereDay(field, day string) *SqlBuild {
	if len(field) == 0 || len(day) == 0 {
		panic("WhereDay param error")
	}
	return this.whereResult( whereFactory.WhereDay(field, day, false))
}

func (this *SqlBuild) OrWhereDay(field, day string) *SqlBuild {
	if len(field) == 0 || len(day) == 0 {
		panic("WhereDay param error")
	}
	return this.whereResult( whereFactory.WhereDay(field, day, true))
}

func (this *SqlBuild) WhereYear(field, year string) *SqlBuild {
	if len(field) == 0 || len(year) == 0 {
		panic("whereYear param error")
	}
	return this.whereResult( whereFactory.whereYear(field, year, false))
}

func (this *SqlBuild) OrWhereYear(field, year string) *SqlBuild {
	if len(field) == 0 || len(year) == 0 {
		panic("whereYear param error")
	}
	return this.whereResult( whereFactory.whereYear(field, year, true))
}

func (this *SqlBuild) WhereTime(field, condition, timestamp string) *SqlBuild {
	if len(field) == 0 || len(condition) == 0 || len(timestamp) == 0 {
		panic("WhereTime param error")
	}
	return this.whereResult( whereFactory.WhereTime(field, condition, timestamp, false))
}

func (this *SqlBuild) OrWhereTime(field, condition, timestamp string) *SqlBuild {
	if len(field) == 0 || len(condition) == 0 || len(timestamp) == 0 {
		panic("WhereTime param error")
	}
	return this.whereResult( whereFactory.WhereTime(field, condition, timestamp, true))
}

func (this *SqlBuild) WhereFunc(callback func(build *SqlBuild) *SqlBuild) *SqlBuild {
	return this.whereResult( whereFactory.WhereFunc(callback, this.newBuild(), false))
}

func (this *SqlBuild) OrWhereFunc(callback func(build *SqlBuild) *SqlBuild) *SqlBuild {
	return this.whereResult( whereFactory.WhereFunc(callback, this.newBuild(), true))
}

func (this *SqlBuild) WhereExists(callback func(build *SqlBuild) *SqlBuild) *SqlBuild {
	return this.whereResult( whereFactory.WhereExists("Exists", callback, false))
}

func (this *SqlBuild) OrWhereExists(callback func(build *SqlBuild) *SqlBuild) *SqlBuild {
	return this.whereResult( whereFactory.WhereExists("Exists", callback, true))
}

func (this *SqlBuild) WhereNotExists(field string, callback func(build *SqlBuild) *SqlBuild) *SqlBuild {
	return this.whereResult( whereFactory.WhereExists("Not Exists", callback, false))
}

func (this *SqlBuild) OrWhereNotExists(field string, callback func(build *SqlBuild) *SqlBuild) *SqlBuild {
	return this.whereResult( whereFactory.WhereExists("Not Exists", callback, true))
}

func (this *SqlBuild) WhereRaw(sql string) *SqlBuild {
	return this.whereResult( whereFactory.WhereRaw(sql, false))
}

func (this *SqlBuild) OrWhereRaw(sql string) *SqlBuild {
	return this.whereResult( whereFactory.WhereRaw(sql, true))
}

func (this *SqlBuild) When(boolean bool, callback func(build *SqlBuild) *SqlBuild) *SqlBuild {
	if boolean {
		return this.whereResult( whereFactory.WhereFunc(callback, this.newBuild(), false))
	}

	return this
}

func (this *SqlBuild) OrWhen(boolean bool, callback func(build *SqlBuild) *SqlBuild) *SqlBuild {
	if boolean {
		return this.whereResult( whereFactory.WhereFunc(callback, this.newBuild(), true))
	}

	return this
}

func (this *SqlBuild) WhenElse(boolean bool, trueCallback func(build *SqlBuild) *SqlBuild, falseCallback func(build *SqlBuild) *SqlBuild) *SqlBuild {
	if boolean {
		return this.whereResult( whereFactory.WhereFunc(trueCallback, this.newBuild(), false))
	} else {
		return this.whereResult( whereFactory.WhereFunc(falseCallback, this.newBuild(), false))
	}
}

func (this *SqlBuild) OrWhenElse(boolean bool, trueCallback func(build *SqlBuild) *SqlBuild, falseCallback func(build *SqlBuild) *SqlBuild) *SqlBuild {
	if boolean {
		return this.whereResult( whereFactory.WhereFunc(trueCallback, this.newBuild(), true))
	} else {
		return this.whereResult( whereFactory.WhereFunc(falseCallback, this.newBuild(), true))
	}
}

/*************************************  WHERE  END  *******************************************************/

/*************************************  Other  *******************************************************/
func (this *SqlBuild) OrderBy(args... string) *SqlBuild {
	if len(args) == 0 {
		panic("OrderBy param error")
	}
	for _, item := range args{
		this.sqlPart.orderData = append(this.sqlPart.orderData, item)
	}
	return this
}

func (this *SqlBuild) OrderByRaw(sql string) *SqlBuild {
	if len(sql) == 0 {
		panic("OrderByRaw param error")
	}
	this.sqlPart.orderData = append(this.sqlPart.orderData, sql)
	return this
}

func (this *SqlBuild) GroupBy(args... string) *SqlBuild {
	if len(args) == 0 {
		panic("GroupBy param error")
	}
	for _, item := range args{
		this.sqlPart.groupData = append(this.sqlPart.groupData, item)
	}
	return this
}

func (this *SqlBuild) GroupByRaw(sql string) *SqlBuild {
	if len(sql) == 0 {
		panic("GroupByRaw param error")
	}
	this.sqlPart.groupData = append(this.sqlPart.groupData, sql)
	return this
}

func (this *SqlBuild) Having(args... string) *SqlBuild {
	if len(args) == 0 {
		panic("Having param error")
	}
	return this.whereResult( whereFactory.Where(args...))
}

func (this *SqlBuild) HavingRaw(sql string) *SqlBuild {
	if len(sql) == 0 {
		panic("HavingRaw param error")
	}
	return this.whereResult( whereFactory.WhereRaw(sql, false))
}

func (this *SqlBuild) Offset(num int) *SqlBuild {
	this.sqlPart.offset = num
	return this
}

func (this *SqlBuild) Limit(num int) *SqlBuild {
	if num == 0 {
		panic("Limit param error")
	}
	this.sqlPart.limit = num
	return this
}

/*************************************  Other END *******************************************************/

/*************************************  SELECT *******************************************************/

func (this *SqlBuild) Get(args... string) *SqlBuild {
	if len(args) > 0 {
		_ = this.Select(args...)
	}

	this.exeType = "SELECT"
	return this.build()
}

func (this *SqlBuild) Value(field string) *SqlBuild {
	_ = this.Select(field)
	this.exeType = "SELECT"
	return this.build()
}

func (this *SqlBuild) First() *SqlBuild {
	_ = this.Limit(1)
	this.exeType = "SELECT"
	return this.build()
}

func (this *SqlBuild) PluckArray(field string) *SqlBuild {
	_ = this.Select(field)
	this.exeType = "SELECT"
	return this.build()
}

func (this *SqlBuild) PluckMap(field, value string) *SqlBuild {
	_ = this.Select(field, value)
	this.exeType = "SELECT"
	return this.build()
}

func (this *SqlBuild) Chunk(num int, callback func()) *SqlBuild {
	if num == 0 {
		panic("Chunk param error")
	}
	_ = this.Limit(num)
	this.exeType = "SELECT"
	return this.build()
}

func (this *SqlBuild) Count() *SqlBuild {
	_ = this.SelectRaw("COUNT(1) AS count")
	this.exeType = "SELECT"
	return this.build()
}

func (this *SqlBuild) Max(field string) *SqlBuild {
	if len(field) == 0 {
		panic("Max param error")
	}
	_ = this.SelectRaw("Max("+field+") AS max")
	this.exeType = "SELECT"
	return this.build()
}

func (this *SqlBuild) Sum(field string) *SqlBuild {
	if len(field) == 0 {
		panic("Sum param error")
	}
	_ = this.SelectRaw("Sum("+field+") AS sum")
	this.exeType = "SELECT"
	return this.build()
}

func (this *SqlBuild) Exists() *SqlBuild {
	_ = this.SelectRaw("COUNT(1) AS count")
	this.exeType = "SELECT"
	return this.build()
}

func (this *SqlBuild) DoesntExists() *SqlBuild {
	_ = this.SelectRaw("COUNT(1) AS count")
	this.exeType = "SELECT"
	return this.build()
}

func (this *SqlBuild) ToSql() string {
	return this.ShowSql
}

/*************************************  SELECT END *******************************************************/

func (this *SqlBuild) Insert(data map[string]interface{}) *SqlBuild {
	if len(data) == 0 {
		panic("Insert param error")
	}

	this.sqlPart.insertData = append(this.sqlPart.insertData, data)
	this.exeType = "INSERT"
	return this.build()
}

func (this *SqlBuild) MultiInsert(data []map[string]interface{}) *SqlBuild {
	if len(data) == 0 {
		panic("MultiInsert param error")
	}
	this.sqlPart.insertData = data
	this.exeType = "INSERT"
	return this.build()
}

func (this *SqlBuild) Update(data map[string]interface{}) *SqlBuild {
	if len(data) == 0 {
		panic("Update param error")
	}
	this.sqlPart.updateData = data
	this.exeType = "UPDATE"
	return this.build()
}

func (this *SqlBuild) Delete() *SqlBuild {
	this.exeType = "DELETE"
	return this.build()
}

