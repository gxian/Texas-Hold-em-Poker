package gate

import (
	"texas-holdem/server/game"
	"texas-holdem/server/login"
	"texas-holdem/server/protocol"
)

func init() {
	protocol.Processor.SetRouter(&protocol.UserLoginInfo{}, login.ChanRPC)
	protocol.Processor.SetRouter(&protocol.Version{}, login.ChanRPC)
	protocol.Processor.SetRouter(&protocol.RoomList{}, login.ChanRPC)

	protocol.Processor.SetRouter(&protocol.JoinRoom{}, game.ChanRPC)
	protocol.Processor.SetRouter(&protocol.LeaveRoom{}, game.ChanRPC)
	protocol.Processor.SetRouter(&protocol.SitDown{}, game.ChanRPC)
	protocol.Processor.SetRouter(&protocol.StandUp{}, game.ChanRPC)
	protocol.Processor.SetRouter(&protocol.Bet{}, game.ChanRPC)
	protocol.Processor.SetRouter(&protocol.Chat{}, game.ChanRPC)
}
