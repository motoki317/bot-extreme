package evaluate

import (
	"math"
	"math/rand"
	"time"
)

const (
	maxStamps = 3
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

	// generates 1 ~ 3 stamps exponentially
	numStamps := int(math.Min(maxStamps, math.Floor(r.ExpFloat64()*1.5+1)))

	stampsList := make([]string, 0, len(stampsMap))
	stampsMapLock.RLock()
	defer stampsMapLock.RUnlock()
	for k := range stampsMap {
		stampsList = append(stampsList, k)
	}

	ret := ""
	for i := 0; i < numStamps; i++ {
		// ランダムにスタンプをpick
		stamp := stampsList[r.Intn(len(stampsList))]

		// 30%の確率でサイズを付加
		sizeIndex := r.Intn(len(sizes)*3 + 1)
		if sizeIndex < len(sizes) {
			stamp += "." + sizes[sizeIndex]
		}

		// generates 0 ~ 5 exponentially
		numEffects := int(math.Min(5, math.Floor(r.ExpFloat64())))
		for j := 0; j < numEffects; j++ {
			stamp += "." + moves[r.Intn(len(moves))]
		}

		ret += ":" + stamp + ":"
	}

	return ret, nil
}
