package repository

import "database/sql"

func (r *RepositoryImpl) GetStampRelation(from, to string) (*StampRelation, error) {
	var relation StampRelation
	if err := r.db.Get(&relation, "SELECT * FROM `stamp_relation` WHERE `id_from` = ? AND `id_to` = ?", from, to); err == nil {
		return &relation, nil
	} else if err != sql.ErrNoRows {
		return nil, err
	}
	if err := r.db.Get(&relation, "SELECT * FROM `stamp_relation` WHERE `id_to` = ? AND `id_from` = ?", from, to); err == nil {
		return &relation, nil
	} else if err != sql.ErrNoRows {
		return nil, err
	}
	return nil, nil
}

func (r *RepositoryImpl) GetStampRelations(id string) (relations []*StampRelation, err error) {
	var relationsSlice []StampRelation

	if err = r.db.Get(&relationsSlice, "SELECT * FROM `stamp_relation` WHERE `id_from` = ?", id); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	for _, r := range relationsSlice {
		relation := r
		relations = append(relations, &relation)
	}

	relationsSlice = nil
	if err = r.db.Get(&relationsSlice, "SELECT * FROM `stamp_relation` WHERE `id_to` = ?", id); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	for _, r := range relationsSlice {
		relation := r
		relations = append(relations, &relation)
	}

	return relations, nil
}

func (r *RepositoryImpl) UpdateStampRelation(relation *StampRelation) error {
	_, err := r.db.NamedExec("INSERT INTO `stamp_relation` (id_from, id_to, point) VALUES (:id_from, :id_to, :point) ON DUPLICATE KEY UPDATE point = :point", relation)
	return err
}
