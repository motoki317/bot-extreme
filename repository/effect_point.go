package repository

type EffectPoint struct {
	Name  string  `db:"name"`
	Point float64 `db:"point"`
}

type EffectPointRepository interface {
	// エフェクトのポイントを取得します
	// 存在しない場合にはnilを返します
	GetEffectPoint(name string) (*EffectPoint, error)
	// 全てのレコードを取得し返します
	GetAllEffectPoints() ([]*EffectPoint, error)
	// エフェクトのポイントを更新します
	// 存在しない場合は新規に作成、
	// 既に存在する場合は更新します。
	UpdateEffectPoint(point *EffectPoint) error
	// 全てのエフェクトのポイントを更新します
	// 存在しない場合は新規に作成、
	// 既に存在する場合は更新します。
	UpdateAllEffectPoints(points []*EffectPoint) error
}
