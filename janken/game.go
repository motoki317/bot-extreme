package janken

import (
	"os"
)

var (
	// "BOT_extreme"
	botName = os.Getenv("BOT_NAME")
	botUuid = os.Getenv("BOT_UUID")
)

type Game struct {
	State
	// 自分の情報
	self *User
	// 対戦相手 Botと対戦しているなら空
	opponent *User
	// 自分と対戦相手の手
	selfResponse     string
	opponentResponse string
}

type User struct {
	Name string
	ID   string
}

func newGame(sender *User) *Game {
	return &Game{
		State: OpponentPick,
		self:  sender,
	}
}
