package handler

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"github.com/motoki317/bot-extreme/repository"
	"log"
	"sync"
	"testing"
	"time"
)

// レーティング、スタンプの使用回数、エフェクト、関係、最後に見たチャンネルを擬似的に保存するrepository
type MockRepository struct {
	lock         sync.Mutex
	channelLock  sync.Mutex
	rating       map[string]float64
	used         map[string]int
	effects      map[string]float64
	relations    map[string]map[string]float64
	seenChannels map[string]time.Time
}

func newMockRepository() *MockRepository {
	return &MockRepository{
		lock:         sync.Mutex{},
		rating:       make(map[string]float64),
		used:         make(map[string]int),
		effects:      make(map[string]float64),
		relations:    make(map[string]map[string]float64),
		seenChannels: make(map[string]time.Time),
	}
}

func (m *MockRepository) Lock() {
	m.lock.Lock()
}

func (m *MockRepository) Unlock() {
	m.lock.Unlock()
}

func (m *MockRepository) ChannelLock() {
	m.channelLock.Lock()
}

func (m *MockRepository) ChannelUnlock() {
	m.channelLock.Unlock()
}

func (m *MockRepository) GetRating(ID string) (*repository.Rating, error) {
	if r, ok := m.rating[ID]; ok {
		return &repository.Rating{
			ID:     ID,
			Rating: r,
		}, nil
	} else {
		return nil, nil
	}
}

func (m *MockRepository) UpdateRating(r *repository.Rating) error {
	m.rating[r.ID] = r.Rating
	return nil
}

func (m *MockRepository) GetEffectPoint(name string) (*repository.EffectPoint, error) {
	if e, ok := m.effects[name]; ok {
		return &repository.EffectPoint{
			Name:  name,
			Point: e,
		}, nil
	} else {
		return nil, nil
	}
}

func (m *MockRepository) GetAllEffectPoints() ([]*repository.EffectPoint, error) {
	ret := make([]*repository.EffectPoint, 0, len(m.effects))
	for name, p := range m.effects {
		ret = append(ret, &repository.EffectPoint{
			Name:  name,
			Point: p,
		})
	}
	return ret, nil
}

func (m *MockRepository) UpdateEffectPoint(point *repository.EffectPoint) error {
	m.effects[point.Name] = point.Point
	return nil
}

func (m *MockRepository) UpdateAllEffectPoints(points []*repository.EffectPoint) error {
	for _, point := range points {
		m.effects[point.Name] = point.Point
	}
	return nil
}

func (m *MockRepository) GetStampRelation(from, to string) (*repository.StampRelation, error) {
	if r, ok := m.relations[from]; ok {
		if ret, ok := r[to]; ok {
			return &repository.StampRelation{
				IDFrom: from,
				IDTo:   to,
				Point:  ret,
			}, nil
		}
	}
	if r, ok := m.relations[to]; ok {
		if ret, ok := r[from]; ok {
			return &repository.StampRelation{
				IDFrom: to,
				IDTo:   from,
				Point:  ret,
			}, nil
		}
	}
	return nil, nil
}

func (m *MockRepository) GetStampRelations(id string) ([]*repository.StampRelation, error) {
	ret := make([]*repository.StampRelation, 0)
	if r, ok := m.relations[id]; ok {
		for to, p := range r {
			ret = append(ret, &repository.StampRelation{
				IDFrom: id,
				IDTo:   to,
				Point:  p,
			})
		}
	}
	for from, r := range m.relations {
		if p, ok := r[id]; ok {
			ret = append(ret, &repository.StampRelation{
				IDFrom: from,
				IDTo:   id,
				Point:  p,
			})
		}
	}
	return ret, nil
}

func (m *MockRepository) UpdateStampRelation(relation *repository.StampRelation) error {
	if _, ok := m.relations[relation.IDFrom]; !ok {
		m.relations[relation.IDFrom] = make(map[string]float64)
	}
	m.relations[relation.IDFrom][relation.IDTo] = relation.Point
	return nil
}

func (m *MockRepository) UpdateStampRelations(relations []*repository.StampRelation) error {
	for _, relation := range relations {
		if _, ok := m.relations[relation.IDFrom]; !ok {
			m.relations[relation.IDFrom] = make(map[string]float64)
		}
		m.relations[relation.IDFrom][relation.IDTo] = relation.Point
	}
	return nil
}

func (m *MockRepository) GetStamp(ID string) (*repository.Stamp, error) {
	if stamp, ok := m.used[ID]; ok {
		return &repository.Stamp{
			ID:   ID,
			Used: stamp,
		}, nil
	}
	return nil, nil
}

func (m *MockRepository) UpdateStamp(stamp *repository.Stamp) error {
	log.Printf("Update stamp %s to %v", stamp.ID, stamp.Used)
	m.used[stamp.ID] = stamp.Used
	return nil
}

func (m *MockRepository) GetSeenChannel(ID string) (*repository.SeenChannel, error) {
	if t, ok := m.seenChannels[ID]; ok {
		return &repository.SeenChannel{
			ID:                   ID,
			LastProcessedMessage: t,
		}, nil
	}
	return nil, nil
}

func (m *MockRepository) UpdateSeenChannel(channel *repository.SeenChannel) error {
	m.seenChannels[channel.ID] = channel.LastProcessedMessage
	return nil
}

func Test_updater_updateRatings(t *testing.T) {
	type fields struct {
		repo repository.Repository
	}
	type args struct {
		channelID string
		from      time.Time
	}
	type ast struct {
		stampsUsed map[string]int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		ast     ast
	}{
		{
			name:   "kashiwade",
			fields: fields{newMockRepository()},
			args: args{
				channelID: "ec513c2c-b105-4a0e-8ccd-86b12ef307d8",
				from:      time.Date(2020, time.March, 1, 0, 0, 0, 0, time.UTC),
			},
			wantErr: false,
			ast: ast{
				stampsUsed: map[string]int{
					// ultrafastparrot
					"4a7ac270-0bfa-4b2a-9ebc-58e3487a23da": 1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &updater{
				repo: tt.fields.repo,
			}
			err := u.repo.UpdateSeenChannel(&repository.SeenChannel{
				ID:                   tt.args.channelID,
				LastProcessedMessage: tt.args.from,
			})
			if err != nil {
				t.Errorf("want err = nil, got = %v", err)
			}
			if err := u.updateRatings(tt.args.channelID); (err != nil) != tt.wantErr {
				t.Errorf("updateRatings() error = %v, wantErr %v", err, tt.wantErr)
			}

			for stampID, count := range tt.ast.stampsUsed {
				stampInfo, _ := u.repo.GetStamp(stampID)
				if stampInfo == nil {
					t.Errorf("got stampInfo nil for stamp %s, want non-nil", stampID)
					return
				}
				assert.Equal(t, stampInfo.Used >= count, true, fmt.Sprintf("want stamp %s to be used at least %v times", stampID, count))
			}
		})
	}
}
