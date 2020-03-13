package janken

type State int

const (
	// 相手を待っている（プレイヤーを選ぶかBotと対戦するか）
	WaitingOpponent State = iota
	// プレイヤー vs Botで返信を待っている
	PvB
	// プレイヤー vs プレイヤーで返信を待っている
	PvP
)
