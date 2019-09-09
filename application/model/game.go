package model

type GameModel struct {
	*modelContext
}

func (g *GameModel) GetGameInfo() string {
	g.ctx.Log.Info("model_gameid", 1001, "model_gamename", "三国战纪")
	g.ctx.Log.Info("model_userinfo", "6666", "model_gameuser", 7777)
	g.ctx.Log.Error("model_gameinfo", "null", "model_gameid", 1001, "model_gamerequest", "timeout")

	g.library.Func.GetGame()

	return "三国战纪"
}
