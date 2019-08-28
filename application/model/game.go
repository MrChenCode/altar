package model

import "altar/application/core/context"

type GameModel struct {
	*context.RequestContext
}

func (g *GameModel) GetGameInfo() string {
	g.Log.Info("gameid", 101, "gamename", "斗地主")
	return "斗地主"
}
