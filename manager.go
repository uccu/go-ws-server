package ws

import (
	"context"
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
	unions     *syncMap     // 标记的连接
	groups     *syncMap     // 分组的连接
	systems    *syncMap     // 系统的连接
	Connect    chan *Client // 连接处理
	DisConnect chan *Client // 断开处理
	cancel     context.CancelFunc
}

func newManager() *Manager {
	return &Manager{
		clients:    new(syncMap),
		unions:     new(syncMap),
		groups:     new(syncMap),
		systems:    new(syncMap),
		Connect:    make(chan *Client, 100),
		DisConnect: make(chan *Client, 100),
	}
}

// 管道处理程序
func (manager *Manager) start() *Manager {
	ctx, cancel := context.WithCancel(context.TODO())
	manager.cancel = cancel
	go func() {
		for {
			select {
			case client := <-manager.Connect:
				manager.eventConnect(client)
			case client := <-manager.DisConnect:
				manager.eventDisconnect(client)
			case <-ctx.Done():
				return
			}
		}
	}()
	return manager
}

func (manager *Manager) close() {
	if manager.cancel != nil {
		manager.cancel()
	}
}

// 建立连接事件
func (manager *Manager) eventConnect(client *Client) {
	logrus.Infof("WS用户连接, clientId: %d", client.ClientId)
	manager.addClient(client)
}

// 断开连接时间
func (manager *Manager) eventDisconnect(client *Client) {
	logrus.Infof("WS用户断开, clientId: %d", client.ClientId)
	client.Socket.Close()
	manager.DelClient(client)
	client.IsDeleted = true
	client = nil
}
