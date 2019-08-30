package context

import (
	"altar/application/logger"
	"context"
	"database/sql"
	"sync/atomic"
	"time"
)

const (
	//执行一条sql的超时时间
	MysqlExecTimeout = 3 * time.Second

	//执行sql超过此时间会记录一条wf日志
	MysqlExecWF = 300 * time.Millisecond
)

var (
	mysqlwf *logger.Logger
)

//Mysql操作API
//把原始的sql api封闭
type mysql struct {
	db *sql.DB
}

//mysql事务操作
type mysqlTx struct {
	txend int32
	tx    *sql.Tx
}

func newMysql(db *sql.DB, wf *logger.Logger) mysqlApi {
	mysqlwf = wf
	return &mysql{db: db}
}

//获取一个sql超时的上下文,超时时间由MYSQL_EXEC_TIMEOUT指定
func mysqlTimeoutCtx() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), MysqlExecTimeout)
	return ctx
}

//记录wf日志
func mysqlLogWF(query string, err error, d time.Duration) {
	if err == nil && d < MysqlExecWF {
		return
	}
	if err != nil {
		mysqlwf.Errorw("", "msg", err.Error(), "device", "mysql", "query", query, "usetime", d.Seconds())
	} else {
		mysqlwf.Warnw("", "msg", "slow", "device", "mysql", "query", query, "usetime", d.Seconds())
	}

}

func queryResult(rows *sql.Rows, irow bool) ([]map[string]string, error) {
	if err := rows.Err(); err != nil {
		return nil, err
	}
	ks, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	rs := make([]interface{}, len(ks))
	for i := 0; i < len(ks); i++ {
		var s string
		rs[i] = &s
	}
	var res []map[string]string
	if irow {
		if !rows.Next() {
			return res, nil
		}
		if err := rows.Scan(rs...); err != nil {
			return nil, err
		}
		vs := make(map[string]string)
		for i := 0; i < len(ks); i++ {
			vs[ks[i]] = *(rs[i].(*string))
		}
		res = append(res, vs)
	} else {
		for rows.Next() {
			if err := rows.Scan(rs...); err != nil {
				return nil, err
			}
			vs := make(map[string]string)
			for i := 0; i < len(ks); i++ {
				vs[ks[i]] = *(rs[i].(*string))
			}
			res = append(res, vs)
		}
	}
	return res, nil
}

//设置连接池打开最大的连接数
//如果n<=0，则不限制(不建议)
func (m *mysql) SetMaxOpenConns(n int) {
	m.db.SetMaxOpenConns(n)
}

//连接池连接空闲的最大时间，如果连接空闲超过此时间d，则会被关闭
//时间d不可大于mysql服务器的wait_timeout，会引发invalid connection(bad connection)错误
func (m *mysql) SetConnMaxLifetime(d time.Duration) {
	m.db.SetConnMaxLifetime(d)
}

//设置连接池空闲连接数
//如果n<=0,则不保留空闲连接（不建议）
func (m *mysql) SetMaxIdleConns(n int) {
	m.db.SetMaxIdleConns(n)
}

//关闭数据库
func (m *mysql) Close() error {
	return m.db.Close()
}

//PING
func (m *mysql) Ping() error {
	return m.db.PingContext(mysqlTimeoutCtx())
}

//获取一些统计信息
func (m *mysql) Stats() sql.DBStats {
	return m.db.Stats()
}

//执行一条非查询的sql语句
func (m *mysql) Exec(query string, args ...interface{}) (sql.Result, error) {
	start := time.Now()
	res, err := m.db.ExecContext(mysqlTimeoutCtx(), query, args...)
	mysqlLogWF(query, err, time.Now().Sub(start))

	return res, err
}

//查询所有结果
func (m *mysql) Query(query string, args ...interface{}) ([]map[string]string, error) {
	start := time.Now()
	v, err := m.query(query, false, args...)
	mysqlLogWF(query, err, time.Now().Sub(start))

	return v, err
}

//查询一条结果
func (m *mysql) QueryRow(query string, args ...interface{}) (map[string]string, error) {
	start := time.Now()
	v, err := m.query(query, true, args...)
	mysqlLogWF(query, err, time.Now().Sub(start))
	if err != nil {
		return nil, err
	}
	if len(v) > 0 {
		return v[0], nil
	}
	return nil, nil
}

//自己操作查询对象
//如果自己操作，必须手动关闭rows结果集，否则会造成结果集泄漏
func (m *mysql) QueryResult(query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()
	rows, err := m.db.QueryContext(mysqlTimeoutCtx(), query, args...)
	mysqlLogWF(query, err, time.Now().Sub(start))

	return rows, err
}

//启动事务，opts传递一个sql.TxOptions事务配置对象(如要默认，传nil):
//	TxOptions.Isolation事务隔离级别
//	TxOptions.ReadOnly是否开启只读事务(true/false)
func (m *mysql) Begin(opts *sql.TxOptions) (MysqlTx, error) {
	start := time.Now()
	tx, err := m.db.BeginTx(mysqlTimeoutCtx(), opts)
	mysqlLogWF("begin", err, time.Now().Sub(start))
	if err != nil {
		return nil, err
	}
	return &mysqlTx{txend: 0, tx: tx}, nil
}

func (m *mysql) query(query string, irow bool, args ...interface{}) ([]map[string]string, error) {
	rows, err := m.db.QueryContext(mysqlTimeoutCtx(), query, args...)
	if err != nil {
		return nil, err
	}
	//此处必须保证安全的关闭rows查询资源
	//对于row单行查询，也不允许再次读取剩余的数据(直接关闭资源)
	defer rows.Close()
	res, err := queryResult(rows, irow)
	if err != nil {
		return nil, err
	}
	return res, rows.Close()
}

//以下为事务操作//

//回滚事务
func (mt *mysqlTx) Rollback() error {
	if !atomic.CompareAndSwapInt32(&mt.txend, 0, 1) {
		return nil
	}
	start := time.Now()
	err := mt.tx.Rollback()
	mysqlLogWF("rollback", err, time.Now().Sub(start))

	return err
}

//提交事务
func (mt *mysqlTx) Commit() error {
	if !atomic.CompareAndSwapInt32(&mt.txend, 0, 1) {
		return nil
	}
	start := time.Now()
	err := mt.tx.Commit()
	mysqlLogWF("commit", err, time.Now().Sub(start))

	return err
}

func (mt *mysqlTx) Query(query string, args ...interface{}) ([]map[string]string, error) {
	start := time.Now()
	v, err := mt.query(query, false, args...)
	mysqlLogWF(query, err, time.Now().Sub(start))

	return v, err
}

func (mt *mysqlTx) QueryRow(query string, args ...interface{}) (map[string]string, error) {
	start := time.Now()
	v, err := mt.query(query, true, args...)
	mysqlLogWF(query, err, time.Now().Sub(start))
	if err != nil {
		return nil, err
	}
	if len(v) > 0 {
		return v[0], nil
	}
	return nil, nil
}

//手动操作结果集
func (mt *mysqlTx) QueryResult(query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()
	rows, err := mt.tx.QueryContext(mysqlTimeoutCtx(), query, args...)
	mysqlLogWF(query, err, time.Now().Sub(start))

	return rows, err
}

func (mt *mysqlTx) Exec(query string, args ...interface{}) (sql.Result, error) {
	start := time.Now()
	res, err := mt.tx.ExecContext(mysqlTimeoutCtx(), query, args...)
	mysqlLogWF(query, err, time.Now().Sub(start))

	//如果执行增删改操作出现错误，自动回滚事务
	if err != nil {
		_ = mt.Rollback()
	}

	return res, err
}

func (mt *mysqlTx) query(query string, irow bool, args ...interface{}) ([]map[string]string, error) {
	rows, err := mt.tx.QueryContext(mysqlTimeoutCtx(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res, err := queryResult(rows, irow)
	if err != nil {
		return nil, err
	}
	return res, rows.Close()
}
