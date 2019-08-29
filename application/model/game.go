package model

type GameModel struct {
	*Model
}

func (g *GameModel) GetGameInfo() string {
	g.Log.Info("model_gameid", 1001, "model_gamename", "三国战纪")
	g.Log.Info("model_userinfo", "6666", "model_gameuser", 7777)
	g.Log.Error("model_gameinfo", "null", "model_gameid", 1001, "model_gamerequest", "timeout")

	g.Library.Func.GetGame()

	return "三国战纪"
}
