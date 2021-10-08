package ws

import "sync"

// 添加到分组
func (manager *Manager) AddGroupClient(groupKey string, client *Client) {
	if _, ok := client.GroupList[groupKey]; ok {
		return
	}
	manager.addGroupClient(groupKey, client.ClientId)
	client.GroupList[groupKey] = true
}

//删除分组里的客户端
func (manager *Manager) DelGroupClient(groupKey string, client *Client) {
	if _, ok := client.GroupList[groupKey]; !ok {
		return
	}
	manager.delGroupClient(groupKey, client.ClientId)
	delete(client.GroupList, groupKey)
}

// 添加到分组
func (manager *Manager) addGroupClient(groupKey string, clientId string) {
	clientIdMap, _ := manager.groups.LoadOrStore(groupKey, new(syncMap))
	sm := clientIdMap.(*syncMap)
	_, loaded := sm.LoadOrStore(clientId, 1)
	if !loaded {
		sm.count++
	}
}

// 删除分组里的客户端
func (manager *Manager) delGroupClient(groupKey string, clientId string) {
	if clientIdMap, ok := manager.groups.Load(groupKey); ok {
		sm := clientIdMap.(*syncMap)
		_, loaded := sm.LoadAndDelete(clientId)
		if loaded {
			sm.count--
		}
	}
}

// 获取分组所有的客户端数量
func (manager *Manager) GetGroupClientCount(groupKey string) int {
	if clientIdMap, ok := manager.groups.Load(groupKey); ok {
		return clientIdMap.(*syncMap).count
	}
	return 0
}

// 获取分组的成员
func (manager *Manager) GetGroupClientList(groupKey string, f func(clientId string) bool) {
	if clientIdMap, ok := manager.groups.Load(groupKey); ok {
		clientIdMap.(*sync.Map).Range(func(key, value interface{}) bool {
			return f(key.(string))
		})
	}
}

// 发送分组消息
func (manager *Manager) SendGroupMessage(groupKey string, ack string, data interface{}) {
	manager.GetGroupClientList(groupKey, func(clientId string) bool {
		client, _ := manager.GetClientByClientId(clientId)
		client.SendMessage(ack, data)
		return true
	})
}
