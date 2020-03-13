package evaluate

import (
	"math/rand"
	"time"
)

const (
	maxStamps = 5
)

// ランダムなスタンプとスタンプエフェクトを返します
// :ultrafastparrot.ex-large.rotate.parrot:
func GetRandomStampResponse() (string, error) {
	if len(stampsMap) == 0 {
		err := reCacheStamps()
		if err != nil {
			return "", err
		}
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	numStamps := r.Intn(maxStamps) + 1
	stampsList := make([]string, 0, len(stampsMap))
	stampsMapLock.RLock()
	defer stampsMapLock.RUnlock()
	for k := range stampsMap {
		stampsList = append(stampsList, k)
	}

	ret := ""
	for i := 0; i < numStamps; i++ {
		stamp := stampsList[r.Intn(len(stampsList))]

		sizeIndex := r.Intn(len(sizes) + 1)
		if sizeIndex < len(sizes) {
			stamp += "." + sizes[sizeIndex]
		}

		numEffects := r.Intn(6)
		for j := 0; j < numEffects; j++ {
			stamp += "." + moves[r.Intn(len(moves))]
		}

		ret += ":" + stamp + ":"
	}

	return ret, nil
}
