package repository

import "database/sql"

func (r *RepositoryImpl) GetSeenChannel(ID string) (*SeenChannel, error) {
	var channel SeenChannel
	if err := r.db.Get(&channel, "SELECT * FROM `seen_channel` WHERE `id` = ?", ID); err != nil && err != sql.ErrNoRows {
		return nil, err
	} else if err == sql.ErrNoRows {
		return nil, nil
	}
	return &channel, nil
}

func (r *RepositoryImpl) UpdateSeenChannel(channel *SeenChannel) error {
	_, err := r.db.NamedExec(
		"INSERT INTO `seen_channel` (id, last_processed_message) VALUES (:id, :last_processed_message) ON DUPLICATE KEY UPDATE `last_processed_message` = :last_processed_message",
		channel)
	return err
}
