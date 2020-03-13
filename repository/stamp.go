package repository

type Stamp struct {
	ID   string `db:"id"`
	Used int    `db:"used"`
}

type StampRepository interface {
	// スタンプの情報を取得します
	// 存在しない場合はnilを返します
	GetStamp(ID string) (*Stamp, error)
	// スタンプの情報を更新します
	// 存在する場合は更新、
	// 存在しない場合は新規に作成します
	UpdateStamp(stamp *Stamp) error
}
