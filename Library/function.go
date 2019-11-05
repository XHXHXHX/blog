package Library

import (
	"strconv"
)

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}


/*
 * 用 char 连接符生成 Mysql 日期格式化字符串
 * %Y%m%d
 */
func SetMysqlDateFormatByChar(char string) string {
	return "%Y" + char + "%m" + char + "%d"
}

/*
 * 转化成字符串
 */
func TransferString(s interface{}) string {
	switch r := s.(type) {
	case string:
		return r
	case int:
		return strconv.Itoa(r)
	default:
		panic("transfer string error")
	}
}

func AddSingleSymbol(s string) string {
	return "'" + string(s) + "'"
}
