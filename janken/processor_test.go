package janken

import (
	"github.com/magiconair/properties/assert"
	"github.com/motoki317/bot-extreme/repository"
	"testing"
)

type EmptyRepository struct {
	// ユーザーレーティング
	repository.RatingRepository
	// スタンプエフェクトのポイント
	repository.EffectPointRepository
	// スタンプ同士の案系
	repository.StampRelationRepository
	// スタンプ単体のポイント
	repository.StampRepository
	// 処理したチャンネルの保存用
	repository.SeenChannelRepository
}

func (m EmptyRepository) Lock() {}

func (m EmptyRepository) Unlock() {}

func (m EmptyRepository) ChannelLock() {}

func (m EmptyRepository) ChannelUnlock() {}

func TestJankenProcessor(t *testing.T) {
	p := NewProcessor(EmptyRepository{})

	t.Run("player versus player", func(t *testing.T) {
		sender := &User{
			Name: "toki",
			ID:   "this_is_totally_a_uuid",
		}
		opponent := &User{
			Name: "xxpoxx",
			ID:   "also_totally_a_uuid",
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

		err = p.handle(opponent, "@BOT_extreme @xxpoxx", []*User{opponent}, respond)
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
			Name: "toki",
			ID:   "this_is_totally_a_uuid",
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
