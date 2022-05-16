package evaluate

import (
	"log"
	"regexp"
	"strings"

	"github.com/motoki317/bot-extreme/api"
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

func filterUnknownStamp(stamps []*stamp) []*stamp {
	stampsMapLock.RLock()
	defer stampsMapLock.RUnlock()
	ret := make([]*stamp, 0, len(stamps))
	for _, s := range stamps {
		if stampID, ok := stampsMap[s.name]; ok {
			s.id = stampID
			ret = append(ret, s)
		}
	}
	return ret
}

func reCacheStamps() error {
	log.Println("Requesting all stamps and users")

	stamps, err := api.GetStamps()
	if err != nil {
		return err
	}

	users, err := api.GetUsers(false)
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
		stampsMap["@"+u.Name] = u.Id
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

// スタンプをパースします。エフェクトの有効性についてvalidateします。
// e.g. "thonk_spin.ex-large.rotate.parrot" -> "thonk_spin", []string{"ex-large", "rotate", "parrot"}
func parseStamp(stampMessage string) *stamp {
	matches := strings.Split(stampMessage, ".")

	ret := &stamp{
		name: matches[0],
	}

	// NOTE: ここではまだスタンプの存在はチェックしない

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
