package janken

import (
	"github.com/magiconair/properties/assert"
	"github.com/motoki317/bot-extreme/repository"
	"testing"
)

type EmptyRepository struct{}

func (m EmptyRepository) Lock() {}

func (m EmptyRepository) Unlock() {}

func (m EmptyRepository) ChannelLock() {}

func (m EmptyRepository) ChannelUnlock() {}

func (m EmptyRepository) GetRating(ID string) (*repository.Rating, error) {
	return &repository.Rating{
		ID:     ID,
		Rating: 1500,
	}, nil
}

func (m EmptyRepository) UpdateRating(*repository.Rating) error {
	return nil
}

func (m EmptyRepository) GetEffectPoint(name string) (*repository.EffectPoint, error) {
	return &repository.EffectPoint{
		Name:  name,
		Point: 1,
	}, nil
}

func (m EmptyRepository) GetAllEffectPoints() ([]*repository.EffectPoint, error) {
	return []*repository.EffectPoint{}, nil
}

func (m EmptyRepository) UpdateEffectPoint(point *repository.EffectPoint) error {
	return nil
}

func (m EmptyRepository) GetStampRelation(from, to string) (*repository.StampRelation, error) {
	return &repository.StampRelation{
		IDFrom: from,
		IDTo:   to,
		Point:  1,
	}, nil
}

func (m EmptyRepository) GetStampRelations(id string) ([]*repository.StampRelation, error) {
	return []*repository.StampRelation{}, nil
}

func (m EmptyRepository) UpdateStampRelation(relation *repository.StampRelation) error {
	return nil
}

func (m EmptyRepository) UpdateStampRelations(relations []*repository.StampRelation) error {
	return nil
}

func (m EmptyRepository) UpdateAllEffectPoints(points []*repository.EffectPoint) error {
	return nil
}

func (m EmptyRepository) GetStamp(ID string) (*repository.Stamp, error) {
	return &repository.Stamp{
		ID:   ID,
		Used: 5,
	}, nil
}

func (m EmptyRepository) UpdateStamp(stamp *repository.Stamp) error {
	return nil
}

func (m EmptyRepository) GetSeenChannel(ID string) (*repository.SeenChannel, error) {
	return nil, nil
}

func (m EmptyRepository) UpdateSeenChannel(channel *repository.SeenChannel) error {
	return nil
}

func TestJankenProcessor(t *testing.T) {
	p := NewProcessor(EmptyRepository{})

	t.Run("player versus player", func(t *testing.T) {
		sender := &User{
			DisplayName: "toki",
			ID:          "this_is_totally_a_uuid",
		}
		opponent := &User{
			DisplayName: "xxpoxx",
			ID:          "also_totally_a_uuid",
		}
		respond := func(s string) {
			t.Log("Got response from processor: " + s)
		}

		err := p.handle(sender, "@BOT_extreme じゃんけんしよう", []*User{}, respond)
		assert.Equal(t, err, nil)
		assert.Equal(t, len(p.games), 1)

		err = p.handle(sender, "@BOT_extreme @xxpoxx", []*User{opponent}, respond)
		assert.Equal(t, err, nil)
		assert.Equal(t, len(p.games), 1)

		err = p.handle(sender, "@BOT_extreme :thonk_spin.ex-large.rotate.parrot:", []*User{}, respond)
		assert.Equal(t, err, nil)
		assert.Equal(t, len(p.games), 1)

		err = p.handle(opponent, "@BOT_extreme :ranpuro_5::oisu-4yoko::ranpuro_1::ranpuro_3::ranpuro_4::ranpuro_2::ranpuro_4::ranpuro_2:", []*User{}, respond)
		assert.Equal(t, err, nil)
		assert.Equal(t, len(p.games), 0)
	})

	t.Run("player versus bot", func(t *testing.T) {
		sender := &User{
			DisplayName: "toki",
			ID:          "this_is_totally_a_uuid",
		}
		respond := func(s string) {
			t.Log("Got response from processor: " + s)
		}

		err := p.handle(sender, "@BOT_extreme じゃんけんしよう", []*User{}, respond)
		assert.Equal(t, err, nil)
		assert.Equal(t, len(p.games), 1)

		err = p.handle(sender, "@BOT_extreme ひとりで", []*User{}, respond)
		assert.Equal(t, err, nil)
		assert.Equal(t, len(p.games), 1)

		err = p.handle(sender, "@BOT_extreme :thonk_spin.ex-large.rotate.parrot:", []*User{}, respond)
		assert.Equal(t, err, nil)
		assert.Equal(t, len(p.games), 0)
	})
}
