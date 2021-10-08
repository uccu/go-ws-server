package ws

import "errors"

// 删除客户端
func (manager *Manager) DelClient(client *Client) {
	manager.delClientByClientId(client.ClientId)
	if client.UnionId != "" {
		manager.delClientUnion(client.UnionId)
	}
	for groupKey := range client.GroupList {
		manager.delGroupClient(groupKey, client.ClientId)
	}
	for systemKey := range client.SystemList {
		manager.delSystemClient(systemKey, client.ClientId)
	}
}

// 添加客户端
func (manager *Manager) addClient(client *Client) {
	_, loaded := manager.clients.LoadOrStore(client.ClientId, client)
	if !loaded {
		manager.clients.count++
	}
}

// 获取所有的客户端
func (manager *Manager) GetClientList(f func(clientId string, client *Client) bool) {
	manager.clients.Range(func(key, value interface{}) bool {
		return f(key.(string), value.(*Client))
	})
}

// 获取所有的客户端
func (manager *Manager) GetClientCount() int {
	return manager.clients.count
}

// 通过clientId删除client
func (manager *Manager) delClientByClientId(clientId string) {
	_, loaded := manager.clients.LoadAndDelete(clientId)
	if loaded {
		manager.clients.count--
	}
}

// 通过clientId获取client
func (manager *Manager) GetClientByClientId(clientId string) (*Client, error) {
	if client, ok := manager.clients.Load(clientId); !ok {
		return nil, errors.New("客户端不存在")
	} else {
		return client.(*Client), nil
	}
}

func (manager *Manager) SendAllMessage(ack string, data interface{}) {
	manager.GetClientList(func(clientId string, client *Client) bool {
		client.SendMessage(ack, data)
		return true
	})
}
