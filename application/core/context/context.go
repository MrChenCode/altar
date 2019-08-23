package context

import (
	"gitlab.baidu-shucheng.com/shaohua/bloc/application/core/config"
	"gitlab.baidu-shucheng.com/shaohua/bloc/logger"
	"sync"
)

type BasicContext struct {
	Config *config.Config

	Mysql *MysqlClient
	Redis *RedisClient
	Mongo *MongoClient

	Logger   *logger.Logger
	Loggerwf *logger.Logger
}

func (ctx *BasicContext) Init() error {
	if ctx.Config == nil {
		panic("Invalid config")
	}
	if ctx.Logger == nil || ctx.Loggerwf == nil {
		panic("Invalid logger")
	}

	if ctx.Config.MysqlEnable && ctx.Mysql == nil {
		//初始化mysql链接
		mc := &MysqlClient{
			servers: &sync.Map{},
		}
		for _, opt := range ctx.Config.MysqlServers {
			mysql, err := newMysqlServer(opt, ctx.Loggerwf)
			if err != nil {
				return err
			}
			mc.servers.Store(opt.Name, mysql)
			if opt.Name == ctx.Config.MysqlDefaultServer {
				mc.mysqlApi = mysql
			}
		}
		ctx.Mysql = mc
	}

	if ctx.Config.RedisEnable && ctx.Redis == nil {
		//初始化redis
		rc := &RedisClient{
			servers: &sync.Map{},
		}
		for _, opt := range ctx.Config.RedisServers {
			redisClient := newRedisServer(ctx.Config, opt, ctx.Loggerwf)
			rc.servers.Store(opt.Name, redisClient)
			if opt.Name == ctx.Config.RedisDefaultServer {
				rc.Redis = redisClient
			}
		}
		ctx.Redis = rc
	}

	if ctx.Config.MongoEnable && ctx.Mongo == nil {
		//初始化mongo
		mc := &MongoClient{
			servers: &sync.Map{},
		}

		for _, opt := range ctx.Config.MongoServers {
			mongoClient, err := newMongoServer(opt, ctx.Loggerwf)
			if err != nil {
				return err
			}
			mc.servers.Store(opt.Name, mongoClient)
			if opt.Name == ctx.Config.MongoDefaultServer {
				mc.mongoApi = mongoClient
			}
		}
		ctx.Mongo = mc
	}

	return nil
}
