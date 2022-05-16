package evaluate

import (
	"testing"

	"github.com/motoki317/bot-extreme/repository"
)

// なにもないrepository
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

func (m EmptyRepository) GetAllRatings() ([]*repository.Rating, error) {
	return nil, nil
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

func (m EmptyRepository) UpdateAllEffectPoints(points []*repository.EffectPoint) error {
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

func (m EmptyRepository) DeleteStampRelations(threshold float64) error {
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

func TestMessagePoint(t *testing.T) {
	type args struct {
		repo    repository.Repository
		content string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "parrot",
			args: args{
				repo:    EmptyRepository{},
				content: ":ultrafastparrot.ex-large.rotate.parrot:",
			},
			wantErr: false,
		},
		{
			name: "thonk_spin",
			args: args{
				repo:    EmptyRepository{},
				content: ":thonk_spin.ex-large.rotate.parrot:",
			},
			wantErr: false,
		},
		{
			name: "multi",
			args: args{
				repo:    EmptyRepository{},
				content: ":ranpuro_5::oisu-4yoko::ranpuro_1::ranpuro_3::ranpuro_4::ranpuro_2::ranpuro_4::ranpuro_2:",
			},
			wantErr: false,
		},
		{
			name: "with some message",
			args: args{
				repo:    EmptyRepository{},
				content: "@BOT_extreme :ranpuro_5::oisu-4yoko::ranpuro_1::ranpuro_3::ranpuro_4::ranpuro_2::ranpuro_4::ranpuro_2:",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPts, err := MessagePoint(tt.args.repo, tt.args.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("MessagePoint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !(0 <= gotPts && gotPts < 65) {
				t.Errorf("MessagePoint() gotPts = %v, want %v", gotPts, "0 <= pts < 65")
			}
			t.Logf("Got %v pts for %s", gotPts, tt.args.content)
		})
	}
}
