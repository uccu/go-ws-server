package ws

import "sync"

// 添加到系统
func (manager *Manager) AddSystemClient(systemKey string, client *Client) {
	if _, ok := client.SystemList[systemKey]; ok {
		return
	}
	manager.addSystemClient(systemKey, client.ClientId)
	client.SystemList[systemKey] = true
}

// 删除系统里的客户端
func (manager *Manager) DelSystemClient(systemKey string, client *Client) {
	if _, ok := client.SystemList[systemKey]; !ok {
		return
	}
	manager.delSystemClient(systemKey, client.ClientId)
	delete(client.SystemList, systemKey)
}

// 添加到系统
func (manager *Manager) addSystemClient(systemKey string, clientId string) {
	clientIdMap, _ := manager.systems.LoadOrStore(systemKey, new(syncMap))
	sm := clientIdMap.(*syncMap)
	_, loaded := sm.LoadOrStore(clientId, 1)
	if !loaded {
		sm.count++
	}
}

// 删除系统里的客户端
func (manager *Manager) delSystemClient(systemKey string, clientId string) {
	if clientIdMap, ok := manager.systems.Load(systemKey); ok {
		sm := clientIdMap.(*syncMap)
		_, loaded := sm.LoadAndDelete(clientId)
		if loaded {
			sm.count--
		}
	}
}

// 获取系统所有的客户端数量
func (manager *Manager) GetSystemClientCount(systemKey string) int {
	if clientIdMap, ok := manager.systems.Load(systemKey); ok {
		return clientIdMap.(*syncMap).count
	}
	return 0
}

// 获取系统的成员
func (manager *Manager) GetSystemClientList(systemKey string, f func(clientId string) bool) {
	if clientIdMap, ok := manager.systems.Load(systemKey); ok {
		clientIdMap.(*sync.Map).Range(func(key, value interface{}) bool {
			return f(key.(string))
		})
	}
}

// 发送系统消息
func (manager *Manager) SendSystemMessage(systemKey string, ack string, data interface{}) {
	manager.GetSystemClientList(systemKey, func(clientId string) bool {
		client, _ := manager.GetClientByClientId(clientId)
		client.SendMessage(ack, data)
		return true
	})
}
