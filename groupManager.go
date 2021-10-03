package ws

import "sync"

// 添加到分组
func (manager *Manager) AddGroupClient(groupKey string, client *Client) {
	if _, ok := client.GroupList[groupKey]; ok {
		return
	}
	manager.addGroupClient(groupKey, client.Uid)
	client.GroupList[groupKey] = true
}

//删除分组里的客户端
func (manager *Manager) DelGroupClient(groupKey string, client *Client) {
	if _, ok := client.GroupList[groupKey]; !ok {
		return
	}
	manager.delGroupClient(groupKey, client.Uid)
	delete(client.GroupList, groupKey)
}

// 添加到分组
func (manager *Manager) addGroupClient(groupKey string, uid UID) {
	uidMap, _ := manager.groups.LoadOrStore(groupKey, new(syncMap))
	sm := uidMap.(*syncMap)
	_, loaded := sm.LoadOrStore(uid, 1)
	if !loaded {
		sm.count++
	}
}

// 删除分组里的客户端
func (manager *Manager) delGroupClient(groupKey string, uid UID) {
	if uidMap, ok := manager.groups.Load(groupKey); ok {
		sm := uidMap.(*syncMap)
		_, loaded := sm.LoadAndDelete(uid)
		if loaded {
			sm.count--
		}
	}
}

// 获取分组所有的客户端数量
func (manager *Manager) GetGroupClientCount(groupKey string) int {
	if uidMap, ok := manager.groups.Load(groupKey); ok {
		return uidMap.(*syncMap).count
	}
	return 0
}

// 获取分组的成员
func (manager *Manager) GetGroupClientList(groupKey string, f func(uid UID) bool) {
	if uidMap, ok := manager.groups.Load(groupKey); ok {
		uidMap.(*sync.Map).Range(func(key, value interface{}) bool {
			return f(key.(UID))
		})
	}
}

// 发送分组消息
func (manager *Manager) SendGroupMessage(groupKey string, ack string, data interface{}) {
	manager.GetGroupClientList(groupKey, func(uid UID) bool {
		client, _ := manager.GetClientByUid(uid)
		client.SendMessage(ack, data)
		return true
	})
}
