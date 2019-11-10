package Orm

import (
	"blog/Orm/clientPool"
	"blog/Orm/result"
	"blog/Orm/sqlBuild"
	"database/sql"
	"errors"
	"strconv"
)

type Manager struct {
	build *sqlBuild.SqlBuild
	client *sql.DB
	tx *sql.Tx
}

var log func(build *sqlBuild.SqlBuild)

func SetDBLog(callback func(build *sqlBuild.SqlBuild)) {
	log = callback
}

func DB() *Manager {
	manage := &Manager{}
	return manage
}

func DbBegin() *Manager {
	client,err := clientPool.GetClient()
	if err != nil {
		panic(err)
	}
	tx, err := client.Begin()
	if err != nil {
		panic(err)
	}
	return &Manager{
		build: sqlBuild.NewBuild(),
		client: client,
		tx: tx,
	}
}

func (this *Manager) DbCommit() (error) {
	if this.tx == nil {
		return errors.New("Please begin transaction by DbBegin() first")
	}
	if err := this.tx.Commit(); err != nil {
		return err
	}
	err := clientPool.CloseClient(this.client)
	if err != nil {
		return err
	}
	this.tx = nil
	this.client = nil
	return nil
}

func (this *Manager) DbRollBack() (error) {
	if this.tx == nil {
		return errors.New("Please begin transaction by DbBegin() first")
	}
	err := this.tx.Rollback()
	if err != sql.ErrTxDone && err != nil {
		return err
	}
	err = clientPool.CloseClient(this.client)
	if err != nil {
		return err
	}
	this.tx = nil
	this.client = nil
	return nil
}

func (this *Manager) LastInsertId(data map[string]interface{}) (int, error) {
	this.build.Insert(data)
	return this.exec(true)
}

func (this *Manager) Insert(data map[string]interface{}) (int, error) {
	this.build.Insert(data)
	return this.exec(false)
}

func (this *Manager) MultiInsert(data []map[string]interface{}) (int, error) {
	this.build.MultiInsert(data)
	return this.exec(false)
}

func (this *Manager) Update(data map[string]interface{}) (int, error) {
	this.build.Update(data)
	return this.exec(false)
}

func (this *Manager) Delete() (int, error) {
	this.build.Delete()
	return this.exec(false)
}

func (this *Manager) Get(args... string) ([]map[string]string, error) {
	this.build.Get(args...)
	return this.query()
}

func (this *Manager) Value(field string) (string, error) {
	this.build.Value(field)
	data, err := this.query()
	if err != nil || len(data) == 0 {
		return "", err
	}
	return data[0][field], nil
}

func (this *Manager) First() (map[string]string, error) {
	this.build.First()
	data, err := this.query()
	if err != nil || len(data) == 0 {
		return nil, err
	}
	return data[0], nil
}

func (this *Manager) PluckArray(field string) ([]string, error) {
	this.build.PluckArray(field)
	data, err := this.query()
	if err != nil || len(data) == 0 {
		return nil, err
	}
	var res []string
	for _, item := range data {
		res = append(res, item[field])
	}
	return res, nil
}

func (this *Manager) PluckMap(field, value string) (map[string]string, error) {
	this.build.PluckMap(field, value)
	data, err := this.query()
	if err != nil || len(data) == 0 {
		return nil, err
	}
	var res = make(map[string]string)
	for _, item := range data {
		res[item[field]] = item[value]
	}
	return res, nil
}

// Todo Chunk
func (this *Manager) Chunk(num int, callback func()) {
	this.build.Chunk(num, callback)
}

func (this *Manager) Count() (int, error) {
	this.build.Count()
	data, err := this.query()
	if err != nil || len(data) == 0 {
		return 0, err
	}
	count, _ := strconv.Atoi(data[0]["count"])
	return count, nil
}

func (this *Manager) Max(field string) (int, error) {
	this.build.Max(field)
	data, err := this.query()
	if err != nil || len(data) == 0 {
		return 0, err
	}
	count, _ := strconv.Atoi(data[0]["max"])
	return count, nil
}

func (this *Manager) Sum(field string) (int, error) {
	this.build.Sum(field)
	data, err := this.query()
	if err != nil || len(data) == 0 {
		return 0, err
	}
	count, _ := strconv.Atoi(data[0]["sum"])
	return count, nil
}

// Todo Exists
func (this *Manager) Exists() {
	this.build.Exists()
}

// Todo DoesntExists
func (this *Manager) DoesntExists() {
	this.build.DoesntExists()
}

func (this *Manager) query() ([]map[string]string, error) {
	var rows *sql.Rows
	var err error
	if this.tx != nil {
		rows, err = this.tx.Query(this.build.ExeSql, this.build.ExeParam...)
	} else {
		client, err := clientPool.GetClient()
		if err != nil {
			return nil, err
		}
		stmt, err := client.Prepare(this.build.ExeSql)
		if err != nil {
			return nil, err
		}
		defer stmt.Close()
		rows, err = stmt.Query(this.build.ExeParam...)
		defer clientPool.CloseClient(client)
	}

	if err != nil {
		return nil, err
	}

	return result.MakeResult(rows)
}

func (this *Manager) exec(InsertId bool) (int, error) {
	var ret sql.Result
	var err error
	if this.tx != nil {
		ret, err = this.tx.Exec(this.build.ExeSql, this.build.ExeParam...)
	} else {
		client, err := clientPool.GetClient()
		if err != nil {
			return 0, err
		}
		stmt, err := client.Prepare(this.build.ExeSql)
		if err != nil {
			return 0, err
		}
		defer stmt.Close()
		ret, err = stmt.Exec(this.build.ExeParam...)
		defer clientPool.CloseClient(client)
	}

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



func (this *Manager) LastInsertIdToSql(data map[string]interface{}) string {
	this.build.Insert(data)
	return this.build.ShowSql
}

func (this *Manager) InsertToSql(data map[string]interface{}) string {
	this.build.Insert(data)
	return this.build.ShowSql
}

func (this *Manager) MultiInsertToSql(data []map[string]interface{}) string {
	this.build.MultiInsert(data)
	return this.build.ShowSql
}

func (this *Manager) UpdateToSql(data map[string]interface{}) string {
	this.build.Update(data)
	return this.build.ShowSql
}

func (this *Manager) DeleteToSql() string {
	this.build.Delete()
	return this.build.ShowSql
}
func (this *Manager) GetToSql(args... string) string {
	this.build.Get(args...)
	return this.build.ShowSql
}
func (this *Manager) ValueToSql(field string) string {
	this.build.Value(field)
	return this.build.ShowSql
}
func (this *Manager) FirstToSql(args... string) string {
	this.build.First()
	return this.build.ShowSql
}
func (this *Manager) PluckArrayToSql(field string) string {
	this.build.PluckArray(field)
	return this.build.ShowSql
}
func (this *Manager) PluckMapToSql(field, value string) string {
	this.build.PluckMap(field, value)
	return this.build.ShowSql
}
func (this *Manager) CountToSql() string {
	this.build.Count()
	return this.build.ShowSql
}
func (this *Manager) MaxToSql(field string) string {
	this.build.Max(field)
	return this.build.ShowSql
}
func (this *Manager) SumToSql(field string) string {
	this.build.Sum(field)
	return this.build.ShowSql
}
func (this *Manager) ChunkToSql(num int) string {
	this.build.Limit(num)
	this.build.Get()
	return this.build.ShowSql
}
func (this *Manager) Table(args... string) *Manager {this.build = sqlBuild.NewBuild();this.build.Table(args...);return this}
func (this *Manager) Select(args... string) *Manager {this.build.Select(args...);return this}
func (this *Manager) SelectRaw(sql string) *Manager {this.build.SelectRaw(sql);return this}
func (this *Manager) Join(table, thatRelationField, relationCondition, thisRelationField string) *Manager {this.build.Join(table, thatRelationField, relationCondition, thisRelationField);return this}
func (this *Manager) LeftJoin(table, thatRelationField, relationCondition, thisRelationField string) *Manager {this.build.LeftJoin(table, thatRelationField, relationCondition, thisRelationField);return this}
func (this *Manager) RightJoin(table, thatRelationField, relationCondition, thisRelationField string) *Manager {this.build.RightJoin(table, thatRelationField, relationCondition, thisRelationField);return this}
func (this *Manager) InnerJoin(table, thatRelationField, relationCondition, thisRelationField string) *Manager {this.build.InnerJoin(table, thatRelationField, relationCondition, thisRelationField);return this}
func (this *Manager) On(thatRelationField, relationCondition, thisRelationField string) *Manager {this.build.On(thatRelationField, relationCondition, thisRelationField);return this}
func (this *Manager) LeftJoinFunc(table string, callback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manager {this.build.LeftJoinFunc(table, callback);return this}
func (this *Manager) RightJoinFunc(table string, callback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manager {this.build.RightJoinFunc(table, callback);return this}
func (this *Manager) InnerJoinFunc(table string, callback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manager {this.build.InnerJoinFunc(table, callback);return this}
func (this *Manager) JoinFunc(table string, callback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manager {this.build.JoinFunc(table, callback);return this}
func (this *Manager) Where(args... interface{}) *Manager {this.build.Where(args...);return this}
func (this *Manager) OrWhere(args... interface{}) *Manager {this.build.OrWhere(args);return this}
func (this *Manager) WhereArray(arrayWhere [][] string) *Manager {this.build.WhereArray(arrayWhere);return this}
func (this *Manager) OrWhereArray(arrayWhere [][] string) *Manager {this.build.OrWhereArray(arrayWhere);return this}
func (this *Manager) WhereMap(mapWhere map[string] interface{}) *Manager {this.build.WhereMap(mapWhere);return this}
func (this *Manager) OrWhereMap(mapWhere map[string] interface{}) *Manager {this.build.OrWhereMap(mapWhere);return this}
func (this *Manager) WhereIn(field string, listValue [] interface{}) *Manager {this.build.WhereIn(field, listValue);return this}
func (this *Manager) OrWhereIn(field string, listValue [] interface{}) *Manager {this.build.OrWhereIn(field, listValue);return this}
func (this *Manager) WhereNotIn(field string, listValue [] interface{}) *Manager {this.build.WhereNotIn(field, listValue);return this}
func (this *Manager) OrWhereNotIn(field string, listValue [] interface{}) *Manager {this.build.OrWhereNotIn(field, listValue);return this}
func (this *Manager) WhereBetween(field string, interval [] interface{}) *Manager {this.build.WhereBetween(field, interval);return this}
func (this *Manager) OrWhereBetween(field string, interval [] interface{}) *Manager {this.build.OrWhereBetween(field, interval);return this}
func (this *Manager) WhereNotBetween(field string, interval [] interface{}) *Manager {this.build.WhereNotBetween(field, interval);return this}
func (this *Manager) OrWhereNotBetween(field string, interval [] interface{}) *Manager {this.build.OrWhereNotBetween(field, interval);return this}
func (this *Manager) WhereNull(field string) *Manager {this.build.WhereNull(field);return this}
func (this *Manager) OrWhereNull(field string) *Manager {this.build.OrWhereNull(field);return this}
func (this *Manager) WhereNotNull(field string) *Manager {this.build.WhereNotNull(field);return this}
func (this *Manager) OrWhereNotNull(field string) *Manager {this.build.OrWhereNotNull(field);return this}
func (this *Manager) WhereDate(field, date string) *Manager {this.build.WhereDate(field, date);return this}
func (this *Manager) OrWhereDate(field, date string) *Manager {this.build.OrWhereDate(field, date);return this}
func (this *Manager) WhereMonth(field, month string) *Manager {this.build.WhereMonth(field, month);return this}
func (this *Manager) OrWhereMonth(field, month string) *Manager {this.build.OrWhereMonth(field, month);return this}
func (this *Manager) WhereDay(field, day string) *Manager {this.build.WhereDay(field, day);return this}
func (this *Manager) OrWhereDay(field, day string) *Manager {this.build.OrWhereDay(field, day);return this}
func (this *Manager) WhereYear(field, year string) *Manager {this.build.WhereYear(field, year);return this}
func (this *Manager) OrWhereYear(field, year string) *Manager {this.build.OrWhereYear(field, year);return this}
func (this *Manager) WhereTime(field, condition, timestamp string) *Manager {this.build.WhereTime(field, condition, timestamp);return this}
func (this *Manager) OrWhereTime(field, condition, timestamp string) *Manager {this.build.OrWhereTime(field, condition, timestamp);return this}
func (this *Manager) WhereFunc(callback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manager {this.build.WhereFunc(callback);return this}
func (this *Manager) OrWhereFunc(callback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manager {this.build.OrWhereFunc(callback);return this}
func (this *Manager) WhereExists(callback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manager {this.build.WhereExists(callback);return this}
func (this *Manager) OrWhereExists(callback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manager {this.build.OrWhereExists(callback);return this}
func (this *Manager) WhereNotExists(field string, callback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manager {this.build.WhereNotExists(field, callback);return this}
func (this *Manager) OrWhereNotExists(field string, callback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manager {this.build.OrWhereNotExists(field, callback);return this}
func (this *Manager) WhereRaw(sql string) *Manager {this.build.WhereRaw(sql);return this}
func (this *Manager) OrWhereRaw(sql string) *Manager {this.build.OrWhereRaw(sql);return this}
func (this *Manager) When(boolean bool, callback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manager {this.build.When(boolean, callback);return this}
func (this *Manager) OrWhen(boolean bool, callback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manager {this.build.OrWhen(boolean, callback);return this}
func (this *Manager) WhenElse(boolean bool, trueCallback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild, falseCallback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manager {this.build.WhenElse(boolean, trueCallback, falseCallback);return this}
func (this *Manager) OrWhenElse(boolean bool, trueCallback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild, falseCallback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manager {this.build.OrWhenElse(boolean, trueCallback, falseCallback);return this}
func (this *Manager) OrderBy(args... string) *Manager {this.build.OrderBy(args...);return this}
func (this *Manager) OrderByRaw(sql string) *Manager {this.build.OrderByRaw(sql);return this}
func (this *Manager) GroupBy(args... string) *Manager {this.build.GroupBy(args...);return this}
func (this *Manager) GroupByRaw(sql string) *Manager {this.build.GroupByRaw(sql);return this}
func (this *Manager) Having(args... string) *Manager {this.build.Having(args...);return this}
func (this *Manager) HavingRaw(sql string) *Manager {this.build.HavingRaw(sql);return this}
func (this *Manager) Offset(num int) *Manager {this.build.Offset(num);return this}
func (this *Manager) Limit(num int) *Manager {this.build.Limit(num);return this}