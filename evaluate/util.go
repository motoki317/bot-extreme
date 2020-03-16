package evaluate

import (
	"github.com/motoki317/bot-extreme/api"
	"log"
	"regexp"
	"strings"
)

func containsUnknownStamp(stamps []*stamp) bool {
	stampsMapLock.RLock()
	defer stampsMapLock.RUnlock()
	for _, s := range stamps {
		if _, ok := stampsMap[s.name]; !ok {
			log.Println("Found unknown stamp " + s.name)
			return true
		}
	}
	return false
}

func reCacheStamps() error {
	log.Println("Requesting all stamps and users")

	stamps, err := api.GetStamps()
	if err != nil {
		return err
	}

	users, err := api.GetUsers()
	if err != nil {
		return err
	}

	stampsMapLock.Lock()
	defer stampsMapLock.Unlock()

	stampsMap = make(map[string]string)

	for _, s := range stamps {
		stampsMap[s.Name] = s.Id
	}
	for _, u := range users {
		stampsMap[u.Name] = u.UserId
	}

	return nil
}

// メッセージ中に含まれるスタンプをパースします
func getMessageStamps(content string) []*stamp {
	messageStamps := regexp.MustCompile(":(.+?):").FindAllStringSubmatch(content, -1)
	if messageStamps == nil {
		return nil
	}
	stamps := make([]*stamp, 0, len(messageStamps))
	for _, stampMessage := range messageStamps {
		parsedStamp := parseStamp(stampMessage[1])
		if parsedStamp != nil {
			stamps = append(stamps, parsedStamp)
		}
	}
	return stamps
}

// スタンプをパースします。存在するスタンプであるかどうか、またエフェクトの有効性についてもvalidateします
// e.g. "thonk_spin.ex-large.rotate.parrot" -> "thonk_spin", []string{"ex-large", "rotate", "parrot"}
func parseStamp(stampMessage string) *stamp {
	matches := strings.Split(stampMessage, ".")

	ret := &stamp{
		name: matches[0],
	}

	if len(stampsMap) == 0 {
		err := reCacheStamps()
		if err != nil {
			log.Println(err)
			return nil
		}
	}

	stampsMapLock.RLock()
	defer stampsMapLock.RUnlock()

	// 存在するスタンプかチェック
	if id, ok := stampsMap[ret.name]; ok {
		ret.id = id
	} else {
		return nil
	}

	for _, eff := range matches[1:] {
		if _, ok := sizeEffects[eff]; ok {
			ret.sizeEffect = eff
		} else if _, ok := moveEffects[eff]; ok {
			ret.moveEffects = append(ret.moveEffects, eff)
		} else {
			// invalid stamp
			return nil
		}
	}

	// invalid stamp move effects
	if len(ret.moveEffects) > 5 {
		return nil
	}
	return ret
}
