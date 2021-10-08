package ws

import "errors"

// 添加客户端unionId
func (manager *Manager) AddClientUnionId(client *Client, unionId string) {
	client.UnionId = unionId
	manager.unions.Store(unionId, client)
}

// 通过unionId删除client
func (manager *Manager) DelClientUnion(client *Client, unionId string) {
	client.UnionId = ""
	manager.delClientUnion(unionId)
}
func (manager *Manager) delClientUnion(unionId string) {
	manager.unions.Delete(unionId)
}

// 通过unionId获取client
func (manager *Manager) GetClientByUnionId(unionId string) (*Client, error) {
	if client, ok := manager.unions.Load(unionId); !ok {
		return nil, errors.New("客户端不存在")
	} else {
		return client.(*Client), nil
	}
}
