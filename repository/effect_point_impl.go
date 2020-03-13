package repository

func (r *RepositoryImpl) GetEffectPoint(name string) (*EffectPoint, error) {
	var effectPoint EffectPoint
	if err := r.db.Get(&effectPoint, "SELECT * FROM `effect_point` WHERE `name` = ?", name); err != nil {
		return nil, err
	}
	return &effectPoint, nil
}

func (r *RepositoryImpl) GetAllEffectPoints() ([]*EffectPoint, error) {
	var effectPoints []EffectPoint
	if err := r.db.Get(&effectPoints, "SELECT * FROM `effect_point`"); err != nil {
		return nil, err
	}
	ret := make([]*EffectPoint, 0, len(effectPoints))
	for _, e := range effectPoints {
		effect := e
		ret = append(ret, &effect)
	}
	return ret, nil
}

func (r *RepositoryImpl) UpdateEffectPoint(point *EffectPoint) error {
	_, err := r.db.NamedExec("INSERT INTO `effect_point` (name, point) VALUES (:name, :point) ON DUPLICATE KEY UPDATE point = :point", point)
	return err
}
