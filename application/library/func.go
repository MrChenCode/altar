package library

type Func struct {
	*BasicLibrary
}

func (f *Func) GetGame() {
	f.Log.Info("func_gameid", true, "func_gamename", "斗地主")
	f.Log.Error("func_getgameinfo", "timeout", "func_response", "dial error")
}
