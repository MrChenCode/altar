package controller

import (
	"altar/application/context/cctx"
	"github.com/gin-gonic/gin"
)

type Game struct{}

func (_ *Game) GameInfo(c *cctx.ControllerContext) {
	c.Log.Info("controller_gameid", 9000, "controller_gamename", "魔兽争霸")
	c.Log.Error("http_get_error", "无效的游戏id")
	c.Log.Info("controller_getgame", 1)
	c.Redis.Get("altar_redis_test")
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "",
		"result": map[string]interface{}{
			"gameid":   100,
			"gamename": c.Model.Game.GetGameInfo(),
		},
	})
	c.Library.Func.GetGame()
	c.Model.Game.GetGameInfo()
}