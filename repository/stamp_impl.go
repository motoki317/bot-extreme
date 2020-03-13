package repository

func (r *RepositoryImpl) GetStamp(ID string) (*Stamp, error) {
	var stamp Stamp
	if err := r.db.Get(&stamp, "SELECT * FROM `stamp` WHERE `id` = ?", ID); err != nil {
		return nil, err
	}
	return &stamp, nil
}

func (r *RepositoryImpl) UpdateStamp(stamp *Stamp) error {
	_, err := r.db.NamedExec("INSERT INTO `stamp` (id, used) VALUES (:id, :used) ON DUPLICATE KEY UPDATE used = :used", stamp)
	return err
}
