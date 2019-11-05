package Orm

import (
	"fmt"
	"testing"
)

type Goods struct {
	Id int
	GoodsName string
	GoodsStock int
	GoodsStatus int
}

func TestGet(t *testing.T) {
	result, err := DB().Table("goods").Select("id", "goods_name").Where("goods_stock", ">", "1").Get()
	if err != nil {
		t.Errorf("Get error: %s", err)
	}
	fmt.Println("Get: ", result)
}

func TestFirst(t *testing.T) {
	result, err := DB().Table("goods").Select("id", "goods_name").Where("goods_stock", ">", "1").First()
	if err != nil {
		t.Errorf("First error: %s", err)
	}
	fmt.Println("First: ", result)
}

func TestValue(t *testing.T) {
	result, err := DB().Table("goods").Where("goods_stock", ">", "1").Value("goods_name")
	if err != nil {
		t.Errorf("Value error: %s", err)
	}
	fmt.Println("Value: ", result)
}

func TestPluckArray(t *testing.T) {
	result, err := DB().Table("goods").Where("goods_stock", ">", "1").PluckArray("goods_name")
	if err != nil {
		t.Errorf("PluckArray error: %s", err)
	}
	fmt.Println("PluckArray: ", result)
}

func TestPluckMap(t *testing.T) {
	result, err := DB().Table("goods").Where("goods_stock", ">", "1").PluckMap("id", "goods_name")
	if err != nil {
		t.Errorf("PluckMap error: %s", err)
	}
	fmt.Println("PluckMap: ", result)
}

func TestCount(t *testing.T) {
	result, err := DB().Table("goods").Where("goods_stock", ">", "1").Count()
	if err != nil {
		t.Errorf("Count error: %s", err)
	}
	fmt.Println("Count: ", result)
}

func TestMax(t *testing.T) {
	result, err := DB().Table("goods").Where("goods_stock", ">", "1").Max("id")
	if err != nil {
		t.Errorf("Max error: %s", err)
	}
	fmt.Println("Max: ", result)
}

func TestSum(t *testing.T) {
	result, err := DB().Table("goods").Where("goods_stock", ">", "1").Sum("id")
	if err != nil {
		t.Errorf("Sum error: %s", err)
	}
	fmt.Println("Sum: ", result)
}

func TestInsert(t *testing.T) {
	var arr = make(map[string]interface{})
	arr["goods_name"] = "小小当家"
	arr["goods_status"] = 1
	arr["goods_stock"] = 250
	s := DB().Table("goods").InsertToSql(arr)
	fmt.Println("Insert sql : ", s)
	affected_num, err := DB().Table("goods").Insert(arr)
	if err != nil {
		t.Errorf("Insert error: %s", err)
	}
	fmt.Println("Insert: ", affected_num)
}

func TestLastInsertId(t *testing.T) {
	var arr = make(map[string]interface{})
	arr["goods_name"] = "中当家"
	arr["goods_status"] = 1
	arr["goods_stock"] = 250
	inert_id, err := DB().Table("goods").LastInsertId(arr)
	if err != nil {
		t.Errorf("LastInsertId error: %s", err)
	}
	fmt.Println("LastInsertId: ", inert_id)
}

func TestMultiInsert(t *testing.T) {
	array := make([]map[string]interface{}, 0)
	var arr = make(map[string]interface{})
	arr["goods_name"] = "小当家"
	arr["goods_status"] = 1
	arr["goods_stock"] = 250
	array = append(array, arr)
	var arr2 = make(map[string]interface{})
	arr2["goods_name"] = "二当家"
	arr2["goods_status"] = 1
	arr2["goods_stock"] = 250
	array = append(array, arr2)
	var arr3 = make(map[string]interface{})
	arr3["goods_name"] = "三当家"
	arr3["goods_status"] = 1
	arr3["goods_stock"] = 250
	array = append(array, arr3)
	var arr4 = make(map[string]interface{})
	arr4["goods_name"] = "大当家"
	arr4["goods_status"] = 1
	arr4["goods_stock"] = 250
	array = append(array, arr4)
	affected_num, err := DB().Table("goods").MultiInsert(array)
	if err != nil {
		t.Errorf("MultiInsert error: %s", err)
	}
	fmt.Println("MultiInsert: ", affected_num)
}

func TestUpdate(t *testing.T) {
	var arr = make(map[string]interface{})
	arr["goods_name"] = "四当家"
	arr["goods_status"] = 2
	arr["goods_stock"] = 350
	s := DB().Table("goods").Where("goods_name", "三当家").UpdateToSql(arr)
	fmt.Println("Update sql : ", s)
	affected_num, err := DB().Table("goods").Where("goods_name", "三当家").Update(arr)
	if err != nil {
		t.Errorf("Update error: %s", err)
	}
	fmt.Println("Update: ", affected_num)
}

func TestDelete(t *testing.T) {
	s := DB().Table("goods").Where("goods_name", "二当家").DeleteToSql()
	fmt.Println("Delete sql : ", s)
	affected_num, err := DB().Table("goods").Where("goods_name", "二当家").Delete()
	if err != nil {
		t.Errorf("Delete error: %s", err)
	}
	fmt.Println("Delete: ", affected_num)
}