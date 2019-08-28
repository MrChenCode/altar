package model

import "altar/application/core/context"

type GameModel struct {
	*context.RequestContext
}

func (g *GameModel) GetGameInfo() string {
	g.Log.Info("model_gameid", 1001, "model_gamename", "三国战纪")
	g.Log.Info("model_userinfo", "6666", "model_gameuser", 7777)
	g.Log.Error("model_gameinfo", "null", "model_gameid", 1001, "model_gamerequest", "timeout")
	return "三国战纪"
}
