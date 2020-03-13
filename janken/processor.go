package janken

import (
	"database/sql"
	"fmt"
	"github.com/motoki317/bot-extreme/evaluate"
	"github.com/motoki317/bot-extreme/repository"
	"log"
	"regexp"
	"strings"
	"time"
)

type Processor struct {
	repo  repository.Repository
	games map[string]*Game
}

func NewProcessor(repo repository.Repository) *Processor {
	return &Processor{
		repo:  repo,
		games: make(map[string]*Game),
	}
}

// ユーザーが参加しているゲームを取得
func (p *Processor) getCorrespondingGame(uuid string) *Game {
	for starter, game := range p.games {
		if starter == uuid {
			return game
		}
		if game.opponent != nil && game.opponent.ID == uuid {
			return game
		}
	}
	return nil
}

// メンションされたメッセージを処理します。何らかのアクションがあった場合、trueを返します。
func (p *Processor) Handle(sender *User, plainText string, mentioned []*User, respond func(string)) {
	err := p.handle(sender, plainText, mentioned, respond)
	if err != nil {
		log.Println(err)
		respond(fmt.Sprintf("エラーが発生しました: %s", err))
	}
}

func (p *Processor) handle(sender *User, plainText string, mentioned []*User, respond func(string)) (err error) {
	game := p.getCorrespondingGame(sender.ID)

	// じゃんけんしよう -> 新しくゲームを始める
	if regexp.MustCompile("\\s*じゃんけんしよう\\s*").MatchString(plainText) {
		if game != nil {
			// 既にユーザーがゲームを初めていた/ゲームに参加していた場合の処理
			respond("既にゲームを初めていませんか...？")
			return
		}

		p.games[sender.ID] = newGame(sender)
		content := strings.Join([]string{
			"いいですよ！",
			"今回対戦する相手を私にメンションで飛ばしてください！「`@" + botName + " @kashiwade`」",
			"もしくはBotと対戦するなら「`@" + botName + " ひとりで`」と",
			"やっぱり対戦をやめるなら「`@" + botName + " やっぱりいい`」と返してください！",
		}, "\n")
		respond(content)
		return
	}

	if game == nil {
		return
	}

	switch game.State {
	case WaitingOpponent:
		p.handlePickOpponent(game, plainText, respond, sender.ID, mentioned)
		return
	case PvB:
		game.opponent = &User{
			DisplayName: botName,
			ID:          botUuid,
		}
		game.opponentResponse, err = evaluate.GetRandomStampResponse()
		if err != nil {
			return err
		}
		fallthrough
	case PvP:
		return p.handlePvP(game, sender, respond, plainText)
	}

	return
}

func (p *Processor) handlePickOpponent(game *Game, plainText string, respond func(string), senderUuid string, mentioned []*User) {
	if regexp.MustCompile("\\s*ひとりで\\s*").MatchString(plainText) {
		// Player versus Bot
		game.State = PvB

		content := strings.Join([]string{
			"分かりました！",
			"私が「じゃーんけーん」と言ったら、じゃんけんの手をリプライしてください！",
		}, "\n")
		respond(content)

		go func() {
			<-time.NewTimer(time.Second * 1).C
			respond("じゃーんけーん")
		}()

		return
	} else if regexp.MustCompile("\\s*やっぱりいい\\s*").MatchString(plainText) {
		// Cancel game
		delete(p.games, senderUuid)

		respond("分かりました。またじゃんけんしましょう！")
	} else {
		// Pick opponent
		if len(mentioned) != 1 {
			respond(strings.Join([]string{
				"じゃんけんしたい相手を必ず**一人**指名してください！",
				"例: `@" + botName + " @kashiwade`",
			}, "\n"))
			return
		}

		if mentioned[0].ID == senderUuid {
			respond(strings.Join([]string{
				"自分自身と対戦ですか？面白い冗談ですね。",
				"じゃんけんしたい相手を必ず**一人**指名してください！",
				"例: `@" + botName + " @kashiwade`",
			}, "\n"))
			return
		}

		// picked opponent
		game.opponent = mentioned[0]
		game.State = PvP
		respond(strings.Join([]string{
			"@" + game.opponent.DisplayName + " 分かりました！",
			"私が「じゃーんけーん」と言ったら、じゃんけんの手をリプライしてください！",
		}, "\n"))

		go func() {
			<-time.NewTimer(time.Second * 1).C
			respond("@" + game.opponent.DisplayName + " じゃーんけーん")
		}()

		return
	}
}

func (p *Processor) handlePvP(game *Game, sender *User, respond func(string), plainText string) (err error) {
	if sender.ID == game.opponent.ID {
		// 相手の手
		if game.opponentResponse != "" {
			respond("二度出しはできません！")
			return
		}
		game.opponentResponse = plainText
	} else {
		// 自分の手
		if game.selfResponse != "" {
			respond("二度出しはできません！")
			return
		}
		game.selfResponse = plainText
	}

	if game.opponentResponse == "" || game.selfResponse == "" {
		return
	}

	// 両方の手が集まったならば、評価する
	selfPoints, err := evaluate.Message(p.repo, game.selfResponse)
	if err != nil {
		return err
	}
	opponentPoints, err := evaluate.Message(p.repo, game.opponentResponse)
	if err != nil {
		return err
	}

	result := evaluate.PickWinner(selfPoints, opponentPoints)

	log.Printf("%s %v pts vs %s %v pts, result: %v\n",
		game.self.DisplayName, selfPoints,
		game.opponent.DisplayName, opponentPoints, result)

	// 引き分け
	if result == evaluate.Even {
		game.selfResponse = ""
		game.opponentResponse = ""
		respond("@" + game.opponent.DisplayName + " あーいこーで")
		return
	}

	// 勝敗が決定した、ゲームを終了する
	delete(p.games, game.self.ID)
	// レーティングを計算
	selfRating, err := p.getRatingOrDefault(game.self.ID)
	if err != nil {
		return err
	}
	opponentRating, err := p.getRatingOrDefault(game.opponent.ID)
	if err != nil {
		return err
	}

	var oldSelfRating, oldOpponentRating float64
	response := []string{"@" + game.opponent.DisplayName, "", ""}

	if result == evaluate.FirstWins {
		// 自分の勝ち
		response = append(response, ":"+game.self.DisplayName+": の勝ちです！")
		selfRating.Rating, opponentRating.Rating = evaluate.ChangeRating(selfRating.Rating, opponentRating.Rating)
	} else {
		// 相手の勝ち
		response = append(response, ":"+game.opponent.DisplayName+": の勝ちです！")
		opponentRating.Rating, selfRating.Rating = evaluate.ChangeRating(opponentRating.Rating, selfRating.Rating)
	}

	err = p.repo.UpdateRating(selfRating)
	if err != nil {
		return
	}
	err = p.repo.UpdateRating(opponentRating)
	if err != nil {
		return
	}

	response = append(response, "", "")
	response = append(response, "新しいレーティングは")
	response = append(response, fmt.Sprintf(":%s: : %v (%+v)", game.self.DisplayName, int(selfRating.Rating), int(selfRating.Rating-oldSelfRating)))
	response = append(response, fmt.Sprintf(":%s: : %v (%+v)", game.opponent.DisplayName, int(opponentRating.Rating), int(opponentRating.Rating-oldOpponentRating)))
	response = append(response, "です！")

	respond(strings.Join(response, "\n"))
	return
}

// ユーザーIDのレーティングを返します。存在しない場合、デフォルトを生成し返します。
func (p *Processor) getRatingOrDefault(ID string) (*repository.Rating, error) {
	if rating, err := p.repo.GetRating(ID); err == nil {
		return rating, nil
	} else if err == sql.ErrNoRows {
		return &repository.Rating{
			ID:     ID,
			Rating: evaluate.DefaultRating,
		}, nil
	} else {
		return nil, err
	}
}
