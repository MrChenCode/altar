package context

import (
	"sync"
)

type MysqlClient struct {
	mysqlApi
	servers *sync.Map
}

type RedisClient struct {
	*Redis
	servers *sync.Map
}

type MongoClient struct {
	mongoApi
	servers *sync.Map
}

//获取指定mysql服务器的连接对象
//注意：不是切换mysql数据库(use db)
func (mc *MysqlClient) SelectDB(key string) (MysqlApi, bool) {
	s, ok := mc.servers.Load(key)
	if ok {
		return s.(mysqlApi), true
	}
	return nil, false
}

//获取指定的redis服务器的连接对象
//注意：不是切换redis数据库(select db)
func (rc *RedisClient) SelectDB(key string) (*Redis, bool) {
	server, ok := rc.servers.Load(key)
	if !ok {
		return nil, false
	}
	return server.(*Redis), true
}

//获取指定的mongo服务器连接对象
//注意：不是切换mongo数据库(use db)
func (mc *MongoClient) SelectDB(key string) (MongoApi, bool) {
	server, ok := mc.servers.Load(key)
	if !ok {
		return nil, false
	}
	return server.(mongoApi), true
}
