package clientPool

import (
	"blog/Config/database"
	"database/sql"
	"errors"
	"time"
)

func init() {
	config, err := database.GetMysqlConfig()
	if err != nil {
		panic(err)
	}
	if config.InitCap < 0 || config.MaxCap < 0 || config.InitCap > config.MaxCap {
		panic("invalid capacity settings")
	}

	mysqlExpire = config.KeepClientTime

	myPool = &Pool{
		useMap: 	make(map[*sql.DB] *Client),
		Config:		InitConfig(config),
		Clients:	make(chan *Client, config.MaxCap),
		ClientNum: 	0,
	}

	myPool.InitClient()

	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(myPool.Config.checkClientAliveInterval))
		defer ticker.Stop()
		for _ = range ticker.C {
			myPool.checkInvalidClient()
		}
	}()
}

/*
 * 获取链接
 */
func GetClient() (*sql.DB, error) {
	// 设置超时时间
	ticker := time.NewTicker(time.Duration(myPool.Config.clientTimeOut) * time.Second)
	for {
		select {
		case client := <- myPool.Clients:
			if time.Now().After(client.expire) {
				_ = myPool.Close(client)
				continue
			}
			myPool.useMap[client.MysqlClient] = client
			return client.MysqlClient, nil
		case <-ticker.C:
			return nil, errors.New("client expire time")
		default:
			// 已有链接数小于最大值生成新链接
			if myPool.ClientNum < myPool.Config.MaxCap {
				client := myPool.CreateClient()
				myPool.useMap[client.MysqlClient] = client
				return client.MysqlClient, nil
			}
		}
	}
}

/*
 * 关闭链接
 * 根据情况 放回链接池 / 关闭连接
 */
func CloseClient(mysqlClient *sql.DB) error {
	client, ok := myPool.useMap[mysqlClient]
	if !ok {
		return errors.New("invalid mysql client")
	}

	myPool.wait.RLock()
	clientLen := myPool.ClientNum
	myPool.wait.RUnlock()
	// 未过期且少于最小链接数，将链接放回链接池
	if time.Now().Before(client.expire) && clientLen < myPool.Config.InitCap {
		myPool.Clients <- client
		return nil
	}

	err := myPool.Close(client)

	if err != nil {
		return err
	}

	return nil
}

func OnClose() [] error {
	// 判断管道是否已经关闭
	_, ok := <- myPool.Clients
	if(!ok) {
		return nil
	}
	close(myPool.Clients)
	var errs [] error
	for client := range myPool.Clients {
		err := myPool.Close(client)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}