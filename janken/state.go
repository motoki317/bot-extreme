package janken

type State int

const (
	// 相手を待っている（プレイヤーを選ぶかBotと対戦するか）
	OpponentPick State = iota
	// 相手の返信を待っている（プレイヤーと対戦する場合、相手の返信が必要）
	WaitingOpponent
	// プレイヤー vs Botで返信を待っている
	PvB
	// プレイヤー vs プレイヤーで返信を待っている
	PvP
)
