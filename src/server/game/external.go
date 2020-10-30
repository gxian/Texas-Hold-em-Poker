package game

import (
	"texas-holdem/server/game/internal"
)

var (
	Module  = new(internal.Module) //建立模块新的
	ChanRPC = internal.ChanRPC
)
