package context

import (
	"altar/application/config"
	"altar/application/logger"
	"sync"
)

//基础上下文资源，服务启动期间，仅初始化一次
type Context struct {
	Config *config.Config

	Mysql *MysqlClient
	Redis *RedisClient
	Mongo *MongoClient

	loginfo *logger.Logger
	logwf   *logger.Logger
}

func NewController(c *config.Config, log, logwf *logger.Logger) (*Context, error) {
	ctx := &Context{Config: c, loginfo: log, logwf: logwf}
	return ctx, ctx.init()
}

func (ctx *Context) WriteLogInfo(id string, kvs ...interface{}) {
	ctx.loginfo.Infow(id, kvs...)
}

func (ctx *Context) WriteLogError(id string, kvs ...interface{}) {
	ctx.logwf.Errorw(id, kvs...)
}

func (ctx *Context) LogSync() {
	_ = ctx.loginfo.Sync()
	_ = ctx.logwf.Sync()
}

func (ctx *Context) init() error {
	if ctx.Config == nil {
		panic("Invalid config")
	}
	if ctx.loginfo == nil || ctx.logwf == nil {
		panic("Invalid logger")
	}

	if ctx.Config.MysqlEnable && ctx.Mysql == nil {
		//初始化mysql链接
		mc := &MysqlClient{
			servers: &sync.Map{},
		}
		for _, opt := range ctx.Config.MysqlServers {
			mysql, err := newMysqlServer(opt, ctx.logwf)
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
			redisClient := newRedisServer(ctx.Config, opt, ctx.logwf)
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
			mongoClient, err := newMongoServer(opt, ctx.logwf)
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
