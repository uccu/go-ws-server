package ws

import "sync"

// 添加到系统
func (manager *Manager) AddSystemClient(systemKey string, client *Client) {
	if _, ok := client.SystemList[systemKey]; ok {
		return
	}
	manager.addSystemClient(systemKey, client.Uid)
	client.SystemList[systemKey] = true
}

// 删除系统里的客户端
func (manager *Manager) DelSystemClient(systemKey string, client *Client) {
	if _, ok := client.SystemList[systemKey]; !ok {
		return
	}
	manager.delSystemClient(systemKey, client.Uid)
	delete(client.SystemList, systemKey)
}

// 添加到系统
func (manager *Manager) addSystemClient(systemKey string, uid UID) {
	uidMap, _ := manager.systems.LoadOrStore(systemKey, new(syncMap))
	sm := uidMap.(*syncMap)
	_, loaded := sm.LoadOrStore(uid, 1)
	if !loaded {
		sm.count++
	}
}

// 删除系统里的客户端
func (manager *Manager) delSystemClient(systemKey string, uid UID) {
	if uidMap, ok := manager.systems.Load(systemKey); ok {
		sm := uidMap.(*syncMap)
		_, loaded := sm.LoadAndDelete(uid)
		if loaded {
			sm.count--
		}
	}
}

// 获取系统所有的客户端数量
func (manager *Manager) GetSystemClientCount(systemKey string) int {
	if uidMap, ok := manager.systems.Load(systemKey); ok {
		return uidMap.(*syncMap).count
	}
	return 0
}

// 获取系统的成员
func (manager *Manager) GetSystemClientList(systemKey string, f func(uid UID) bool) {
	if uidMap, ok := manager.systems.Load(systemKey); ok {
		uidMap.(*sync.Map).Range(func(key, value interface{}) bool {
			return f(key.(UID))
		})
	}
}

// 发送系统消息
func (manager *Manager) SendSystemMessage(systemKey string, ack string, data interface{}) {
	manager.GetSystemClientList(systemKey, func(uid UID) bool {
		client, _ := manager.GetClientByUid(uid)
		client.SendMessage(ack, data)
		return true
	})
}
