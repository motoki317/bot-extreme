package repository

import (
	"database/sql"
	"strings"
)

func (r *RepositoryImpl) GetEffectPoint(name string) (*EffectPoint, error) {
	var effectPoint EffectPoint
	if err := r.db.Get(&effectPoint, "SELECT * FROM `effect_point` WHERE `name` = ?", name); err != nil && err != sql.ErrNoRows {
		return nil, err
	} else if err == sql.ErrNoRows {
		return nil, nil
	}
	return &effectPoint, nil
}

func (r *RepositoryImpl) GetAllEffectPoints() ([]*EffectPoint, error) {
	var effectPoints []EffectPoint
	if err := r.db.Select(&effectPoints, "SELECT * FROM `effect_point`"); err != nil && err != sql.ErrNoRows {
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
	_, err := r.db.NamedExec("INSERT INTO `effect_point` (name, point) VALUES (:name, :point) ON DUPLICATE KEY UPDATE point = VALUES(point)", point)
	return err
}

func (r *RepositoryImpl) UpdateAllEffectPoints(points []*EffectPoint) error {
	if len(points) == 0 {
		return nil
	}

	placeHolders := make([]string, 0, len(points))
	for range points {
		placeHolders = append(placeHolders, "(?, ?)")
	}
	placeHolder := strings.Join(placeHolders, ", ")

	values := make([]interface{}, 0, len(points)*2)
	for _, p := range points {
		values = append(values, p.Name)
		values = append(values, p.Point)
	}

	_, err := r.db.Exec("INSERT INTO `effect_point` (name, point) VALUES "+placeHolder+" ON DUPLICATE KEY UPDATE point = VALUES(point)", values...)
	return err
}
