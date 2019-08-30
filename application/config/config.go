package config

import (
	"altar/application/logger"
	"altar/application/system"
	"errors"
	"fmt"
	"github.com/go-ini/ini"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	RootPath    string
	DebugOnline = "online"
	DebugQA     = "qa"
)

type Config struct {
	//运行模式, online/qa
	Debug string

	//http服务监听地址, ip:port(10.3.138.32:8888)
	HttpServerAddr string

	//日志路径
	LogPath string

	//日志文件名称
	//wf日志会在末尾加上.wf
	LogFileName string

	//日志切割周期
	//	0-每小时切割
	//	1-每天切割
	//	2-每周切割
	//	3-每月切割
	//	4-永不切割
	LogCatTime logger.Cattime

	//日志保存天数
	LogRetainDay int

	//pid
	PidFile string

	//是否启用redis
	RedisEnable bool
	//是否启用mysql
	MysqlEnable bool
	//是否启用mongo
	MongoEnable bool

	//如果连接池所有的连接都繁忙，等待空闲连接的时间
	RedisPoolTimeout time.Duration
	//关闭空闲连接等待时间
	RedisIdleTimeout time.Duration
	//检查空闲连接时间
	RedisIdleCheckFrequency time.Duration

	//redis服务器
	RedisDefaultServer string
	RedisServers       []*RedisServer

	//mysql服务器
	MysqlDefaultServer string
	MysqlServers       []*MysqlServer

	//mongo服务器
	MongoDefaultServer string
	MongoServers       []*MongoServer
}

type RedisServer struct {
	Name          string
	Addr          string
	Pwd           string
	DB            int
	RedisPoolSize int
	RedisMinIdle  int
}

type MysqlServer struct {
	Name        string
	Host        string
	Port        int
	User        string
	Pwd         string
	DB          string
	PoolSize    int
	IdleSize    int
	IdleTimeout time.Duration
}

type MongoServer struct {
	Name  string
	Addrs []string
	DB    string

	ReplicaSet     string
	ReadConcern    string
	WriteConcernW  string
	WriteConcernJ  bool
	ReadPreference readpref.Mode

	Poolsize     int
	IdlePoolSize int
	IdleTimeout  time.Duration
}

//测试语法
func SyntaxCheck(path string) error {
	file, err := getDebugFile(path)
	if err != nil {
		return err
	}
	_, err = ini.Load(file)
	return err
}

//获取具体的配置文件
func getDebugFile(path string) (string, error) {
	cfg, err := ini.Load(path)
	if err != nil {
		return "", err
	}
	se := cfg.Section("running")
	runDebug := se.Key("running").String()
	inifile := ""
	switch runDebug {
	case DebugQA:
		inifile = se.Key("qa_ini").String()
	case DebugOnline:
		inifile = se.Key("online_ini").String()
	default:
		return "", fmt.Errorf("config: nnknown key running_debug(%s)", runDebug)
	}
	if inifile == "" {
		return "", errors.New("config: no valid configuration file was found")
	}
	if filepath.IsAbs(inifile) {
		return inifile, nil
	}
	files := []string{
		inifile,
		filepath.Join(RootPath, inifile),
	}
	for _, ff := range files {
		fi, err := os.Stat(ff)
		//如果存在，且不是目录，则使用这个
		if err == nil {
			if !fi.IsDir() {
				return ff, nil
			}
		}
	}
	return filepath.Join(RootPath, inifile), nil
}

func NewConfig(path string) (*Config, error) {
	inifile, err := getDebugFile(path)
	if err != nil {
		return nil, err
	}
	cfg, err := ini.Load(inifile)
	if err != nil {
		return nil, err
	}
	c := new(Config)

	//读取services配置
	section := cfg.Section("services")

	c.PidFile = section.Key("pid_file").String()
	c.HttpServerAddr = section.Key("http_server_addr").String()
	c.LogFileName = section.Key("log_filename").String()
	cattime := section.Key("log_cat_time").String()
	switch cattime {
	case "hour":
		c.LogCatTime = logger.CAT_HOUR
	case "day":
		c.LogCatTime = logger.CAT_DAY
	case "week":
		c.LogCatTime = logger.CAT_WEEK
	case "method":
		c.LogCatTime = logger.CAT_MONTH
	}
	c.LogRetainDay, _ = section.Key("log_retain_day").Int()

	logpath := section.Key("log_path").String()
	logpath, err = directoryPermissions(logpath)
	if err != nil {
		return nil, err
	}
	c.LogPath = logpath

	section = cfg.Section("redis")
	c.RedisEnable, err = section.Key("redis_enable").Bool()
	if err != nil {
		return nil, err
	}
	//如果开启了redis，则读取redis配置
	if c.RedisEnable {
		if err := getRedisClient(section, c); err != nil {
			return nil, err
		}
	}

	section = cfg.Section("mysql")
	c.MysqlEnable, err = section.Key("mysql_enable").Bool()
	if err != nil {
		return nil, err
	}
	//如果开启了mysql，则读取mysql配置
	if c.MysqlEnable {
		if err := getMysqlClient(section, c); err != nil {
			return nil, err
		}
	}

	section = cfg.Section("mongo")
	c.MongoEnable, err = section.Key("mongo_enable").Bool()
	if err != nil {
		return nil, err
	}
	//如果开启mongo，则读取mongo配置
	if c.MongoEnable {
		if err := getMongoClient(section, c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func getMongoClient(s *ini.Section, c *Config) error {
	replicaSet := s.Key("mongo_replica_set").String()
	if replicaSet == "" {
		return errors.New("config mongo error: invalid mongo_replica_set")
	}
	readConcern := s.Key("mongo_read_concern").String()
	if readConcern == "" {
		readConcern = "local"
	}
	if readConcern != "local" && readConcern != "majority" {
		return fmt.Errorf("config mongo error: unknown mongo_read_concern configuration item (%s)", readConcern)
	}
	writeConcernW := s.Key("mongo_write_concern_w").String()
	if writeConcernW == "" {
		writeConcernW = "1"
	}
	if writeConcernW != "0" && writeConcernW != "1" && writeConcernW != "majority" {
		return fmt.Errorf("config mongo error: unknown mongo_write_concern_w configuration item (%s)", writeConcernW)
	}
	writeConcernJ, err := s.Key("mongo_write_concern_j").Bool()
	if err != nil {
		return err
	}
	readPreference := s.Key("mongo_read_preference").String()
	if readPreference == "" {
		readPreference = "primary"
	}
	mode, err := readpref.ModeFromString(readPreference)
	if err != nil {
		return err
	}
	poolsize, _ := s.Key("mongo_poolsize").Int()
	if poolsize < 1 || poolsize > 65535 {
		poolsize = 128
	}
	idlesize, _ := s.Key("mongo_idle_poolsize").Int()
	if idlesize < 1 || idlesize > poolsize {
		idlesize = 10
	}
	idletimeout, _ := s.Key("mongo_idle_timeout").Int()
	if idletimeout < 1 || idletimeout > 3600 {
		idletimeout = 300
	}

	defaultServer := s.Key("mongo_default_server").String()
	if defaultServer == "" {
		return errors.New("config mongo error: the default mongo server must be set")
	}
	defaultIsExist := false
	mongoServers := s.ChildSections()
	for _, ms := range mongoServers {
		servername := ms.Key("mongo_server_name").String()
		if servername == "" {
			return errors.New("config mongo error: mongo_server_name invalid")
		}
		addrsStr := ms.Key("mongo_addrs").String()
		if addrsStr == "" {
			return errors.New("config mongo error: mongo_addrs invalid")
		}
		addrsTemp := strings.Split(addrsStr, ",")
		var addrs []string
		for _, addr := range addrsTemp {
			if addr == "" {
				continue
			}
			if _, _, err := net.SplitHostPort(addr); err != nil {
				return err
			}
			addrs = append(addrs, addr)
		}
		if len(addrs) == 0 {
			return errors.New("config mongo error: invalid mongo_addrs")
		}
		db := ms.Key("mongo_db").String()
		if db == "" {
			return errors.New("config mongo error: invalid mongo_db")
		}

		//覆盖配置
		msreplicaSet := ms.Key("mongo_replica_set").String()
		if msreplicaSet == "" {
			msreplicaSet = replicaSet
		}
		msreadConcern := ms.Key("mongo_read_concern").String()
		if msreadConcern == "" {
			msreadConcern = readConcern
		}
		if msreadConcern != "local" && msreadConcern != "majority" {
			return fmt.Errorf("config mongo error: unknown mongo_read_concern configuration item (%s)", msreadConcern)
		}
		mswriteConcernW := ms.Key("mongo_write_concern_w").String()
		if mswriteConcernW == "" {
			mswriteConcernW = writeConcernW
		}
		if mswriteConcernW != "0" && mswriteConcernW != "1" && mswriteConcernW != "majority" {
			return fmt.Errorf("config mongo error: unknown mongo_write_concern_w configuration item (%s)", mswriteConcernW)
		}
		mswriteConcernJStr := ms.Key("mongo_write_concern_j")
		var mswriteConcernJ bool
		if mswriteConcernJStr.String() == "" {
			mswriteConcernJ = writeConcernJ
		} else {
			mswriteConcernJ, err = mswriteConcernJStr.Bool()
			if err != nil {
				return err
			}
		}
		msreadPreference := ms.Key("mongo_read_preference").String()
		var msmode readpref.Mode
		if readPreference == "" {
			msmode = mode
		} else {
			msmode, err = readpref.ModeFromString(msreadPreference)
			if err != nil {
				return err
			}
		}
		mspoolsize, _ := ms.Key("mongo_poolsize").Int()
		if mspoolsize < 1 || mspoolsize > 65535 {
			mspoolsize = poolsize
		}
		msidlesize, _ := ms.Key("mongo_idle_poolsize").Int()
		if msidlesize < 1 || msidlesize > mspoolsize {
			msidlesize = idlesize
		}
		//TODO:mongodb官方发生了错误，MinPoolSize会引发一个死锁问题
		//TODO:具体参照：https://jira.mongodb.org/browse/GODRIVER-1234
		//TODO:我们这里首先把最小连接池设置为0，后续再做更改
		msidlesize = 0

		msidletimeout, _ := ms.Key("mongo_idle_timeout").Int()
		if msidletimeout < 1 || msidletimeout > 3600 {
			msidletimeout = idletimeout
		}
		//end覆盖配置

		c.MongoServers = append(c.MongoServers, &MongoServer{
			Name:           servername,
			Addrs:          addrs,
			DB:             db,
			ReplicaSet:     msreplicaSet,
			ReadConcern:    msreadConcern,
			WriteConcernW:  mswriteConcernW,
			WriteConcernJ:  mswriteConcernJ,
			ReadPreference: msmode,
			Poolsize:       mspoolsize,
			IdlePoolSize:   msidlesize,
			IdleTimeout:    time.Duration(msidletimeout) * time.Second,
		})
		if !defaultIsExist && defaultServer == servername {
			defaultIsExist = true
		}
	}
	if !defaultIsExist {
		return fmt.Errorf("config mongo error: The default mongo server (%s) does not exist", defaultServer)
	}
	c.MongoDefaultServer = defaultServer
	return nil
}

func getMysqlClient(s *ini.Section, c *Config) error {
	//处理默认连接池配置
	poolsize, _ := s.Key("mysql_poolsize").Int()
	if poolsize < 1 || poolsize > 65535 {
		poolsize = 128
	}
	idlesize, _ := s.Key("mysql_idle_poolsize").Int()
	if idlesize < 1 || idlesize > poolsize {
		idlesize = 10
	}
	idletimeout, _ := s.Key("mysql_idletimeout").Int()
	//mysql server默认8h
	if idletimeout < 1 || idletimeout > (3600*8) {
		idletimeout = 300
	}
	defaultServer := s.Key("mysql_default_server").String()
	if defaultServer == "" {
		return errors.New("config mysql error: the default mysql server must be set")
	}
	defaultIsExist := false

	//处理mysql服务器信息
	mysqlServers := s.ChildSections()
	for _, ms := range mysqlServers {
		servername := ms.Key("mysql_server_name").String()
		if servername == "" {
			return errors.New("config mysql error: mysql_server_name invalid")
		}
		host := ms.Key("mysql_host").String()
		port, err := ms.Key("mysql_port").Int()
		if err != nil {
			return err
		}
		username := ms.Key("mysql_username").String()
		password := ms.Key("mysql_password").String()
		dbname := ms.Key("mysql_dbname").String()

		//处理连接池大小信息
		mspoolsize, _ := ms.Key("mysql_poolsize").Int()
		if mspoolsize < 1 || mspoolsize > 65535 {
			mspoolsize = poolsize
		}
		msidlesize, _ := ms.Key("mysql_idlesize").Int()
		if msidlesize < 1 || msidlesize > mspoolsize {
			msidlesize = idlesize
		}
		msidletimeout, _ := ms.Key("mysql_idletimeout").Int()
		if msidletimeout < 1 || msidletimeout > (3600*8) {
			msidletimeout = idletimeout
		}

		c.MysqlServers = append(c.MysqlServers, &MysqlServer{
			Name:        servername,
			Host:        host,
			Port:        port,
			User:        username,
			Pwd:         password,
			DB:          dbname,
			PoolSize:    mspoolsize,
			IdleSize:    msidlesize,
			IdleTimeout: time.Duration(msidletimeout) * time.Second,
		})
		if !defaultIsExist && defaultServer == servername {
			defaultIsExist = true
		}
	}
	if !defaultIsExist {
		return fmt.Errorf("config mysql error: The default mysql server (%s) does not exist", defaultServer)
	}
	c.MysqlDefaultServer = defaultServer
	return nil
}

func getRedisClient(s *ini.Section, c *Config) error {
	//处理连接池大小信息
	poolsize, _ := s.Key("redis_poolsize").Int()
	if poolsize < 1 || poolsize > 65535 {
		poolsize = 128
	}

	minIdleNum, _ := s.Key("redis_idle_poolsize").Int()
	if minIdleNum < 1 || minIdleNum > poolsize {
		minIdleNum = 10
	}

	//处理控制连接池的配置信息
	pooltimeout, _ := s.Key("redis_pooltimeout").Int()
	idletimeout, _ := s.Key("redis_idletimeout").Int()
	idlecheckfrequency, _ := s.Key("redis_idlecheckfrequency").Int()
	if pooltimeout <= 0 || pooltimeout > 120 {
		pooltimeout = 5
	}
	if idletimeout <= 0 || idletimeout > 1800 {
		idletimeout = 300
	}
	if idlecheckfrequency <= 0 || idlecheckfrequency > 600 {
		idlecheckfrequency = 60
	}
	c.RedisPoolTimeout = time.Duration(pooltimeout) * time.Second
	c.RedisIdleTimeout = time.Duration(idletimeout) * time.Second
	c.RedisIdleCheckFrequency = time.Duration(idlecheckfrequency) * time.Second

	defaultServer := s.Key("redis_default_server").String()
	if defaultServer == "" {
		return errors.New("config redis error: the default redis server must be set")
	}
	defaultIsExist := false

	//处理redis服务器信息
	redisServers := s.ChildSections()
	for _, rs := range redisServers {
		servername := rs.Key("redis_server_name").String()
		if servername == "" {
			return errors.New("config redis error: redis_server_name invalid")
		}
		addr := rs.Key("redis_address").String()
		pwd := rs.Key("redis_password").String()
		db, err := rs.Key("redis_db").Int()
		if err != nil {
			return err
		}
		//处理连接池大小信息
		rspoolsize, _ := rs.Key("redis_poolsize").Int()
		if rspoolsize < 1 || rspoolsize > 65535 {
			rspoolsize = poolsize
		}

		rsminIdleNum, _ := rs.Key("redis_idle_poolsize").Int()
		if rsminIdleNum < 1 || rsminIdleNum > rspoolsize {
			rsminIdleNum = minIdleNum
		}

		c.RedisServers = append(c.RedisServers, &RedisServer{
			Name:          servername,
			Addr:          addr,
			Pwd:           pwd,
			DB:            db,
			RedisPoolSize: rspoolsize,
			RedisMinIdle:  rsminIdleNum,
		})
		if !defaultIsExist && defaultServer == servername {
			defaultIsExist = true
		}
	}
	if !defaultIsExist {
		return fmt.Errorf("config redis error: The default redis server (%s) does not exist", defaultServer)
	}
	c.RedisDefaultServer = defaultServer
	return nil
}

func directoryPermissions(dir string) (path string, err error) {
	if dir == "" {
		err = errors.New("config: invalid path")
		return
	}
	if !filepath.IsAbs(dir) && RootPath != "" {
		dir = filepath.Join(RootPath, dir)
	}
	path, err = filepath.Abs(dir)
	if err != nil {
		return
	}
	fileinfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = fmt.Errorf("config: path(%s) directory does not exist", path)
			return
		}
		return
	}
	if !fileinfo.IsDir() {
		err = fmt.Errorf("config: path(%s) is not a directory", path)
		return
	}
	if system.SyscallAccess(path, system.O_RDWR) != nil {
		err = fmt.Errorf("config: path(%s) lack of directory permissions", path)
		return
	}

	return path, nil
}
