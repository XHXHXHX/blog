package clientPool

import (
	"fmt"
	"sync"
	"testing"
	"container/list"
	"blog/Config/database"
	"database/sql"
)

var wait sync.WaitGroup

func TestGetClient(t *testing.T) {
	config, err := database.GetMysqlConfig()
	if err != nil {
		t.Errorf("config error %s", err)
	}

	if myPool.Len() != config.InitCap {
		t.Errorf("init client num error, config init_cap : %d, client num : %d", config.InitCap, myPool.Len())
	}

	client, err := GetClient()
	if err != nil {
		fmt.Println(client)
	}

	if myPool.Len() != config.InitCap - 1 {
		t.Errorf("client num error after get a client, expect num : %d, actual num : %d", config.InitCap - 1, myPool.Len())
	}

	err = CloseClient(client)
	if err != nil {
		t.Errorf("close client error : %s", err)
	}

	if myPool.Len() != config.InitCap {
		t.Errorf("client num error after close the client, expect num : %d, actual num : %d", config.InitCap, myPool.Len())
	}
}


func TestMore(t *testing.T) {
	config, err := database.GetMysqlConfig()
	if err != nil {
		t.Errorf("config error %s", err)
	}

	client1, err := GetClient()
	if err != nil {
		t.Errorf("config error %s", err)
	}

	client2, err := GetClient()
	if err != nil {
		t.Errorf("config error %s", err)
	}
	client3, err := GetClient()
	if err != nil {
		t.Errorf("config error %s", err)
	}
	client4, err := GetClient()
	if err != nil {
		t.Errorf("config error %s", err)
	}

	if myPool.ClientNum != config.InitCap + 1 {
		t.Errorf("client num error after get more client with init_cap, expect num %d, acutal num %d", config.InitCap + 1, myPool.ClientNum)
	}

	_ = CloseClient(client1)
	_ = CloseClient(client2)
	_ = CloseClient(client3)
	_ = CloseClient(client4)

}

func TestOverTime(t *testing.T) {
	config, err := database.GetMysqlConfig()
	if err != nil {
		t.Errorf("config error %s", err)
	}

	var clientList [] *sql.DB
	wait.Add(config.MaxCap)
	for i := 0; i < config.MaxCap; i++ {
		go func() {
			client, err := GetClient()
			if err != nil {
				t.Errorf("get client error %s at %d", err, i)
			}
			clientList = append(clientList, client)
			wait.Done()
		}()
	}
	wait.Wait()

	if client, err := GetClient(); err == nil {
		err := CloseClient(client)
		if err == nil {
			t.Errorf("close client error %s", err)
		}
		t.Errorf("should over time but no, error: %s", err)
	}

	wait.Add(config.MaxCap)
	for i := range clientList {
		go func() {
			err := CloseClient(clientList[i])
			if err != nil {
				t.Errorf("close client error %s at %d", err, i)
			}
			wait.Done()
		}()
	}
	wait.Wait()
}

func TestWait(t *testing.T) {
	defer OnClose()

	config, err := database.GetMysqlConfig()
	if err != nil {
		t.Errorf("config error %s", err)
	}

	var clientList = list.New()
	wait.Add(config.MaxCap)
	for i := 0; i < config.MaxCap; i++ {
		go func() {
			client, err := GetClient()
			if err != nil {
				t.Errorf("get client error %s at %d", err, i)
				return
			}
			clientList.PushBack(client)
			wait.Done()
		}()
	}
	wait.Wait()

	wait.Add(2)
	var client *sql.DB
	go func(client *sql.DB) {
		client, err := GetClient()
		if err != nil {
			t.Errorf("wait error : %s", err)
		}
		clientList.PushBack(client)
		wait.Done()
	}(client)

	go func() {
		ele := clientList.Front()
		if client1, ok := ele.Value.(*sql.DB); ok {
			_ = CloseClient(client1)
			clientList.Remove(ele)
		} else {
			t.Errorf("clientList element type error : %v", ele.Value)
		}
		wait.Done()
	}()

	wait.Wait()

	wait.Add(clientList.Len())
	for ele := clientList.Front(); ele != nil; ele = ele.Next() {
		go func(ele *list.Element) {
			if client, ok := ele.Value.(*sql.DB); ok {
				_ = CloseClient(client)
			}
			wait.Done()
		}(ele)
	}
	wait.Wait()
}