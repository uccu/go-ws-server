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

type EventFunc func(*Manager, *Client) bool

// 连接管理
type Manager struct {
	clients    *syncMap     // 全部的连接
	unions     *syncMap     // 标记的连接
	groups     *syncMap     // 分组的连接
	systems    *syncMap     // 系统的连接
	connect    chan *Client // 连接处理
	DisConnect chan *Client // 断开处理
	cancel     context.CancelFunc

	connectFunc    []EventFunc
	disconnectFunc []EventFunc
}

func newManager() *Manager {
	return &Manager{
		clients:        new(syncMap),
		unions:         new(syncMap),
		groups:         new(syncMap),
		systems:        new(syncMap),
		connect:        make(chan *Client, 100),
		DisConnect:     make(chan *Client, 100),
		connectFunc:    make([]EventFunc, 0),
		disconnectFunc: make([]EventFunc, 0),
	}
}

func (manager *Manager) WithConnect(f EventFunc) *Manager {
	manager.connectFunc = append(manager.connectFunc, f)
	return manager
}

func (manager *Manager) WithDisConnect(f EventFunc) *Manager {
	manager.disconnectFunc = append(manager.disconnectFunc, f)
	return manager
}

// 管道处理程序
func (manager *Manager) start() *Manager {
	ctx, cancel := context.WithCancel(context.TODO())
	manager.cancel = cancel
	go func() {
		for {
			select {
			case client := <-manager.connect:
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

func (manager *Manager) Close() {
	if manager.cancel != nil {
		manager.cancel()
	}
}

// 建立连接事件
func (manager *Manager) eventConnect(client *Client) {
	for _, f := range manager.connectFunc {
		if !f(manager, client) {
			return
		}
	}

	logrus.Infof("WS用户连接, clientId: %s", client.ClientId)
	manager.addClient(client)
}

// 断开连接时间
func (manager *Manager) eventDisconnect(client *Client) {
	for _, f := range manager.disconnectFunc {
		if !f(manager, client) {
			return
		}
	}
	logrus.Infof("WS用户断开, clientId: %s", client.ClientId)
	client.Socket.Close()
	manager.DelClient(client)
	client.IsDeleted = true
	client = nil
}
