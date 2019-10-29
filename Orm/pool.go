package Orm

import (
	"blog/Config/database"
	_ "github.com/go-sql-driver/mysql"
	"blog/Library"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"
)

type Client struct {
	MysqlClient *sql.DB			// Mysql 连接
	expire time.Time			// 超时时间 	0：不超时
}

type MysqlConfig struct {
	InitCap int
	MaxCap int
	DBName string
	Dsn string
	clientTimeOut int
	checkClientAliveInterval int
	keepClientTime int
	host string
	port int
	username string
	passwd string
}

type Pool struct {
	wait sync.RWMutex
	useMap map[*sql.DB] *Client		// 使用中链接
	Config MysqlConfig				// 配置信息
	Clients chan *Client			// 空闲链接池
	ClientNum int					// 已生成链接数
}

var mysqlExpire int
var myPool *Pool
var waitGroup sync.WaitGroup


/*
 * 初始化连接池
 */
func (this *Pool) InitClient() {
	waitGroup.Add(this.Config.InitCap)
	for i := 0; i < this.Config.InitCap; i++ {
		go func() {
			this.Clients <- this.CreateClient()
			waitGroup.Done()
		}()
	}
	waitGroup.Wait()
}

/*
 * 连接Mysql
 */
func (this *Pool) clientMysql() (*sql.DB, error) {
	this.wait.RLock()
	defer this.wait.RUnlock()

	if this.ClientNum == this.Config.MaxCap {
		return nil, errors.New("Max")
	}

	db, err := sql.Open("mysql", this.Config.Dsn)
	Library.CheckErr(err)

	return db, nil
}

func (this *Pool) CreateClient() *Client {
	db, err := this.clientMysql()
	Library.CheckErr(err)
	this.ClientNum++
	time_unit, _ := time.ParseDuration("1h")

	return &Client{
		MysqlClient:	db,
		expire:			time.Now().Add(this.Config.keepClientTime * time_unit),
	}
}

/*
 * 定期检查失效链接
 */
func (this *Pool) checkInvalidClient() {
	this.wait.Lock()
	defer this.wait.Unlock()

	tmp_client := make(chan *Client, len(this.Clients))
	for client := range this.Clients {
		if this.Len() >= this.Config.MaxCap {
			_=this.Close(client)
			break
		}

		if time.Now().After(client.expire) {
			_=this.Close(client)
			break
		}

		if err := client.MysqlClient.Ping(); err != nil {
			_=this.Close(client)
			break
		}

		tmp_client <- client
	}

	for client := range tmp_client {
		this.Clients <- client
	}
}

func (this *Pool) Len() int {
	return len(this.Clients)
}

func (this *Pool) Close(client *Client) error {
	err := client.MysqlClient.Close()
	this.ClientNum--
	if err != nil {
		return err
	}

	return nil
}

func InitConfig(config database.Mysql) MysqlConfig {
	return MysqlConfig{
		InitCap:		config.InitCap,
		MaxCap:			config.MaxCap,
		clientTimeOut:  config.ClientTimeOut,
		checkClientAliveInterval:  config.CheckClientAliveInterval,
		keepClientTime: config.KeepClientTime,
		host:			config.Host,
		port:			config.Port,
		username:		config.Username,
		passwd:			config.Passwd,
		DBName:			config.DBName,
		Dsn:			fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", config.Username, config.Passwd, config.Host, config.Port, config.DBName),
	}
}