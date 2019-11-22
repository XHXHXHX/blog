package orm

import (
	"blog/orm/clientPool"
	"blog/orm/result"
	"blog/orm/sqlBuild"
	"blog/library/log"
	"database/sql"
	"errors"
	"github.com/sirupsen/logrus"
	"strconv"
)

type Manager struct {
	build *sqlBuild.SqlBuild
	client *sql.DB
	tx *sql.Tx
	model interface{}
}

func SetErrorLog(err error, build *sqlBuild.SqlBuild) {
	log.New().WithFields(logrus.Fields{
		"sql" : build.ToSql(),
		"error": err,
	})
}

func DB() *Manager {
	managege := &Manager{}
	return managege
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

func (manage *Manager) DbCommit() (error) {
	if manage.tx == nil {
		return errors.New("Please begin transaction by DbBegin() first")
	}
	if err := manage.tx.Commit(); err != nil {
		return err
	}
	err := clientPool.CloseClient(manage.client)
	if err != nil {
		return err
	}
	manage.tx = nil
	manage.client = nil
	return nil
}

func (manage *Manager) DbRollBack() (error) {
	if manage.tx == nil {
		return errors.New("Please begin transaction by DbBegin() first")
	}
	err := manage.tx.Rollback()
	if err != sql.ErrTxDone && err != nil {
		return err
	}
	err = clientPool.CloseClient(manage.client)
	if err != nil {
		return err
	}
	manage.tx = nil
	manage.client = nil
	return nil
}

func (manage *Manager) LastInsertId(data map[string]interface{}) (int, error) {
	manage.build.Insert(data)
	return manage.exec(true)
}

func (manage *Manager) Insert(data map[string]interface{}) (int, error) {
	manage.build.Insert(data)
	return manage.exec(false)
}

func (manage *Manager) MultiInsert(data []map[string]interface{}) (int, error) {
	manage.build.MultiInsert(data)
	return manage.exec(false)
}

func (manage *Manager) Update(data map[string]interface{}) (int, error) {
	manage.build.Update(data)
	return manage.exec(false)
}

func (manage *Manager) Delete() (int, error) {
	manage.build.Delete()
	return manage.exec(false)
}

func (manage *Manager) Get(args... string) ([]map[string]string, error) {
	manage.build.Get(args...)
	return manage.query()
}

func (manage *Manager) GetModel(model interface{}, args... string) (error) {
	manage.build.Get(args...)
	res, err :=  manage.query()
	if err != nil {
		return err
	}

	return result.ModelResult(model, res)
}

func (manage *Manager) Value(field string) (string, error) {
	manage.build.Value(field)
	data, err := manage.query()
	if err != nil || len(data) == 0 {
		return "", err
	}
	return data[0][field], nil
}

func (manage *Manager) First() (map[string]string, error) {
	manage.build.First()
	data, err := manage.query()
	if err != nil || len(data) == 0 {
		return nil, err
	}
	return data[0], nil
}

func (manage *Manager) FirstModel(model interface{}) (error) {
	manage.build.First()
	res, err := manage.query()
	if err != nil || len(res) == 0 {
		return err
	}

	return result.ModelResult(model, res)
}

func (manage *Manager) PluckArray(field string) ([]string, error) {
	manage.build.PluckArray(field)
	data, err := manage.query()
	if err != nil || len(data) == 0 {
		return nil, err
	}
	var res []string
	for _, item := range data {
		res = append(res, item[field])
	}
	return res, nil
}

func (manage *Manager) PluckMap(field, value string) (map[string]string, error) {
	manage.build.PluckMap(field, value)
	data, err := manage.query()
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
func (manage *Manager) Chunk(num int, callback func()) {
	manage.build.Chunk(num, callback)
}

func (manage *Manager) Count() (int, error) {
	manage.build.Count()
	data, err := manage.query()
	if err != nil || len(data) == 0 {
		return 0, err
	}
	count, _ := strconv.Atoi(data[0]["count"])
	return count, nil
}

func (manage *Manager) Max(field string) (int, error) {
	manage.build.Max(field)
	data, err := manage.query()
	if err != nil || len(data) == 0 {
		return 0, err
	}
	count, _ := strconv.Atoi(data[0]["max"])
	return count, nil
}

func (manage *Manager) Sum(field string) (int, error) {
	manage.build.Sum(field)
	data, err := manage.query()
	if err != nil || len(data) == 0 {
		return 0, err
	}
	count, _ := strconv.Atoi(data[0]["sum"])
	return count, nil
}

// Todo Exists
func (manage *Manager) Exists() {
	manage.build.Exists()
}

// Todo DoesntExists
func (manage *Manager) DoesntExists() {
	manage.build.DoesntExists()
}

func (manage *Manager) query() ([]map[string]string, error) {
	var rows *sql.Rows
	var err error
	if manage.tx != nil {
		rows, err = manage.tx.Query(manage.build.ExeSql, manage.build.ExeParam...)
	} else {
		client, err := clientPool.GetClient()
		if err != nil {
			return nil, err
		}
		stmt, err := client.Prepare(manage.build.ExeSql)
		if err != nil {
			SetErrorLog(err, manage.build)
			return nil, err
		}
		defer stmt.Close()
		rows, err = stmt.Query(manage.build.ExeParam...)
		defer clientPool.CloseClient(client)
	}

	if err != nil {
		SetErrorLog(err, manage.build)
		return nil, err
	}

	return result.MakeResult(rows)
}

func (manage *Manager) exec(InsertId bool) (int, error) {
	var ret sql.Result
	var err error
	if manage.tx != nil {
		ret, err = manage.tx.Exec(manage.build.ExeSql, manage.build.ExeParam...)
	} else {
		client, err := clientPool.GetClient()
		if err != nil {
			return 0, err
		}
		stmt, err := client.Prepare(manage.build.ExeSql)
		if err != nil {
			SetErrorLog(err, manage.build)
			return 0, err
		}
		defer stmt.Close()
		ret, err = stmt.Exec(manage.build.ExeParam...)
		defer clientPool.CloseClient(client)
	}

	if err != nil {
		SetErrorLog(err, manage.build)
		return 0, err
	}

	var num int64
	if InsertId {
		num, err = ret.LastInsertId()
	} else {
		num, err = ret.RowsAffected()
	}
	if err != nil {
		SetErrorLog(err, manage.build)
		return 0, err
	}

	return int(num), nil
}

func (manage *Manager) LastInsertIdToSql(data map[string]interface{}) string {
	manage.build.Insert(data)
	return manage.build.ShowSql
}

func (manage *Manager) InsertToSql(data map[string]interface{}) string {
	manage.build.Insert(data)
	return manage.build.ShowSql
}

func (manage *Manager) MultiInsertToSql(data []map[string]interface{}) string {
	manage.build.MultiInsert(data)
	return manage.build.ShowSql
}

func (manage *Manager) UpdateToSql(data map[string]interface{}) string {
	manage.build.Update(data)
	return manage.build.ShowSql
}

func (manage *Manager) DeleteToSql() string {
	manage.build.Delete()
	return manage.build.ShowSql
}
func (manage *Manager) GetToSql(args... string) string {
	manage.build.Get(args...)
	return manage.build.ShowSql
}
func (manage *Manager) ValueToSql(field string) string {
	manage.build.Value(field)
	return manage.build.ShowSql
}
func (manage *Manager) FirstToSql(args... string) string {
	manage.build.First()
	return manage.build.ShowSql
}
func (manage *Manager) PluckArrayToSql(field string) string {
	manage.build.PluckArray(field)
	return manage.build.ShowSql
}
func (manage *Manager) PluckMapToSql(field, value string) string {
	manage.build.PluckMap(field, value)
	return manage.build.ShowSql
}
func (manage *Manager) CountToSql() string {
	manage.build.Count()
	return manage.build.ShowSql
}
func (manage *Manager) MaxToSql(field string) string {
	manage.build.Max(field)
	return manage.build.ShowSql
}
func (manage *Manager) SumToSql(field string) string {
	manage.build.Sum(field)
	return manage.build.ShowSql
}
func (manage *Manager) ChunkToSql(num int) string {
	manage.build.Limit(num)
	manage.build.Get()
	return manage.build.ShowSql
}
func (manage *Manager) Table(args... string) *Manager {manage.build = sqlBuild.NewBuild();manage.build.Table(args...);return manage}
func (manage *Manager) Select(args... string) *Manager {manage.build.Select(args...);return manage}
func (manage *Manager) SelectRaw(sql string) *Manager {manage.build.SelectRaw(sql);return manage}
func (manage *Manager) Join(table, thatRelationField, relationCondition, manageRelationField string) *Manager {manage.build.Join(table, thatRelationField, relationCondition, manageRelationField);return manage}
func (manage *Manager) LeftJoin(table, thatRelationField, relationCondition, manageRelationField string) *Manager {manage.build.LeftJoin(table, thatRelationField, relationCondition, manageRelationField);return manage}
func (manage *Manager) RightJoin(table, thatRelationField, relationCondition, manageRelationField string) *Manager {manage.build.RightJoin(table, thatRelationField, relationCondition, manageRelationField);return manage}
func (manage *Manager) InnerJoin(table, thatRelationField, relationCondition, manageRelationField string) *Manager {manage.build.InnerJoin(table, thatRelationField, relationCondition, manageRelationField);return manage}
func (manage *Manager) On(thatRelationField, relationCondition, manageRelationField string) *Manager {manage.build.On(thatRelationField, relationCondition, manageRelationField);return manage}
func (manage *Manager) LeftJoinFunc(table string, callback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manager {manage.build.LeftJoinFunc(table, callback);return manage}
func (manage *Manager) RightJoinFunc(table string, callback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manager {manage.build.RightJoinFunc(table, callback);return manage}
func (manage *Manager) InnerJoinFunc(table string, callback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manager {manage.build.InnerJoinFunc(table, callback);return manage}
func (manage *Manager) JoinFunc(table string, callback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manager {manage.build.JoinFunc(table, callback);return manage}
func (manage *Manager) Where(args... interface{}) *Manager {manage.build.Where(args...);return manage}
func (manage *Manager) OrWhere(args... interface{}) *Manager {manage.build.OrWhere(args);return manage}
func (manage *Manager) WhereArray(arrayWhere [][] string) *Manager {manage.build.WhereArray(arrayWhere);return manage}
func (manage *Manager) OrWhereArray(arrayWhere [][] string) *Manager {manage.build.OrWhereArray(arrayWhere);return manage}
func (manage *Manager) WhereMap(mapWhere map[string] interface{}) *Manager {manage.build.WhereMap(mapWhere);return manage}
func (manage *Manager) OrWhereMap(mapWhere map[string] interface{}) *Manager {manage.build.OrWhereMap(mapWhere);return manage}
func (manage *Manager) WhereIn(field string, listValue [] interface{}) *Manager {manage.build.WhereIn(field, listValue);return manage}
func (manage *Manager) OrWhereIn(field string, listValue [] interface{}) *Manager {manage.build.OrWhereIn(field, listValue);return manage}
func (manage *Manager) WhereNotIn(field string, listValue [] interface{}) *Manager {manage.build.WhereNotIn(field, listValue);return manage}
func (manage *Manager) OrWhereNotIn(field string, listValue [] interface{}) *Manager {manage.build.OrWhereNotIn(field, listValue);return manage}
func (manage *Manager) WhereBetween(field string, interval [] interface{}) *Manager {manage.build.WhereBetween(field, interval);return manage}
func (manage *Manager) OrWhereBetween(field string, interval [] interface{}) *Manager {manage.build.OrWhereBetween(field, interval);return manage}
func (manage *Manager) WhereNotBetween(field string, interval [] interface{}) *Manager {manage.build.WhereNotBetween(field, interval);return manage}
func (manage *Manager) OrWhereNotBetween(field string, interval [] interface{}) *Manager {manage.build.OrWhereNotBetween(field, interval);return manage}
func (manage *Manager) WhereNull(field string) *Manager {manage.build.WhereNull(field);return manage}
func (manage *Manager) OrWhereNull(field string) *Manager {manage.build.OrWhereNull(field);return manage}
func (manage *Manager) WhereNotNull(field string) *Manager {manage.build.WhereNotNull(field);return manage}
func (manage *Manager) OrWhereNotNull(field string) *Manager {manage.build.OrWhereNotNull(field);return manage}
func (manage *Manager) WhereDate(field, date string) *Manager {manage.build.WhereDate(field, date);return manage}
func (manage *Manager) OrWhereDate(field, date string) *Manager {manage.build.OrWhereDate(field, date);return manage}
func (manage *Manager) WhereMonth(field, month string) *Manager {manage.build.WhereMonth(field, month);return manage}
func (manage *Manager) OrWhereMonth(field, month string) *Manager {manage.build.OrWhereMonth(field, month);return manage}
func (manage *Manager) WhereDay(field, day string) *Manager {manage.build.WhereDay(field, day);return manage}
func (manage *Manager) OrWhereDay(field, day string) *Manager {manage.build.OrWhereDay(field, day);return manage}
func (manage *Manager) WhereYear(field, year string) *Manager {manage.build.WhereYear(field, year);return manage}
func (manage *Manager) OrWhereYear(field, year string) *Manager {manage.build.OrWhereYear(field, year);return manage}
func (manage *Manager) WhereTime(field, condition, timestamp string) *Manager {manage.build.WhereTime(field, condition, timestamp);return manage}
func (manage *Manager) OrWhereTime(field, condition, timestamp string) *Manager {manage.build.OrWhereTime(field, condition, timestamp);return manage}
func (manage *Manager) WhereFunc(callback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manager {manage.build.WhereFunc(callback);return manage}
func (manage *Manager) OrWhereFunc(callback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manager {manage.build.OrWhereFunc(callback);return manage}
func (manage *Manager) WhereExists(callback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manager {manage.build.WhereExists(callback);return manage}
func (manage *Manager) OrWhereExists(callback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manager {manage.build.OrWhereExists(callback);return manage}
func (manage *Manager) WhereNotExists(field string, callback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manager {manage.build.WhereNotExists(field, callback);return manage}
func (manage *Manager) OrWhereNotExists(field string, callback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manager {manage.build.OrWhereNotExists(field, callback);return manage}
func (manage *Manager) WhereRaw(sql string) *Manager {manage.build.WhereRaw(sql);return manage}
func (manage *Manager) OrWhereRaw(sql string) *Manager {manage.build.OrWhereRaw(sql);return manage}
func (manage *Manager) When(boolean bool, callback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manager {manage.build.When(boolean, callback);return manage}
func (manage *Manager) OrWhen(boolean bool, callback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manager {manage.build.OrWhen(boolean, callback);return manage}
func (manage *Manager) WhenElse(boolean bool, trueCallback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild, falseCallback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manager {manage.build.WhenElse(boolean, trueCallback, falseCallback);return manage}
func (manage *Manager) OrWhenElse(boolean bool, trueCallback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild, falseCallback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manager {manage.build.OrWhenElse(boolean, trueCallback, falseCallback);return manage}
func (manage *Manager) OrderBy(args... string) *Manager {manage.build.OrderBy(args...);return manage}
func (manage *Manager) OrderByRaw(sql string) *Manager {manage.build.OrderByRaw(sql);return manage}
func (manage *Manager) GroupBy(args... string) *Manager {manage.build.GroupBy(args...);return manage}
func (manage *Manager) GroupByRaw(sql string) *Manager {manage.build.GroupByRaw(sql);return manage}
func (manage *Manager) Having(args... string) *Manager {manage.build.Having(args...);return manage}
func (manage *Manager) HavingRaw(sql string) *Manager {manage.build.HavingRaw(sql);return manage}
func (manage *Manager) Offset(num int) *Manager {manage.build.Offset(num);return manage}
func (manage *Manager) Limit(num int) *Manager {manage.build.Limit(num);return manage}