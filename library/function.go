package library

import (
	"crypto/md5"
	"fmt"
	"path"
	"runtime"
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

func GetCurrentPath() string {
	_, filename, _, _ := runtime.Caller(1)

	return path.Dir(filename)
}

func ArrayToMayFromKey(arr []string) map[string]string {
	arr_map := make(map[string]string)
	for _, item := range arr {
		key := string(item)
		arr_map[key] = key
	}
	return arr_map
}

func Md5Encryption(s string) string {
	data := []byte(s)
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has) //将[]byte转成16进制
}