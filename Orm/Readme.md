## 仿 Laravel DB Face

```$golang
    import Orm
    
    func Main() {
    	result, err := DB().Table("goods").Select("id", "goods_name").Where("goods_stock", ">", "1").Get()
    	if err != nil {
    		t.Errorf("Get error: %s", err)
    	}
    	fmt.Println("Get: ", result)
    }
```

- `DB()`方法
    
    返回 Manager 对象
    
- `Table(args... string)`方法  两个参数
    
    接受两个参数，第一个为表名，第二个为别名，第二个参数可不传
    Eg: `Table("goods")` || `Table("goods", "g")`
    
-  `Select(args... string)` 方法 若干参数

    Eg: `Select("id", "goods_name")`
    
- `SelectRaw(sql string)`
- Join
    - `Join(table, thatRelationField, relationCondition, thisRelationField string)`
    - `LeftJoin(table, thatRelationField, relationCondition, thisRelationField string)`
    - `RightJoin(table, thatRelationField, relationCondition, thisRelationField string)`
    - `InnerJoin(table, thatRelationField, relationCondition, thisRelationField string)`
    - `On(thatRelationField, relationCondition, thisRelationField string)`
    - `LeftJoinFunc(table string, callback func(build *sqlBuild.SqlBuild)`
    - `RightJoinFunc(table string, callback func(build *sqlBuild.SqlBuild)`
    - `InnerJoinFunc(table string, callback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild)`
    - `JoinFunc(table string, callback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild)`
    
- Where

    - `Where(args... interface{})` 方法  最多三个参数
    
        Eg: `Where("id", 1)`    || `Where("id", ">", 1)`
        
    - `WhereArray(arrayWhere [][] string)`
        
        Eg: ```[
            ["id", 1],
            ["goods_name", "like", "test%"]
        ]```
        
    - `WhereMap(mapWhere map[string] interface{})`
    
        Eg: ```{
            "id": 1,
            "goods_stock": 0
        }```
        
    - `WhereIn(field string, listValue [] interface{})`
    - `WhereNotIn(field string, listValue [] interface{})`
    - `WhereBetween(field string, interval [] interface{})`
    - `WhereNotBetween(field string, interval [] interface{})`
    - `WhereNull(field string)`
    - `WhereNotNull(field string)`
    - `WhereDate(field, date string)`
    - `WhereMonth(field, month string)`
    - `WhereDay(field, day string)`
    - `WhereYear(field, year string)`
    - `WhereTime(field, condition, timestamp string)`
    - `WhereFunc(callback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild)`
    - `WhereRaw(sql string)`
- Other
    - `When(boolean bool, callback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild)`
    - `WhenElse(boolean bool, trueCallback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild, falseCallback func(build *sqlBuild.SqlBuild) *sqlBuild.SqlBuild) *Manager`
    - `OrderBy(args... string)`
    - `OrderByRaw(sql string)`
    - `GroupBy(args... string)`
    - `GroupByRaw(sql string)`
    - `Having(args... string)`
    - `HavingRaw(sql string)`
    - `Offset(num int)`
    - `Limit(num int)`
- Query
    - `Get(args... string)`     可代替Select
    - `Value(field string)`     返回单个字段值
    - `First()`                 返回第一行数据
    - `PluckArray(field string)`返回单独一列字段
    - `PluckMap(field, value string)`   返回键值对数据
    - `Count()`
    - `Max(field string)`
    - `Sum(field string)`
- Exec
    - `Insert(data map[string]interface{})`
    - `MultiInsert(data []map[string]interface{})`
    - `LastInsertId(data map[string]interface{})`
    - `Update(data map[string]interface{})`
    - `Delete() (int, error)`
    - 上述方法后加 `ToSql` 返回执行sql
    
- 事务
    - `DbBegin() *Manager` 返回事务句柄
    - `DbCommit()`
    - `DbRollBack()`