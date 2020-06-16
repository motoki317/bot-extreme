package repository

import "database/sql"

func (r *RepositoryImpl) GetRating(ID string) (*Rating, error) {
	var rating Rating
	if err := r.db.Get(&rating, "SELECT * FROM `rating` WHERE `id` = ?", ID); err != nil && err != sql.ErrNoRows {
		return nil, err
	} else if err == sql.ErrNoRows {
		return nil, nil
	}
	return &rating, nil
}

func (r *RepositoryImpl) GetAllRatings() ([]*Rating, error) {
	var ratings []Rating
	if err := r.db.Select(&ratings, "SELECT * FROM `rating`"); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	ret := make([]*Rating, 0, len(ratings))
	for _, e := range ratings {
		rating := e
		ret = append(ret, &rating)
	}
	return ret, nil
}

func (r *RepositoryImpl) UpdateRating(rating *Rating) error {
	_, err := r.db.NamedExec(
		"INSERT INTO `rating` (`id`, `rating`) VALUES (:id, :rating) ON DUPLICATE KEY UPDATE `rating` = :rating",
		rating)
	return err
}
