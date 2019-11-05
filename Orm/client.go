package Orm

import (
	"blog/Orm/clientPool"
	"blog/Orm/result"
	"blog/Orm/sqlBuild"
	"strconv"
)

type Manage struct {
	build *sqlBuild.SqlBuild
	toSql bool
}

func DB() *Manage {
	manage := &Manage{}
	manage.build = sqlBuild.NewBuild()
	return manage
}

func (this *Manage) LastInsertId(data map[string]interface{}) (int, error) {
	this.build.Insert(data)
	return this.Exec(true)
}

func (this *Manage) Insert(data map[string]interface{}) (int, error) {
	this.build.Insert(data)
	return this.Exec(false)
}

func (this *Manage) MultiInsert(data []map[string]interface{}) (int, error) {
	this.build.MultiInsert(data)
	return this.Exec(false)
}

func (this *Manage) Update(data map[string]interface{}) (int, error) {
	this.build.Update(data)
	return this.Exec(false)
}

func (this *Manage) Delete() (int, error) {
	this.build.Delete()
	return this.Exec(false)
}

func (this *Manage) Get(args... string) ([]map[string]string, error) {
	this.build.Get(args...)
	return this.Query()
}

func (this *Manage) Value(field string) (string, error) {
	this.build.Value(field)
	data, err := this.Query()
	if err != nil || len(data) == 0 {
		return "", err
	}
	return data[0][field], nil
}

func (this *Manage) First() (map[string]string, error) {
	this.build.First()
	data, err := this.Query()
	if err != nil || len(data) == 0 {
		return nil, err
	}
	return data[0], nil
}

func (this *Manage) PluckArray(field string) ([]string, error) {
	this.build.PluckArray(field)
	data, err := this.Query()
	if err != nil || len(data) == 0 {
		return nil, err
	}
	var res []string
	for _, item := range data {
		res = append(res, item[field])
	}
	return res, nil
}

func (this *Manage) PluckMap(field, value string) (map[string]string, error) {
	this.build.PluckMap(field, value)
	data, err := this.Query()
	if err != nil || len(data) == 0 {
		return nil, err
	}
	var res = make(map[string]string)
	for _, item := range data {
		res[item[field]] = item[value]
	}
	return res, nil
}

func (this *Manage) Chunk(num int, callback func()) {
	this.build.Chunk(num, callback)
}

func (this *Manage) Count() (int, error) {
	this.build.Count()
	data, err := this.Query()
	if err != nil || len(data) == 0 {
		return 0, err
	}
	count, _ := strconv.Atoi(data[0]["count"])
	return count, nil
}

func (this *Manage) Max(field string) (int, error) {
	this.build.Max(field)
	data, err := this.Query()
	if err != nil || len(data) == 0 {
		return 0, err
	}
	count, _ := strconv.Atoi(data[0]["max"])
	return count, nil
}

func (this *Manage) Sum(field string) (int, error) {
	this.build.Sum(field)
	data, err := this.Query()
	if err != nil || len(data) == 0 {
		return 0, err
	}
	count, _ := strconv.Atoi(data[0]["sum"])
	return count, nil
}

func (this *Manage) Exists() {
	this.build.Exists()
}

func (this *Manage) DoesntExists() {
	this.build.DoesntExists()
}

func (this *Manage) Query() ([]map[string]string, error) {
	client, err := clientPool.GetClient()
	if err != nil {
		return nil, err
	}
	stmt, err := client.Prepare(this.build.ExeSql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(this.build.ExeParam...)
	if err != nil {
		return nil, err
	}
	err = clientPool.CloseClient(client)
	if err != nil {
		return nil, err
	}

	return result.MakeResult(rows)
}

func (this *Manage) Exec(InsertId bool) (int, error) {
	client, err := clientPool.GetClient()
	if err != nil {
		return 0, err
	}
	stmt, err := client.Prepare(this.build.ExeSql)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	ret, err := stmt.Exec(this.build.ExeParam...)
	if err != nil {
		return 0, err
	}

	err = clientPool.CloseClient(client)
	if err != nil {
		return 0, err
	}

	var num int64
	if InsertId {
		num, err = ret.LastInsertId()
	} else {
		num, err = ret.RowsAffected()
	}
	if err != nil {
		return 0, err
	}

	return int(num), nil
}



func (this *Manage) LastInsertIdToSql(data map[string]interface{}) string {
	this.build.Insert(data)
	return this.build.ShowSql
}

func (this *Manage) InsertToSql(data map[string]interface{}) string {
	this.build.Insert(data)
	return this.build.ShowSql
}

func (this *Manage) MultiInsertToSql(data []map[string]interface{}) string {
	this.build.MultiInsert(data)
	return this.build.ShowSql
}

func (this *Manage) UpdateToSql(data map[string]interface{}) string {
	this.build.Update(data)
	return this.build.ShowSql
}

func (this *Manage) DeleteToSql() string {
	this.build.Delete()
	return this.build.ShowSql
}
func (this *Manage) GetToSql(args... string) string {
	this.build.Get(args...)
	return this.build.ShowSql
}
func (this *Manage) ValueToSql(field string) string {
	this.build.Value(field)
	return this.build.ShowSql
}
func (this *Manage) FirstToSql(args... string) string {
	this.build.First()
	return this.build.ShowSql
}
func (this *Manage) PluckArrayToSql(field string) string {
	this.build.PluckArray(field)
	return this.build.ShowSql
}
func (this *Manage) PluckMapToSql(field, value string) string {
	this.build.PluckMap(field, value)
	return this.build.ShowSql
}
func (this *Manage) CountToSql() string {
	this.build.Count()
	return this.build.ShowSql
}
func (this *Manage) MaxToSql(field string) string {
	this.build.Max(field)
	return this.build.ShowSql
}
func (this *Manage) SumToSql(field string) string {
	this.build.Sum(field)
	return this.build.ShowSql
}
func (this *Manage) ChunkToSql(num int) string {
	this.build.Limit(num)
	this.build.Get()
	return this.build.ShowSql
}
func (this *Manage) Table(args... string) *Manage {this.build.Table(args...);return this}
func (this *Manage) Select(args... string) *Manage {this.build.Select(args...);return this}
func (this *Manage) SelectRaw(sql string) *Manage {this.build.SelectRaw(sql);return this}
func (this *Manage) Join(table, thatRelationField, relationCondition, thisRelationField string) *Manage {this.build.Join(table, thatRelationField, relationCondition, thisRelationField);return this}
func (this *Manage) LeftJoin(table, thatRelationField, relationCondition, thisRelationField string) *Manage {this.build.LeftJoin(table, thatRelationField, relationCondition, thisRelationField);return this}
func (this *Manage) RightJoin(table, thatRelationField, relationCondition, thisRelationField string) *Manage {this.build.RightJoin(table, thatRelationField, relationCondition, thisRelationField);return this}
func (this *Manage) InnerJoin(table, thatRelationField, relationCondition, thisRelationField string) *Manage {this.build.InnerJoin(table, thatRelationField, relationCondition, thisRelationField);return this}
func (this *Manage) On(thatRelationField, relationCondition, thisRelationField string) *Manage {this.build.On(thatRelationField, relationCondition, thisRelationField);return this}
func (this *Manage) LeftJoinFunc(table string, callback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manage {this.build.LeftJoinFunc(table, callback);return this}
func (this *Manage) RightJoinFunc(table string, callback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manage {this.build.RightJoinFunc(table, callback);return this}
func (this *Manage) InnerJoinFunc(table string, callback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manage {this.build.InnerJoinFunc(table, callback);return this}
func (this *Manage) JoinFunc(table string, callback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manage {this.build.JoinFunc(table, callback);return this}
func (this *Manage) Where(args... interface{}) *Manage {this.build.Where(args...);return this}
func (this *Manage) OrWhere(args... interface{}) *Manage {this.build.OrWhere(args);return this}
func (this *Manage) WhereArray(arrayWhere [][] string) *Manage {this.build.WhereArray(arrayWhere);return this}
func (this *Manage) OrWhereArray(arrayWhere [][] string) *Manage {this.build.OrWhereArray(arrayWhere);return this}
func (this *Manage) WhereMap(mapWhere map[string] interface{}) *Manage {this.build.WhereMap(mapWhere);return this}
func (this *Manage) OrWhereMap(mapWhere map[string] interface{}) *Manage {this.build.OrWhereMap(mapWhere);return this}
func (this *Manage) WhereIn(field string, listValue [] interface{}) *Manage {this.build.WhereIn(field, listValue);return this}
func (this *Manage) OrWhereIn(field string, listValue [] interface{}) *Manage {this.build.OrWhereIn(field, listValue);return this}
func (this *Manage) WhereNotIn(field string, listValue [] interface{}) *Manage {this.build.WhereNotIn(field, listValue);return this}
func (this *Manage) OrWhereNotIn(field string, listValue [] interface{}) *Manage {this.build.OrWhereNotIn(field, listValue);return this}
func (this *Manage) WhereBetween(field string, interval [] interface{}) *Manage {this.build.WhereBetween(field, interval);return this}
func (this *Manage) OrWhereBetween(field string, interval [] interface{}) *Manage {this.build.OrWhereBetween(field, interval);return this}
func (this *Manage) WhereNotBetween(field string, interval [] interface{}) *Manage {this.build.WhereNotBetween(field, interval);return this}
func (this *Manage) OrWhereNotBetween(field string, interval [] interface{}) *Manage {this.build.OrWhereNotBetween(field, interval);return this}
func (this *Manage) WhereNull(field string) *Manage {this.build.WhereNull(field);return this}
func (this *Manage) OrWhereNull(field string) *Manage {this.build.OrWhereNull(field);return this}
func (this *Manage) WhereNotNull(field string) *Manage {this.build.WhereNotNull(field);return this}
func (this *Manage) OrWhereNotNull(field string) *Manage {this.build.OrWhereNotNull(field);return this}
func (this *Manage) WhereDate(field, date string) *Manage {this.build.WhereDate(field, date);return this}
func (this *Manage) OrWhereDate(field, date string) *Manage {this.build.OrWhereDate(field, date);return this}
func (this *Manage) WhereMonth(field, month string) *Manage {this.build.WhereMonth(field, month);return this}
func (this *Manage) OrWhereMonth(field, month string) *Manage {this.build.OrWhereMonth(field, month);return this}
func (this *Manage) WhereDay(field, day string) *Manage {this.build.WhereDay(field, day);return this}
func (this *Manage) OrWhereDay(field, day string) *Manage {this.build.OrWhereDay(field, day);return this}
func (this *Manage) WhereYear(field, year string) *Manage {this.build.WhereYear(field, year);return this}
func (this *Manage) OrWhereYear(field, year string) *Manage {this.build.OrWhereYear(field, year);return this}
func (this *Manage) WhereTime(field, condition, timestamp string) *Manage {this.build.WhereTime(field, condition, timestamp);return this}
func (this *Manage) OrWhereTime(field, condition, timestamp string) *Manage {this.build.OrWhereTime(field, condition, timestamp);return this}
func (this *Manage) WhereFunc(callback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manage {this.build.WhereFunc(callback);return this}
func (this *Manage) OrWhereFunc(callback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manage {this.build.OrWhereFunc(callback);return this}
func (this *Manage) WhereExists(callback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manage {this.build.WhereExists(callback);return this}
func (this *Manage) OrWhereExists(callback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manage {this.build.OrWhereExists(callback);return this}
func (this *Manage) WhereNotExists(field string, callback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manage {this.build.WhereNotExists(field, callback);return this}
func (this *Manage) OrWhereNotExists(field string, callback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manage {this.build.OrWhereNotExists(field, callback);return this}
func (this *Manage) WhereRaw(sql string) *Manage {this.build.WhereRaw(sql);return this}
func (this *Manage) OrWhereRaw(sql string) *Manage {this.build.OrWhereRaw(sql);return this}
func (this *Manage) When(boolean bool, callback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manage {this.build.When(boolean, callback);return this}
func (this *Manage) OrWhen(boolean bool, callback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manage {this.build.OrWhen(boolean, callback);return this}
func (this *Manage) WhenElse(boolean bool, trueCallback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild, falseCallback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manage {this.build.WhenElse(boolean, trueCallback, falseCallback);return this}
func (this *Manage) OrWhenElse(boolean bool, trueCallback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild, falseCallback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manage {this.build.OrWhenElse(boolean, trueCallback, falseCallback);return this}
func (this *Manage) OrderBy(args... string) *Manage {this.build.OrderBy(args...);return this}
func (this *Manage) OrderByRaw(sql string) *Manage {this.build.OrderByRaw(sql);return this}
func (this *Manage) GroupBy(args... string) *Manage {this.build.GroupBy(args...);return this}
func (this *Manage) GroupByRaw(sql string) *Manage {this.build.GroupByRaw(sql);return this}
func (this *Manage) Having(args... string) *Manage {this.build.Having(args...);return this}
func (this *Manage) HavingRaw(sql string) *Manage {this.build.HavingRaw(sql);return this}
func (this *Manage) Offset(num int) *Manage {this.build.Offset(num);return this}
func (this *Manage) Limit(num int) *Manage {this.build.Limit(num);return this}