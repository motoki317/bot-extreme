package repository

type Rating struct {
	ID     string  `db:"id"`
	Rating float64 `db:"rating"`
}

type RatingRepository interface {
	// 該当ユーザーのRatingを取得します。
	// 存在しない場合、nilを返します。
	GetRating(ID string) (*Rating, error)
	// Ratingを更新します。
	// 存在しない場合は新規に作成、
	// 既に存在する場合は更新します。
	UpdateRating(*Rating) error
}
