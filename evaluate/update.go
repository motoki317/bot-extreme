package evaluate

import (
	"math"
	"sort"

	openapi "github.com/sapphi-red/go-traq"

	"github.com/motoki317/bot-extreme/repository"
)

// メッセージを取得し、スタンプの評価などを変更します

const (
	defaultEffectPoint = 1
)

type reaction struct {
	// stamp ID
	id string
}

type Message struct {
	MessageStamps []*stamp
	// users -> list of stamps
	UserReactions [][]*reaction
}

func ParseMessage(message *openapi.Message) *Message {
	userReactions := make(map[string][]*reaction)
	for _, s := range message.Stamps {
		userReactions[s.UserId] = append(userReactions[s.UserId], &reaction{
			id: s.StampId,
		})
	}

	users := make([][]*reaction, 0, len(userReactions))
	for _, reactions := range userReactions {
		users = append(users, reactions)
	}

	return &Message{
		// 本当は存在するスタンプを再fetchした方が良いがあえてしない
		MessageStamps: filterUnknownStamp(getMessageStamps(message.Content)),
		UserReactions: users,
	}
}

// 更新用
type updater struct {
	repo    repository.Repository
	message *Message
}

// 1つのメッセージについて処理し、スタンプの評価などを変更します
func ProcessMessage(repo repository.Repository, message *Message) error {
	if len(message.UserReactions) == 0 && len(message.MessageStamps) == 0 {
		return nil
	}

	updater := &updater{
		repo:    repo,
		message: message,
	}

	// スタンプ使用数を更新
	err := updater.updateStampsUsed()
	if err != nil {
		return err
	}

	// スタンプエフェクトを更新
	err = updater.updateStampEffects()
	if err != nil {
		return err
	}

	// スタンプの関係を更新
	err = updater.updateStampRelations()
	if err != nil {
		return err
	}

	return nil
}

func (u *updater) updateStampsUsed() error {
	u.repo.Lock()
	defer u.repo.Unlock()

	// このメッセージの中でそれぞれのスタンプIDが何回使われたか
	counts := make(map[string]int)
	for _, s := range u.message.MessageStamps {
		count := 0
		if _, ok := counts[s.id]; ok {
			count = counts[s.id]
		}
		count++
		counts[s.id] = count
	}
	for _, userReaction := range u.message.UserReactions {
		for _, r := range userReaction {
			count := 0
			if _, ok := counts[r.id]; ok {
				count = counts[r.id]
			}
			count++
			counts[r.id] = count
		}
	}

	// 更新
	for stampID, count := range counts {
		stampInfo, err := u.repo.GetStamp(stampID)
		if err != nil {
			return err
		}
		if stampInfo == nil {
			stampInfo = &repository.Stamp{
				ID:   stampID,
				Used: 0,
			}
		}

		stampInfo.Used += count
		err = u.repo.UpdateStamp(stampInfo)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *updater) updateStampEffects() error {
	u.repo.Lock()
	defer u.repo.Unlock()

	// このメッセージ内で使われた大きさと動きに関するエフェクトの数
	sizeCounts := make(map[string]int)
	moveCounts := make(map[string]int)
	for _, s := range u.message.MessageStamps {
		if s.sizeEffect != "" {
			count := 0
			if _, ok := sizeCounts[s.sizeEffect]; ok {
				count = sizeCounts[s.sizeEffect]
			}
			count++
			sizeCounts[s.sizeEffect] = count
		}

		for _, eff := range s.moveEffects {
			count := 0
			if _, ok := moveCounts[eff]; ok {
				count = moveCounts[eff]
			}
			count++
			moveCounts[eff] = count
		}
	}

	// 現在のポイントを取得
	allEffects, err := u.repo.GetAllEffectPoints()
	if err != nil {
		return err
	}
	sizeEffectPoints := make(map[string]float64)
	moveEffectPoints := make(map[string]float64)

	for _, e := range allEffects {
		if _, ok := sizeEffects[e.Name]; ok {
			sizeEffectPoints[e.Name] = e.Point
		} else {
			moveEffectPoints[e.Name] = e.Point
		}
	}
	// (存在しなかったら)デフォルトのポイントを設定
	for _, size := range sizes {
		if _, ok := sizeEffectPoints[size]; !ok {
			sizeEffectPoints[size] = defaultEffectPoint
		}
	}
	for _, move := range moves {
		if _, ok := moveEffectPoints[move]; !ok {
			moveEffectPoints[move] = defaultEffectPoint
		}
	}

	// 更新
	// 暫定的にkeyをソートしてその順で処理する
	sizeCountKeys := make([]string, 0, len(sizeCounts))
	for k := range sizeCounts {
		sizeCountKeys = append(sizeCountKeys, k)
	}
	sort.Strings(sizeCountKeys)
	for _, eff := range sizeCountKeys {
		for i := 0; i < sizeCounts[eff]; i++ {
			for otherEff, point := range sizeEffectPoints {
				if otherEff == eff {
					continue
				}
				diff := point * 0.01
				sizeEffectPoints[otherEff] -= diff
				sizeEffectPoints[eff] += diff
			}
		}
	}
	moveCountKeys := make([]string, 0, len(moveCounts))
	for k := range moveCounts {
		moveCountKeys = append(moveCountKeys, k)
	}
	sort.Strings(moveCountKeys)
	for _, eff := range moveCountKeys {
		for i := 0; i < moveCounts[eff]; i++ {
			for otherEff, point := range moveEffectPoints {
				if otherEff == eff {
					continue
				}
				diff := point * 0.01
				moveEffectPoints[otherEff] -= diff
				moveEffectPoints[eff] += diff
			}
		}
	}

	// 更新を反映
	allEffects = make([]*repository.EffectPoint, 0, len(allEffects))
	for eff, point := range sizeEffectPoints {
		allEffects = append(allEffects, &repository.EffectPoint{
			Name:  eff,
			Point: point,
		})
	}
	for eff, point := range moveEffectPoints {
		allEffects = append(allEffects, &repository.EffectPoint{
			Name:  eff,
			Point: point,
		})
	}
	return u.repo.UpdateAllEffectPoints(allEffects)
}

func (u *updater) updateStampRelations() error {
	u.repo.Lock()
	defer u.repo.Unlock()

	// 使われたスタンプ
	stampsUsed := make(map[string]bool)
	for _, s := range u.message.MessageStamps {
		stampsUsed[s.id] = true
	}
	for _, userReaction := range u.message.UserReactions {
		for _, r := range userReaction {
			stampsUsed[r.id] = true
		}
	}

	// 現在の関係を取得
	relations := make(map[string]map[string]*repository.StampRelation)
	for stampID := range stampsUsed {
		stampRelations, err := u.repo.GetStampRelations(stampID)
		if err != nil {
			return err
		}

		for _, toAdd := range stampRelations {
			if _, ok := relations[toAdd.IDFrom]; !ok {
				relations[toAdd.IDFrom] = make(map[string]*repository.StampRelation)
			}
			if _, ok := relations[toAdd.IDTo]; !ok {
				relations[toAdd.IDTo] = make(map[string]*repository.StampRelation)
			}
			// 逆方向にも同じポインタを張る
			relations[toAdd.IDFrom][toAdd.IDTo] = toAdd
			relations[toAdd.IDTo][toAdd.IDFrom] = toAdd
		}
	}

	// 更新
	messageStampSet := make(map[string]bool)
	for _, s := range u.message.MessageStamps {
		messageStampSet[s.id] = true
	}
	u.updateRelationForStampSet(messageStampSet, relations)

	for _, userReaction := range u.message.UserReactions {
		stampSet := make(map[string]bool)
		for _, r := range userReaction {
			stampSet[r.id] = true
		}
		u.updateRelationForStampSet(stampSet, relations)
	}

	// 更新を保存
	relationSet := make(map[*repository.StampRelation]bool)
	for _, relation := range relations {
		for _, r := range relation {
			relationSet[r] = true
		}
	}

	toSave := make([]*repository.StampRelation, 0, len(relationSet))
	for r := range relationSet {
		toSave = append(toSave, r)
	}

	err := u.repo.UpdateStampRelations(toSave)
	if err != nil {
		return err
	}
	return u.repo.DeleteStampRelations(0)
}

func (u *updater) updateRelationForStampSet(stampSet map[string]bool, relations map[string]map[string]*repository.StampRelation) {
	stamps := make([]string, 0, len(stampSet))
	for stampID := range stampSet {
		stamps = append(stamps, stampID)
	}

	getStampRelation := func(from, to string) *repository.StampRelation {
		if _, ok := relations[from]; ok {
			if r, ok := relations[from][to]; ok {
				return r
			}
		}
		return nil
	}

	for i, s1 := range stamps {
		if _, ok := relations[s1]; !ok {
			relations[s1] = make(map[string]*repository.StampRelation)
		}
		stampRelation := relations[s1]

		// stamp setに含まれるものを更新
		for j := i + 1; j < len(stamps); j++ {
			s2 := stamps[j]
			relation := getStampRelation(s1, s2)
			if relation == nil {
				relation = &repository.StampRelation{
					IDFrom: s1,
					IDTo:   s2,
					Point:  0,
				}
				if _, ok := relations[s2]; !ok {
					relations[s2] = make(map[string]*repository.StampRelation)
				}
				relations[s1][s2] = relation
				relations[s2][s1] = relation
			}

			relations[s1][s2].Point += 1
		}

		// stamp setに含まれないrelationについて更新
		for s2, r := range stampRelation {
			if _, ok := stampSet[s2]; ok {
				continue
			}
			r.Point = math.Max(0, r.Point-0.1)
		}
	}
}
