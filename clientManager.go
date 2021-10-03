package ws

import "errors"

// 删除客户端
func (manager *Manager) DelClient(client *Client) {
	manager.delclientsByUid(client.Uid)
	for groupKey := range client.GroupList {
		manager.delGroupClient(groupKey, client.Uid)
	}
	for systemKey := range client.SystemList {
		manager.delSystemClient(systemKey, client.Uid)
	}
}

// 添加客户端
func (manager *Manager) addClient(client *Client) {
	_, loaded := manager.clients.LoadOrStore(client.Uid, client)
	if !loaded {
		manager.clients.count++
	}
}

// 获取所有的客户端
func (manager *Manager) GetClientList(f func(uid UID, client *Client) bool) {
	manager.clients.Range(func(key, value interface{}) bool {
		return f(key.(UID), value.(*Client))
	})
}

// 获取所有的客户端
func (manager *Manager) GetClientCount() int {
	return manager.clients.count
}

// 通过uid删除client
func (manager *Manager) delclientsByUid(uid UID) {
	_, loaded := manager.clients.LoadAndDelete(uid)
	if loaded {
		manager.clients.count--
	}
}

// 通过uid获取client
func (manager *Manager) GetClientByUid(uid UID) (*Client, error) {
	if client, ok := manager.clients.Load(uid); !ok {
		return nil, errors.New("客户端不存在")
	} else {
		return client.(*Client), nil
	}
}

func (manager *Manager) SendAllMessage(ack string, data interface{}) {
	manager.GetClientList(func(uid UID, client *Client) bool {
		client.SendMessage(ack, data)
		return true
	})
}
