package redis

import (
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"encoding/json"
)

type Redis struct {
	Host string
	Port int
	Timeout int
	Passwd string
}

func GetRedisConfig() (Redis, error) {
	configPath := getCurrentPath()
	//configPath := "../"
	configName := "db.json"

	var RedisConfig Redis

	configFile, err := os.Open(configPath + "/" + configName)
	if err != nil {
		return RedisConfig, err
	}

	bytes, _ := ioutil.ReadAll(configFile)
	err = json.Unmarshal(bytes, &RedisConfig)

	if err != nil {
		return RedisConfig, err
	}

	return RedisConfig, nil
}

func getCurrentPath() string {
	_, filename, _, _ := runtime.Caller(1)

	return path.Dir(filename)
}