package repository

import "time"

type SeenChannel struct {
	ID                   string    `db:"id"`
	LastProcessedMessage time.Time `db:"last_processed_message"`
}

type SeenChannelRepository interface {
	// チャンネルIDについて、そのチャンネルで最後に処理したメッセージの日時を取得します
	// 存在しない場合はnilを返します
	GetSeenChannel(ID string) (*SeenChannel, error)
	// チャンネルを処理した情報を更新します
	// 存在する場合は更新
	// 存在しない場合は新規に作成します
	UpdateSeenChannel(channel *SeenChannel) error
}
