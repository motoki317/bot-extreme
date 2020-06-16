package janken

import (
	"fmt"
	"github.com/motoki317/bot-extreme/evaluate"
	"github.com/motoki317/bot-extreme/repository"
	"log"
	"math"
	"regexp"
	"strings"
	"time"
)

type Processor struct {
	repo repository.Repository
	// ゲームを始めた人のuuid -> game
	games map[string]*Game
}

func NewProcessor(repo repository.Repository) *Processor {
	return &Processor{
		repo:  repo,
		games: make(map[string]*Game),
	}
}

// ユーザーが参加しているゲームを取得 自分がスタートしたゲームを優先的に取得します
func (p *Processor) getCorrespondingGame(uuid string) *Game {
	if game, ok := p.games[uuid]; ok {
		return game
	}
	for _, game := range p.games {
		if game.opponent != nil && game.opponent.ID == uuid {
			return game
		}
	}
	return nil
}

// メンションされたメッセージを処理します。
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
	case OpponentPick:
		p.handlePickOpponent(game, plainText, respond, sender.ID, mentioned)
		return
	case WaitingOpponent:
		p.handleOpponentResponse(game, plainText, respond, sender.ID)
		return
	case PvB:
		game.opponent = &User{
			Name: botName,
			ID:   botUuid,
		}
		game.opponentResponse, err = evaluate.GetRandomStampResponse()
		respond(game.opponentResponse)
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

		respond(strings.Join([]string{
			"分かりました！",
			"私が「じゃーんけーん」と言ったら、選んだじゃんけんの手を私にリプライしてください！",
			"`@BOT_extreme :ultrafastparrot:`",
		}, "\n"))

		go func() {
			<-time.NewTimer(time.Second * 3).C
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

		// Check validity
		if mentioned[0].ID == senderUuid {
			respond(strings.Join([]string{
				"自分自身と対戦ですか？面白い冗談ですね。",
				"じゃんけんしたい相手を必ず**一人**指名してください！",
				"例: `@" + botName + " @kashiwade`",
			}, "\n"))
			return
		}

		// Check bot
		if strings.HasPrefix(mentioned[0].Name, "BOT_") {
			respond("Botと対戦ですか？考えましたね、でも残念、それはできません。")
			return
		}

		// picked opponent
		game.opponent = mentioned[0]
		game.State = WaitingOpponent
		// wait for opponent response
		respond(strings.Join([]string{
			"分かりました！",
			"",
			"",
			"@" + game.opponent.Name + " さん、準備ができたら私にリプライを飛ばしてください！",
			"`@" + botName + "`",
			"",
			"",
			"@" + game.self.Name + " さん、準備ができないようだったら私に「やっぱりいい」とリプライしてください。",
			"`@" + botName + " やっぱりいい`",
		}, "\n"))

		return
	}
}

func (p *Processor) handleOpponentResponse(game *Game, plainText string, respond func(string), senderUuid string) {
	if senderUuid == game.opponent.ID {
		// opponent responded
		game.State = PvP
		respond(strings.Join([]string{
			"@" + game.self.Name + " 分かりました！",
			"私が「じゃーんけーん」と言ったら、二人ともじゃんけんの手を私にリプライしてください！",
			"`@BOT_extreme :ultrafastparrot:`",
		}, "\n"))

		go func() {
			<-time.NewTimer(time.Second * 3).C
			respond("@" + game.self.Name + " じゃーんけーん")
		}()
	} else if senderUuid == game.self.ID {
		if regexp.MustCompile("\\s*やっぱりいい\\s*").MatchString(plainText) {
			// Cancel game
			delete(p.games, senderUuid)

			respond("分かりました。またじゃんけんしましょう！")
		}
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
	selfPoints, err := evaluate.MessagePoint(p.repo, game.selfResponse)
	if err != nil {
		return err
	}
	opponentPoints, err := evaluate.MessagePoint(p.repo, game.opponentResponse)
	if err != nil {
		return err
	}

	result := evaluate.PickWinner(selfPoints, opponentPoints)

	log.Printf("%s %v pts vs %s %v pts, result: %v\n",
		game.self.Name, selfPoints,
		game.opponent.Name, opponentPoints, result)

	// 引き分け
	if result == evaluate.Even {
		game.selfResponse = ""
		game.opponentResponse = ""
		respond("@" + game.opponent.Name + " あーいこーで")
		return
	}

	// 勝敗が決定した、ゲームを終了する
	delete(p.games, game.self.ID)

	// respond関数はsenderにタグするので、両方にタグするため
	var otherName string
	if sender.ID == game.opponent.ID {
		otherName = game.self.Name
	} else {
		otherName = game.opponent.Name
	}

	response := []string{
		"@" + otherName,
		"",
		"",
		fmt.Sprintf(":@%s: %.2f pts - :@%s: %.2f pts で...", game.self.Name, selfPoints, game.opponent.Name, opponentPoints),
		"",
		"",
	}
	if result == evaluate.FirstWins {
		// 自分の勝ち
		response = append(response, ":@"+game.self.Name+": の勝ちです！")
	} else {
		// 相手の勝ち
		response = append(response, ":@"+game.opponent.Name+": の勝ちです！")
	}

	if game.State == PvP {
		// PvPならレーティングを計算
		selfRating, err := p.getRatingOrDefault(game.self.ID)
		if err != nil {
			return err
		}
		opponentRating, err := p.getRatingOrDefault(game.opponent.ID)
		if err != nil {
			return err
		}

		oldSelfRating, oldOpponentRating := selfRating.Rating, opponentRating.Rating

		if result == evaluate.FirstWins {
			// 自分の勝ち
			selfRating.Rating, opponentRating.Rating = evaluate.ChangeRating(selfRating.Rating, opponentRating.Rating)
		} else {
			// 相手の勝ち
			opponentRating.Rating, selfRating.Rating = evaluate.ChangeRating(opponentRating.Rating, selfRating.Rating)
		}

		err = p.repo.UpdateRating(selfRating)
		if err != nil {
			return err
		}
		err = p.repo.UpdateRating(opponentRating)
		if err != nil {
			return err
		}

		response = append(response, "")
		response = append(response, "")
		response = append(response, "新しいレーティングは")
		response = append(response, fmt.Sprintf(":@%s: %d (%+d)",
			game.self.Name, int(math.Round(selfRating.Rating)), int(math.Round(selfRating.Rating-oldSelfRating))))
		response = append(response, fmt.Sprintf(":@%s: %d (%+d)",
			game.opponent.Name, int(math.Round(opponentRating.Rating)), int(math.Round(opponentRating.Rating-oldOpponentRating))))
		response = append(response, "です！")
	} else {
		response = append(response, "私との対戦なのでレーティング変動はありません。")
		response = append(response, "ちなみにいまのレーティングは")
		selfRating, err := p.getRatingOrDefault(game.self.ID)
		if err != nil {
			return err
		}
		response = append(response, fmt.Sprintf(":@%s: %v", game.self.Name, int(selfRating.Rating)))
		response = append(response, "です！")
	}

	go func() {
		<-time.NewTimer(time.Second * 3).C
		respond(strings.Join(response, "\n"))
	}()
	return
}

// ユーザーIDのレーティングを返します。存在しない場合、デフォルトを生成し返します。
func (p *Processor) getRatingOrDefault(ID string) (*repository.Rating, error) {
	if rating, err := p.repo.GetRating(ID); err == nil {
		if rating == nil {
			return &repository.Rating{
				ID:     ID,
				Rating: evaluate.DefaultRating,
			}, nil
		}
		return rating, nil
	} else {
		return nil, err
	}
}
