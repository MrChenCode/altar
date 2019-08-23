package core

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func (core *Core) GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		c.Next()

		latency := time.Now().Sub(start)

		errinfo := c.Errors.ByType(gin.ErrorTypePrivate).String()
		code := c.Writer.Status()

		//2xx and 3xx success!
		if code >= 200 && code < 400 {
			core.logger.Infow(path, "method", c.Request.Method, "httpcode", c.Writer.Status(), "errinfo", errinfo,
				"clientip", c.ClientIP(), "latency", latency)
		} else if code >= 400 && code < 600 {
			//4xx and 5xx error!
			core.logger.Warnw(path, "method", c.Request.Method, "httpcode", c.Writer.Status(), "errinfo", errinfo,
				"clientip", c.ClientIP(), "latency", latency)
		} else {
			//other warning!
			core.logger.Errorw(path, "method", c.Request.Method, "httpcode", c.Writer.Status(), "errinfo", errinfo,
				"clientip", c.ClientIP(), "latency", latency)
		}
	}

}

func (core *Core) GinRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("---")
				fmt.Println(err)
				fmt.Println("---")
			}
		}()
		c.Next()
	}
}
