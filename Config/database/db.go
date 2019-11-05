package database

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"runtime"
)

type Mysql struct {
	InitCap int
	MaxCap int
	ClientTimeOut int
	CheckClientAliveInterval int
	KeepClientTime int
	Host string
	Port int
	Username string
	Passwd string
	DBName string
	Prifex string
}

func GetMysqlConfig() (Mysql, error) {
	configPath := getCurrentPath()
	//configPath := "../"
	configName := "db.json"

	var MysqlConfig Mysql

	configFile, err := os.Open(configPath + "/" + configName)
	if err != nil {
		return MysqlConfig, err
	}

	bytes, _ := ioutil.ReadAll(configFile)
	err = json.Unmarshal(bytes, &MysqlConfig)

	if err != nil {
		return MysqlConfig, err
	}

	return MysqlConfig, nil
}

func getCurrentPath() string {
	_, filename, _, _ := runtime.Caller(1)

	return path.Dir(filename)
}