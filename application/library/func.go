package library

type Func struct {
	*libraryContext
}

func (f *Func) GetGame() {
	f.ctx.Log.Info("func_gameid", true, "func_gamename", "斗地主")
	f.ctx.Log.Error("func_getgameinfo", "timeout", "func_response", "dial error")
}
