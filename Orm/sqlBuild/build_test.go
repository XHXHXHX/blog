package sqlBuild

import (
	"fmt"
	"strings"
	"testing"
)

const (
	TextHeaderNum = 20
)

func TestDB(t *testing.T) {
	result := DB().Table("class", "c").Where("a", "1").OrWhere("b", 2).Get()
	myPrintln("Where", result.ShowSql)
}

func TestSqlBuild_WhereIn(t *testing.T) {
	var arr []interface{}
	arr = append(arr, 1)
	arr = append(arr, "2")
	arr = append(arr, 4)
	arr = append(arr, 5)

	result := DB().Table("class").Where("a", "1").OrWhereNotIn("type", arr).Get()
	myPrintln("OrWhereNotIn", result.ShowSql)
}

func TestSqlBuild_WhereNull(t *testing.T) {
	result := DB().Table("goods").WhereNull("exrat_data").Get()
	myPrintln("WhereNull", result.ShowSql)
}

func TestSqlBuild_WhereBetween(t *testing.T) {
	var arr []interface{}
	arr = append(arr, 123)
	arr = append(arr, "456")

	result := DB().Table("class").WhereNotBetween("type", arr).Get()
	myPrintln("WhereNotBetween", result.ShowSql)
}

func TestSqlBuild_WhereArray(t *testing.T) {
	var arr = [][] string {
		{
			"id", "2",
		},{
			"is_del", "!=", "0",
		},
	}

	result := DB().Table("class").WhereArray(arr).Get()
	myPrintln("WhereArray", result.ShowSql)
}

func TestSqlBuild_WhereMap(t *testing.T) {
	myMap := make(map[string]interface{})
	myMap["id"] = 2
	myMap["name"] = 1
	myMap["age"] = "23"

	result := DB().Table("class").WhereMap(myMap).Get()
	myPrintln("WhereMap", result.ShowSql)
}

func TestSqlBuild_WhereDate(t *testing.T) {
	result := DB().Table("goods").WhereDate("add_time", "2019-09-21").Get()
	myPrintln("WhereDate", result.ShowSql)
}

func TestSqlBuild_WhereMonth(t *testing.T) {
	result := DB().Table("goods").WhereMonth("add_time", "10").Get()
	myPrintln("WhereMonth", result.ShowSql)
}

func TestSqlBuild_WhereDay(t *testing.T) {
	result := DB().Table("goods").WhereDay("add_time", "31").Get()
	myPrintln("WhereDay", result.ShowSql)
}

func TestSqlBuild_WhereYear(t *testing.T) {
	result := DB().Table("goods").WhereYear("add_time", "2020").Get()
	myPrintln("WhereYear", result.ShowSql)
}

func TestSqlBuild_WhereTime(t *testing.T) {
	result := DB().Table("goods").WhereTime("add_time", "<", "13:20:11").Get()
	myPrintln("WhereTime", result.ShowSql)
}

func TestSqlBuild_WhereFunc(t *testing.T) {
	var aaa = 1
	result := DB().Table("goods", "g").WhereFunc(func(query *SqlBuild) *SqlBuild {
		aaa += 2
		var arr []interface{}
		arr = append(arr, aaa)
		arr = append(arr, "99")
		return query.Where("goods_stock", 1).WhereBetween("goods_stock", arr)
	}).Where("id", 23).Get()

	myPrintln("WhereFunc", result.ShowSql)
}

func TestSqlBuild_WhereRaw(t *testing.T) {
	result := DB().Table("goods").WhereRaw("goods_name = '花卷'").Get()
	myPrintln("WhereRaw", result.ShowSql)
}

func TestSqlBuild_GroupBy(t *testing.T) {
	result := DB().Table("goods").GroupBy("goods_type").Get()
	myPrintln("GroupBy", result.ShowSql)
}

func TestSqlBuild_OrderBy(t *testing.T) {
	result := DB().Table("goods").OrderBy("id desc", "add_time asc").Get()
	myPrintln("OrderBy", result.ShowSql)
}

func TestSqlBuild_Join(t *testing.T) {
	result := DB().Table("goods", "g").Join("goods_score", "goods_id", "=", "goods_id").Get()
	myPrintln("Join", result.ShowSql)
}

func TestSqlBuild_LeftJoin(t *testing.T) {
	result := DB().Table("goods", "g").LeftJoin("goods_score as s", "goods_id", "=", "goods_id").Get()
	myPrintln("LeftJoin", result.ShowSql)
}

func TestSqlBuild_LeftJoinFunc(t *testing.T) {
	result := DB().Table("goods", "g").LeftJoinFunc("goods_score as s", func(build *SqlBuild) *SqlBuild {
		return build.On("goods_id", "=", "goods_id").Where("is_del", 0).WhereNull("g.add_time")
	}).Where("g.goods_stock", ">", 0).Get()
	myPrintln("LeftJoinFunc", result.ShowSql)
}

func myPrintln(s, sql string) {
	l := len(s)
	fmt.Println(s, strings.Repeat(" ", TextHeaderNum - l), sql)
}