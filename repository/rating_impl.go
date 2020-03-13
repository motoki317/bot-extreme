package repository

func (r RepositoryImpl) GetRating(ID string) (*Rating, error) {
	var rating Rating
	if err := r.db.Get(&rating, "SELECT * FROM `rating` WHERE `id` = ?", ID); err != nil {
		return nil, err
	}
	return &rating, nil
}

func (r RepositoryImpl) UpdateRating(rating *Rating) error {
	_, err := r.db.NamedExec(
		"INSERT INTO `rating` (`id`, `rating`) VALUES (:id, :rating) ON DUPLICATE KEY UPDATE `rating` = :rating",
		rating)
	return err
}
