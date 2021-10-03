package ws

import (
	"sync"

	"github.com/sirupsen/logrus"
)

type syncMap struct {
	sync.Map
	count int
}

// 连接管理
type Manager struct {
	clients    *syncMap     // 全部的连接
	groups     *syncMap     // 分组的连接
	systems    *syncMap     // 系统的连接
	Connect    chan *Client // 连接处理
	DisConnect chan *Client // 断开处理
}

func newManager() *Manager {
	return &Manager{
		clients:    new(syncMap),
		groups:     new(syncMap),
		systems:    new(syncMap),
		Connect:    make(chan *Client, 100),
		DisConnect: make(chan *Client, 100),
	}
}

// 管道处理程序
func (manager *Manager) start() *Manager {
	go func() {
		for {
			select {
			case client := <-manager.Connect:
				manager.eventConnect(client)
			case client := <-manager.DisConnect:
				manager.eventDisconnect(client)
			}
		}
	}()
	return manager
}

// 建立连接事件
func (manager *Manager) eventConnect(client *Client) {
	logrus.Infof("WS用户连接, uid: %d", client.Uid)
	manager.addClient(client)
}

// 断开连接时间
func (manager *Manager) eventDisconnect(client *Client) {

	client.engine.DelWsUserAddr(client.Uid)
	logrus.Infof("WS用户断开, uid: %d", client.Uid)
	//关闭连接
	client.Socket.Close()
	manager.DelClient(client)
	//标记销毁
	client.IsDeleted = true
	client = nil
}
