package context

import (
	"fmt"
	"github.com/go-redis/redis"
	"gitlab.baidu-shucheng.com/shaohua/bloc/logger"
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
		if err != nil || usetime > RedisExecTimeout {
			cmdinfo := cmd.Name()
			for _, c := range cmd.Args() {
				cmdinfo += " "
				cmdinfo += fmt.Sprint(c)
			}
			if err != nil {
				r.logwf.Errorw(err.Error(), "device", "redis", "query", cmdinfo, "usetime", usetime)
			} else {
				r.logwf.Warnw("slow", "device", "redis", "query", cmdinfo, "usetime", usetime)
			}
		}
		return err
	}
}
