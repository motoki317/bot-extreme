package evaluate

import (
	ast "github.com/magiconair/properties/assert"
	"github.com/motoki317/bot-extreme/repository"
	"math"
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

func TestProcessMessage(t *testing.T) {
	type args struct {
		repo     repository.Repository
		messages []*Message
	}
	type assert struct {
		stampUsed   map[string]int
		effectPoint map[string]float64
		relations   map[string]map[string]float64
	}
	tests := []struct {
		name    string
		args    args
		assert  assert
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				repo: newMockRepository(),
				messages: []*Message{{
					MessageStamps: []*stamp{
						{
							name:        "ultrafastparrot",
							id:          "ultrafastparrot",
							sizeEffect:  "ex-large",
							moveEffects: []string{"rotate", "parrot"},
						},
					},
					UserReactions: [][]*reaction{
						{
							{
								id: "ultrafastparrot",
							},
						},
						{
							{
								id: "ultrafastparrot",
							},
						},
					},
				}},
			},
			assert: assert{
				stampUsed: map[string]int{"ultrafastparrot": 3},
				effectPoint: map[string]float64{
					"small":    0.99,
					"large":    0.99,
					"ex-large": 1.02,
					// move effects
					"inversion": 0.9801,
					"conga":     0.9801,
					// ... and other move effects
					// (1 + 0.01 * 19) * 0.99
					"rotate": 1.1801,
					// 1 * 0.99 + (18 * 0.99 * 0.01) + (1.19 * 0.01)
					"parrot": 1.1781,
				},
				relations: map[string]map[string]float64{},
			},
			wantErr: false,
		},
		{
			name: "empty",
			args: args{
				repo: newMockRepository(),
				messages: []*Message{{
					MessageStamps: nil,
					UserReactions: nil,
				}},
			},
			assert:  assert{},
			wantErr: false,
		},
		{
			name: "relations",
			args: args{
				repo: newMockRepository(),
				messages: []*Message{{
					MessageStamps: []*stamp{
						{
							name:        "bigultrafastparrot_1",
							id:          "bigultrafastparrot_1",
							sizeEffect:  "",
							moveEffects: nil,
						},
						{
							name:        "bigultrafastparrot_2",
							id:          "bigultrafastparrot_2",
							sizeEffect:  "",
							moveEffects: nil,
						},
						{
							name:        "bigultrafastparrot_3",
							id:          "bigultrafastparrot_3",
							sizeEffect:  "",
							moveEffects: nil,
						},
						{
							name:        "bigultrafastparrot_4",
							id:          "bigultrafastparrot_4",
							sizeEffect:  "",
							moveEffects: nil,
						},
					},
					UserReactions: nil,
				}},
			},
			assert: assert{
				stampUsed: map[string]int{
					"bigultrafastparrot_1": 1, "bigultrafastparrot_2": 1, "bigultrafastparrot_3": 1, "bigultrafastparrot_4": 1,
				},
				effectPoint: nil,
				relations: map[string]map[string]float64{
					"bigultrafastparrot_1": {
						"bigultrafastparrot_2": 1,
						"bigultrafastparrot_3": 1,
						"bigultrafastparrot_4": 1,
					},
					"bigultrafastparrot_2": {
						"bigultrafastparrot_1": 1,
						"bigultrafastparrot_3": 1,
						"bigultrafastparrot_4": 1,
					},
					"bigultrafastparrot_3": {
						"bigultrafastparrot_1": 1,
						"bigultrafastparrot_2": 1,
						"bigultrafastparrot_4": 1,
					},
					"bigultrafastparrot_4": {
						"bigultrafastparrot_1": 1,
						"bigultrafastparrot_2": 1,
						"bigultrafastparrot_3": 1,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "minus_relation",
			args: args{
				repo: newMockRepository(),
				messages: []*Message{
					{
						MessageStamps: []*stamp{
							{
								name:        "bigultrafastparrot_1",
								id:          "bigultrafastparrot_1",
								sizeEffect:  "",
								moveEffects: nil,
							},
							{
								name:        "bigultrafastparrot_2",
								id:          "bigultrafastparrot_2",
								sizeEffect:  "",
								moveEffects: nil,
							},
						},
						UserReactions: nil,
					},
					{
						MessageStamps: []*stamp{
							{
								name:        "bigultrafastparrot_1",
								id:          "bigultrafastparrot_1",
								sizeEffect:  "",
								moveEffects: nil,
							},
						},
						UserReactions: nil,
					},
				},
			},
			assert: assert{
				stampUsed:   nil,
				effectPoint: nil,
				relations: map[string]map[string]float64{
					"bigultrafastparrot_1": {
						"bigultrafastparrot_2": 0.9,
					},
					"bigultrafastparrot_2": {
						"bigultrafastparrot_1": 0.9,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, message := range tt.args.messages {
				if err := ProcessMessage(tt.args.repo, message); (err != nil) != tt.wantErr {
					t.Errorf("processMessage() error = %v, wantErr %v", err, tt.wantErr)
				}
			}

			for stamp, used := range tt.assert.stampUsed {
				stampInfo, _ := tt.args.repo.GetStamp(stamp)
				if stampInfo == nil {
					t.Fatalf("got nil stamp info for %s", stamp)
				}
				ast.Equal(t, stampInfo.Used, used)
			}
			for eff, p := range tt.assert.effectPoint {
				effectInfo, _ := tt.args.repo.GetEffectPoint(eff)
				if effectInfo == nil {
					t.Fatalf("got nil effect info for %s", eff)
				}
				floatEquals(t, effectInfo.Point, p, 1e-6)
			}
			for from, relation := range tt.assert.relations {
				for to, p := range relation {
					relationInfo, _ := tt.args.repo.GetStampRelation(from, to)
					if relationInfo == nil {
						t.Fatalf("got nil relation for %s and %s", from, to)
					}
					floatEquals(t, relationInfo.Point, p, 1e-6)
				}
			}
		})
	}
}

func floatEquals(t *testing.T, got float64, want float64, precision float64) {
	if math.Abs(got-want) > precision {
		t.Fatalf("got %v, want %v at precision %v", got, want, precision)
	}
}
