package context

import (
	"context"
	"fmt"
	"gitlab.baidu-shucheng.com/shaohua/bloc/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	MongoExecTimeout = 3 * time.Second
	MongoExecWF      = 300 * time.Millisecond
)

type mongodb struct {
	db *mongo.Database
	wf *logger.Logger
}

func newMongo(db *mongo.Database, wf *logger.Logger) mongoApi {
	return &mongodb{db: db, wf: wf}
}

func mongoTimeoutCtx() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), MongoExecTimeout)
	return ctx
}

//记录wf日志
func (m *mongodb) logwf(err error, d time.Duration, colname, exectype string, filter interface{}) {

	if err != nil {
		mysqlwf.Errorw(err.Error(),
			"device", "mongo",
			"colname", colname,
			"exectype", exectype,
			"usetime", d.Seconds(),
			"filter", fmt.Sprint(filter),
		)
	} else {
		mysqlwf.Warnw("slow",
			"device", "mongo",
			"colname", colname,
			"usetime", d.Seconds(),
			"filter", fmt.Sprint(filter),
		)
	}
}

//查询一条数据
//colname: 文档名称(类似于mysql表名)
func (m *mongodb) FindOne(colname string, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	start := time.Now()
	v := m.db.Collection(colname).FindOne(mongoTimeoutCtx(), filter, opts...)
	usetime := time.Now().Sub(start)

	err := v.Err()
	if err != nil || usetime > MongoExecWF {
		m.logwf(err, usetime, colname, "FindOne", filter)
	}
	return v
}

//查询全部结果集
//colname: 文档名称(类似于mysql表名)
func (m *mongodb) Find(colname string, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	start := time.Now()
	cur, err := m.db.Collection(colname).Find(mongoTimeoutCtx(), filter, opts...)
	usetime := time.Now().Sub(start)

	if err != nil || usetime > MongoExecWF {
		m.logwf(err, usetime, colname, "Find", filter)
	}

	return cur, err
}

//插入单条
//colname: 文档名称(类似于mysql表名)
func (m *mongodb) InsertOne(colname string, document interface{},
	opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {

	start := time.Now()
	res, err := m.db.Collection(colname).InsertOne(mongoTimeoutCtx(), document, opts...)
	usetime := time.Now().Sub(start)

	if err != nil || usetime > MongoExecWF {
		m.logwf(err, usetime, colname, "InsertOne", nil)
	}

	return res, err
}

//插入多条
//colname: 文档名称(类似于mysql表名)
func (m *mongodb) InsertMany(colname string, documents []interface{},
	opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {

	start := time.Now()
	res, err := m.db.Collection(colname).InsertMany(mongoTimeoutCtx(), documents, opts...)
	usetime := time.Now().Sub(start)

	if err != nil || usetime > MongoExecWF {
		m.logwf(err, usetime, colname, "InsertMany", nil)
	}

	return res, err
}

//统计文档
func (m *mongodb) CountDocuments(colname string, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	start := time.Now()
	n, err := m.db.Collection(colname).CountDocuments(mongoTimeoutCtx(), filter, opts...)
	usetime := time.Now().Sub(start)

	if err != nil || usetime > MongoExecWF {
		m.logwf(err, usetime, colname, "CountDocuments", filter)
	}

	return n, err
}

//删除单条
//colname: 文档名称(类似于mysql表名)
func (m *mongodb) DeleteOne(colname string, filter interface{},
	opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {

	start := time.Now()
	res, err := m.db.Collection(colname).DeleteOne(mongoTimeoutCtx(), filter, opts...)
	usetime := time.Now().Sub(start)

	if err != nil || usetime > MongoExecWF {
		m.logwf(err, usetime, colname, "DeleteOne", filter)
	}

	return res, err
}

//删除多条
//colname: 文档名称(类似于mysql表名)
func (m *mongodb) DeleteMany(colname string, filter []interface{},
	opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {

	start := time.Now()
	res, err := m.db.Collection(colname).DeleteMany(mongoTimeoutCtx(), filter, opts...)
	usetime := time.Now().Sub(start)

	if err != nil || usetime > MongoExecWF {
		m.logwf(err, usetime, colname, "DeleteMany", filter)
	}

	return res, err
}

//查询单个文档后删除，返回被删除的文档
//colname: 文档名称(类似于mysql表名)
func (m *mongodb) FindOneAndDelete(colname string, filter interface{},
	opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult {

	start := time.Now()
	res := m.db.Collection(colname).FindOneAndDelete(mongoTimeoutCtx(), filter, opts...)
	usetime := time.Now().Sub(start)

	err := res.Err()
	if err != nil || usetime > MongoExecWF {
		m.logwf(err, usetime, colname, "FindOneAndDelete", filter)
	}

	return res
}

//查询单个文档后替换，返回被替换的文档或者新文档
//colname: 文档名称(类似于mysql表名)
func (m *mongodb) FindOneAndReplace(colname string, filter interface{},
	replacement interface{}, opts ...*options.FindOneAndReplaceOptions) *mongo.SingleResult {

	start := time.Now()
	res := m.db.Collection(colname).FindOneAndReplace(mongoTimeoutCtx(), filter, replacement, opts...)
	usetime := time.Now().Sub(start)

	err := res.Err()
	if err != nil || usetime > MongoExecWF {
		m.logwf(err, usetime, colname, "FindOneAndReplace", filter)
	}

	return res
}

//查询单个文档后修改，返回被修改的文档或者新文档
//colname: 文档名称(类似于mysql表名)
func (m *mongodb) FindOneAndUpdate(colname string, filter interface{},
	update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {

	start := time.Now()
	res := m.db.Collection(colname).FindOneAndUpdate(mongoTimeoutCtx(), filter, update, opts...)
	usetime := time.Now().Sub(start)

	err := res.Err()
	if err != nil || usetime > MongoExecWF {
		m.logwf(err, usetime, colname, "FindOneAndUpdate", filter)
	}

	return res
}

//更新一个文档
//colname: 文档名称(类似于mysql表名)
func (m *mongodb) UpdateOne(colname string, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {

	start := time.Now()
	res, err := m.db.Collection(colname).UpdateOne(mongoTimeoutCtx(), filter, update, opts...)
	usetime := time.Now().Sub(start)

	if err != nil || usetime > MongoExecWF {
		m.logwf(err, usetime, colname, "UpdateOne", filter)
	}

	return res, err
}

//替换一个符合查询添加的文档
//colname: 文档名称(类似于mysql表名)
func (m *mongodb) ReplaceOne(colname string, filter interface{},
	replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {

	start := time.Now()
	res, err := m.db.Collection(colname).ReplaceOne(mongoTimeoutCtx(), filter, replacement, opts...)
	usetime := time.Now().Sub(start)

	if err != nil || usetime > MongoExecWF {
		m.logwf(err, usetime, colname, "ReplaceOne", filter)
	}

	return res, err
}

//修改所有符合查询条件的文档
//colname: 文档名称(类似于mysql表名)
func (m *mongodb) UpdateMany(colname string, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {

	start := time.Now()
	res, err := m.db.Collection(colname).UpdateMany(mongoTimeoutCtx(), filter, update, opts...)
	usetime := time.Now().Sub(start)

	if err != nil || usetime > MongoExecWF {
		m.logwf(err, usetime, colname, "UpdateMany", filter)
	}

	return res, err
}
