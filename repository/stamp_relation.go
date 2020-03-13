package repository

// from and to are interchangeable
type StampRelation struct {
	IDFrom string  `db:"id_from"`
	IDTo   string  `db:"id_to"`
	Point  float64 `db:"point"`
}

type StampRelationRepository interface {
	// fromとto（スタンプID）のスタンプの関係を取得します
	// fromとtoが入れ替わる場合があります
	// 存在しない場合はnilを返します
	GetStampRelation(from, to string) (*StampRelation, error)
	// スタンプIDの関係を全て取得します
	GetStampRelations(id string) ([]*StampRelation, error)
	// スタンプの関係を更新します
	// 存在する場合は更新
	// 存在しない場合は新規に作成します
	UpdateStampRelation(relation *StampRelation) error
}
