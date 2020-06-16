package evaluate

import (
	"github.com/motoki317/bot-extreme/repository"
	"math"
	"sync"
)

var (
	// stamps cache - stamp name to stamp uuid
	// ユーザースタンプも含む（ユーザーID）
	stampsMapLock sync.RWMutex
	stampsMap     map[string]string
	sizes         []string
	moves         []string
	sizeEffects   map[string]bool
	moveEffects   map[string]bool
)

type stamp struct {
	name        string
	id          string
	sizeEffect  string
	moveEffects []string
}

type messageEvaluator struct {
	repo   repository.Repository
	stamps []*stamp
}

func init() {
	sizes = []string{
		"ex-large", "large", "small",
	}
	moves = []string{
		"rotate",
		"rotate-inv",
		"wiggle",
		"parrot",
		"zoom",
		"inversion",
		"turn",
		"turn-v",
		"happa",
		"pyon",
		"flashy",
		"pull",
		"atsumori",
		"stretch",
		"stretch-v",
		// doesn't consider that
		// conga == marquee
		"conga",
		"marquee",
		// conga-inv == marquee-inv
		"conga-inv",
		"marquee-inv",
		"rainbow",
		"ascension",
		"shake",
		"party",
		"attract",
	}

	sizeEffects = make(map[string]bool)
	moveEffects = make(map[string]bool)

	for _, s := range sizes {
		sizeEffects[s] = true
	}
	for _, m := range moves {
		moveEffects[m] = true
	}

	stampsMap = make(map[string]string)
}

// stampsMapに新しいスタンプを追加します
func AddStamp(stampName, stampID string) {
	stampsMap[stampName] = stampID
}

// じゃんけんの手を評価し、点数を返します。
func MessagePoint(repo repository.Repository, content string) (pts float64, err error) {
	// parse stamps
	stamps := getMessageStamps(content)

	if containsUnknownStamp(stamps) {
		if err := reCacheStamps(); err != nil {
			return 0, err
		}
	}

	// check non existent stamps
	filteredStamps := make([]*stamp, 0)
	stampsMapLock.RLock()
	for _, s := range stamps {
		if stampId, ok := stampsMap[s.name]; ok {
			s.id = stampId
			filteredStamps = append(filteredStamps, s)
		}
	}
	stampsMapLock.RUnlock()

	eval := &messageEvaluator{
		repo:   repo,
		stamps: filteredStamps,
	}

	return eval.calculatePoints()
}

func (e *messageEvaluator) calculatePoints() (float64, error) {
	if len(e.stamps) == 0 {
		return 0, nil
	}

	effectPoints := make(map[string]float64)
	if points, err := e.repo.GetAllEffectPoints(); err == nil {
		for _, p := range points {
			effectPoints[p.Name] = p.Point
		}
	} else {
		return 0, err
	}

	var stampPoint float64
	var effectPoint float64
	for _, s := range e.stamps {
		p, err := e.stampPoint(s)
		if err != nil {
			return 0, err
		}
		stampPoint += p
		p = e.effectPoint(s, effectPoints)
		effectPoint += p
	}
	// take average of all stamps
	stampPoint /= float64(len(e.stamps))
	effectPoint /= float64(len(e.stamps))
	combinationPoint, err := e.combinationPoint()
	if err != nil {
		return 0, err
	}
	return stampPoint + effectPoint + combinationPoint, nil
}

func (e *messageEvaluator) stampPoint(stamp *stamp) (float64, error) {
	used, err := e.repo.GetStamp(stamp.id)
	if err != nil {
		return 0, err
	}
	if used == nil {
		used = &repository.Stamp{
			ID:   stamp.id,
			Used: 0,
		}
	}

	tmp := math.Log(float64(used.Used + 2))
	p := 24 - 2*(tmp/5+5/tmp)
	return math.Min(20, math.Max(0, p)), nil
}

func (e *messageEvaluator) effectPoint(stamp *stamp, effectPoints map[string]float64) float64 {
	var sizeEffectPoint float64
	if stamp.sizeEffect != "" {
		if p, ok := effectPoints[stamp.sizeEffect]; ok {
			sizeEffectPoint += p
		}
	}
	var moveEffectPoint float64
	for _, e := range stamp.moveEffects {
		if p, ok := effectPoints[e]; ok {
			moveEffectPoint += p
		}
	}
	moveEffectPoint /= math.Max(1, float64(len(stamp.moveEffects)))

	sizeEffectPoint = math.Min(5, math.Max(0, 6-(sizeEffectPoint/1.5+1.5/sizeEffectPoint)/2))
	moveEffectPoint = math.Min(20, math.Max(0, 24-2*(moveEffectPoint/3+3/moveEffectPoint)))
	return sizeEffectPoint + moveEffectPoint
}

func (e *messageEvaluator) combinationPoint() (float64, error) {
	stampIdsMap := make(map[string]bool)
	for _, s := range e.stamps {
		stampIdsMap[s.id] = true
	}

	// id to id to point map
	stampRelations := make(map[string]map[string]float64)
	for stampId := range stampIdsMap {
		relations, err := e.repo.GetStampRelations(stampId)
		if err != nil {
			return 0, err
		}
		stampRelations[stampId] = make(map[string]float64)
		for _, r := range relations {
			if r.IDFrom != stampId {
				stampRelations[stampId][r.IDFrom] = r.Point
			} else {
				stampRelations[stampId][r.IDTo] = r.Point
			}
		}
	}

	var ret float64
	count := 0
	for i, s1 := range e.stamps {
		for j := i + 1; j < len(e.stamps); j++ {
			s2 := e.stamps[j]
			if s1.name == s2.name {
				continue
			}

			count++
			relations, ok := stampRelations[s1.id]
			if !ok {
				continue
			}
			pt, ok := relations[s2.id]
			if !ok {
				continue
			}
			ret += pt
		}
	}

	if count == 0 {
		return 0, nil
	}
	return math.Min(20, math.Max(0, ret/float64(count))), nil
}
