package evaluate

import (
	"math"
	"math/rand"
	"time"
)

type Winner int

const (
	FirstWins Winner = iota
	Even
	SecondWins
)

const (
	DefaultRating = 1500
	ratingTier    = 400
	ratingSpeed   = 32
)

// 評価された点数に応じ、勝敗を決定します。
func PickWinner(first, second float64) Winner {
	numFirst := math.Pow(math.E, first/10)
	numEven := math.Pow(math.E, (first+second)/30)
	numSecond := math.Pow(math.E, second/10)

	probFirst := numFirst / (numFirst + numEven + numSecond)
	probEven := numEven / (numFirst + numEven + numSecond)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if picked := r.Float64(); picked < probFirst {
		return FirstWins
	} else if picked < (probFirst + probEven) {
		return Even
	} else {
		return SecondWins
	}
}

func ChangeRating(winnerRating, loserRating float64) (newWinnerRating float64, newLoserRating float64) {
	probOfLoserWinning := 1 / (math.Pow(10, (winnerRating-loserRating)/ratingTier) + 1)
	newWinnerRating = winnerRating + probOfLoserWinning*ratingSpeed
	newLoserRating = loserRating - probOfLoserWinning*ratingSpeed
	return
}
