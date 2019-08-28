package context

import (
	"altar/logger"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

const (
	RedisExecTimeout = 100 * time.Millisecond
)

type Redis struct {
	*redis.Client
	logwf *logger.Logger
}

func newRedis(client *redis.Client, wf *logger.Logger) *Redis {
	r := &Redis{Client: client, logwf: wf}
	r.WrapProcess(r.execRedisCmd)

	return r
}

func (r *Redis) execRedisCmd(oldProcess func(cmd redis.Cmder) error) func(cmd redis.Cmder) error {
	return func(cmd redis.Cmder) error {
		start := time.Now()
		err := oldProcess(cmd)
		usetime := time.Now().Sub(start)
		//如果err == redis.Nil, 认为是查找的key不存在，不认为是错误
		if err == redis.Nil {
			err = nil
		}
		if err != nil || usetime > RedisExecTimeout {
			cmdinfo := cmd.Name()
			for _, c := range cmd.Args() {
				cmdinfo += " "
				cmdinfo += fmt.Sprint(c)
			}
			if err != nil {
				r.logwf.Errorw("", "msg", err.Error(), "device", "redis", "query", cmdinfo, "usetime", usetime)
			} else {
				r.logwf.Warnw("", "msg", "slow", "device", "redis", "query", cmdinfo, "usetime", usetime)
			}
		}
		return err
	}
}
