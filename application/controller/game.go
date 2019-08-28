package controller

import (
	"altar/application/core/basic"
	"github.com/gin-gonic/gin"
)

type Game struct{}

func (_ *Game) GameInfo(ctx *basic.Controller) {
	ctx.Log.Info("controller_gameid", 9000, "controller_gamename", "魔兽争霸")
	ctx.Log.Error("http_get_error", "无效的游戏id")
	ctx.Log.Info("controller_getgame", 1)
	ctx.JSON(200, gin.H{
		"code": 0,
		"msg":  ctx.Model.Game.GetGameInfo(),
	})
}
