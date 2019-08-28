package context

import (
	"altar/application/core/config"
	"altar/logger"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"time"
)

//mysql curd相关基本操作
type mysqlCurd interface {
	//增删改查相关sql
	//query, 执行一条查询的sql，返回查询结果集的所有数据
	Query(query string, args ...interface{}) ([]map[string]string, error)

	//queryRow, 执行一条查询的sql，只返回结果集的第一条数据
	QueryRow(query string, args ...interface{}) (map[string]string, error)

	//queryResult,手动处理mysql结果集
	//TODO: 使用此函数，调用方必须手动正确关闭rows，否则会造成mysql资源泄漏
	QueryResult(query string, args ...interface{}) (*sql.Rows, error)

	//执行一条非查询类sql，比如insert、delete、update、alter等操作
	//返回insertid和影响行数接口
	Exec(query string, args ...interface{}) (sql.Result, error)
}

//mysql基础api
type mysqlApi interface {
	//设置mysql连接池最大连接数
	SetMaxOpenConns(n int)

	//连接池连接空闲的最大时间，如果连接空闲超过此时间d，则会被关闭
	//时间d不可大于mysql服务器的wait_timeout，会引发invalid connection(bad connection)错误
	SetConnMaxLifetime(d time.Duration)

	//连接池最大空闲链接数量
	SetMaxIdleConns(n int)

	//关闭，不能再用于任何mysql查询等操作
	Close() error

	//发送一个ping，检测活跃状态
	Ping() error

	//获取一些数据库统计信息
	Stats() sql.DBStats

	//启用mysql事务
	Begin(opts *sql.TxOptions) (MysqlTx, error)

	mysqlCurd
}

type MysqlApi interface {
	mysqlApi
}

//mysql事务
//如果增删改操作出现错误，会自动执行Rollback事务
//如果已经执行过Rollback或Commit，再次执行无效(但不会返回错误)
type MysqlTx interface {
	//回滚事务
	Rollback() error
	//提交事务
	Commit() error

	mysqlCurd
}

//mongo操作API
type mongoApi interface {
	//查询单条mongo文档
	//TODO: 必须对结果操作，否则会导致泄漏
	FindOne(colname string, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult

	//查询多条mongo文档
	//TODO: 必须手动对返回结果Close,否则会导致泄漏
	Find(colname string, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)

	//插入单条文档
	InsertOne(colname string, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)

	//插入多条文档
	InsertMany(colname string, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error)

	//统计符合查询条件的文档条数
	CountDocuments(colname string, filter interface{}, opts ...*options.CountOptions) (int64, error)

	//删除单条文档数据
	DeleteOne(colname string, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)

	//删除多条文档数据
	DeleteMany(colname string, filter []interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)

	//查找一条文档并删除, 并返回删除的文档
	//TODO: 必须对结果进行操作，否则会导致泄漏
	FindOneAndDelete(colname string, filter interface{}, opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult

	//查找一条文档并替换, 返回被替换的文档或返回新文档
	//TODO: 必须对结果进行操作，否则会导致泄漏
	FindOneAndReplace(colname string, filter interface{}, replacement interface{},
		opts ...*options.FindOneAndReplaceOptions) *mongo.SingleResult

	//查找一条文档并修改, 返回被修改的文档或者新文档
	//TODO: 必须对结果进行操作，否则会导致泄漏
	FindOneAndUpdate(colname string, filter interface{}, update interface{},
		opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult

	//根据查询条件修改一个文档
	UpdateOne(colname string, filter interface{}, update interface{},
		opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)

	//根据查询条件替换一个文档
	ReplaceOne(colname string, filter interface{},
		replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error)

	//修改所有符合查询条件的文档
	UpdateMany(colname string, filter interface{}, update interface{},
		opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
}

type MongoApi interface {
	mongoApi
}

func newRedisServer(c *config.Config, opt *config.RedisServer, logwf *logger.Logger) *Redis {
	redisOpt := &redis.Options{
		Network:            "tcp",
		Addr:               opt.Addr,
		Password:           opt.Pwd,
		DB:                 opt.DB,
		MaxRetries:         3,
		MinRetryBackoff:    8 * time.Millisecond,
		MaxRetryBackoff:    512,
		DialTimeout:        5 * time.Second,
		ReadTimeout:        3 * time.Second,
		WriteTimeout:       3 * time.Second,
		PoolSize:           opt.RedisPoolSize,
		MinIdleConns:       opt.RedisMinIdle,
		PoolTimeout:        c.RedisPoolTimeout,
		IdleTimeout:        c.RedisIdleTimeout,
		IdleCheckFrequency: c.RedisIdleCheckFrequency,
	}
	return newRedis(redis.NewClient(redisOpt), logwf)
}

func newMongoServer(msc *config.MongoServer, logwf *logger.Logger) (mongoApi, error) {
	var (
		connectTimeout = 10 * time.Second
		socketTimeout  = 5 * time.Second
		poolsize       = uint64(msc.Poolsize)
		minpoolsize    = uint64(msc.IdlePoolSize)

		readConcern  *readconcern.ReadConcern
		readPref     *readpref.ReadPref
		writeConcern *writeconcern.WriteConcern

		err error
	)
	switch msc.ReadConcern {
	case "local":
		readConcern = readconcern.Local()
	case "majority":
		readConcern = readconcern.Majority()
	default:
		return nil, errors.New("other readConcern types are not supported for the time being")
	}
	readPref, err = readpref.New(msc.ReadPreference)
	if err != nil {
		return nil, err
	}
	writeConcern = writeconcern.New(
		writeconcern.WTagSet(msc.WriteConcernW),
		writeconcern.J(msc.WriteConcernJ),
		writeconcern.WTimeout(5*time.Second),
	)

	timeoutCtx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(timeoutCtx, &options.ClientOptions{
		Hosts:           msc.Addrs,
		ConnectTimeout:  &connectTimeout,
		MaxConnIdleTime: &msc.IdleTimeout,
		MaxPoolSize:     &poolsize,
		MinPoolSize:     &minpoolsize,
		ReadConcern:     readConcern,
		ReadPreference:  readPref,
		ReplicaSet:      &msc.ReplicaSet,
		SocketTimeout:   &socketTimeout,
		WriteConcern:    writeConcern,
	})
	if err != nil {
		return nil, err
	}

	db := client.Database(msc.DB, &options.DatabaseOptions{
		ReadConcern:    readConcern,
		WriteConcern:   writeConcern,
		ReadPreference: readPref,
	})

	return newMongo(db, logwf), nil
}

func newMysqlServer(c *config.MysqlServer, logwf *logger.Logger) (mysqlApi, error) {
	h := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", c.User, c.Pwd, "tcp", c.Host, c.Port, c.DB)
	db, err := sql.Open("mysql", h)
	if err != nil {
		return nil, err
	}
	mysql := newMysql(db, logwf)
	mysql.SetMaxOpenConns(c.PoolSize)
	mysql.SetMaxIdleConns(c.IdleSize)
	mysql.SetConnMaxLifetime(c.IdleTimeout)
	return mysql, nil
}
