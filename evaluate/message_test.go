package evaluate

import (
	"github.com/motoki317/bot-extreme/repository"
	"testing"
)

type MockRepository struct {
}

func (m MockRepository) GetRating(ID string) (*repository.Rating, error) {
	return &repository.Rating{
		ID:     ID,
		Rating: 1500,
	}, nil
}

func (m MockRepository) UpdateRating(*repository.Rating) error {
	return nil
}

func (m MockRepository) GetEffectPoint(name string) (*repository.EffectPoint, error) {
	return &repository.EffectPoint{
		Name:  name,
		Point: 1,
	}, nil
}

func (m MockRepository) GetAllEffectPoints() ([]*repository.EffectPoint, error) {
	return []*repository.EffectPoint{}, nil
}

func (m MockRepository) UpdateEffectPoint(point *repository.EffectPoint) error {
	return nil
}

func (m MockRepository) GetStampRelation(from, to string) (*repository.StampRelation, error) {
	return &repository.StampRelation{
		IDFrom: from,
		IDTo:   to,
		Point:  1,
	}, nil
}

func (m MockRepository) GetStampRelations(id string) ([]*repository.StampRelation, error) {
	return []*repository.StampRelation{}, nil
}

func (m MockRepository) UpdateStampRelation(relation *repository.StampRelation) error {
	return nil
}

func (m MockRepository) GetStamp(ID string) (*repository.Stamp, error) {
	return &repository.Stamp{
		ID:   ID,
		Used: 5,
	}, nil
}

func (m MockRepository) UpdateStamp(stamp *repository.Stamp) error {
	return nil
}

func (m MockRepository) GetSeenChannel(ID string) (*repository.SeenChannel, error) {
	return nil, nil
}

func (m MockRepository) UpdateSeenChannel(channel *repository.SeenChannel) error {
	return nil
}

func TestMessage(t *testing.T) {
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
				repo:    MockRepository{},
				content: ":ultrafastparrot.ex-large.rotate.parrot:",
			},
			wantErr: false,
		},
		{
			name: "thonk_spin",
			args: args{
				repo:    MockRepository{},
				content: ":thonk_spin.ex-large.rotate.parrot:",
			},
			wantErr: false,
		},
		{
			name: "multi",
			args: args{
				repo:    MockRepository{},
				content: ":ranpuro_5::oisu-4yoko::ranpuro_1::ranpuro_3::ranpuro_4::ranpuro_2::ranpuro_4::ranpuro_2:",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPts, err := Message(tt.args.repo, tt.args.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("Message() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !(0 <= gotPts && gotPts < 65) {
				t.Errorf("Message() gotPts = %v, want %v", gotPts, "0 <= pts < 65")
			}
			t.Logf("Got %v pts for %s", gotPts, tt.args.content)
		})
	}
}
