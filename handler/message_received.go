package handler

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/traPtitech/traq-ws-bot/payload"

	"github.com/motoki317/bot-extreme/api"
	"github.com/motoki317/bot-extreme/janken"
	"github.com/motoki317/bot-extreme/repository"
)

var (
	botUuid = os.Getenv("BOT_UUID")
)

func MessageReceived(repo repository.Repository) func(p *payload.MessageCreated) {
	processor := janken.NewProcessor(repo)
	updater := &updater{
		repo: repo,
	}

	return func(p *payload.MessageCreated) {
		log.Println(fmt.Sprintf("[%s]: %s", p.Message.User.DisplayName, p.Message.PlainText))

		// Botからのメッセージは処理しない
		if p.Message.User.Bot {
			return
		}

		// メンションされたときのみコマンドを処理する
		if !isMentioned(p) {
			return
		}

		// レーティング表示
		if regexp.MustCompile("\\s*ランキング\\s*").MatchString(p.Message.PlainText) {
			handleShowRating(repo, p)
			return
		}

		handleJanken(processor, p)

		// より古いメッセージを処理しスタンプのレーティングを更新する
		go func() {
			err := updater.updateRatings(p.Message.ChannelID)
			if err != nil {
				log.Println(err)
			}
		}()
	}
}

func handleJanken(processor *janken.Processor, p *payload.MessageCreated) {
	sender := &janken.User{
		Name: p.Message.User.Name,
		ID:   p.Message.User.ID,
	}
	plainText := p.Message.PlainText

	processor.Handle(sender, plainText, getMentionedUsers(p), func(content string) {
		err := respond(p, content)
		if err != nil {
			log.Println(err)
		}
	})
}

// メッセージ内でBotがメンションされたかを判定
func isMentioned(p *payload.MessageCreated) bool {
	for _, e := range p.Message.Embedded {
		if e.Type == "user" && e.ID == botUuid {
			return true
		}
	}
	return false
}

// Bot自身を除くメンションされたユーザー一覧を取得
func getMentionedUsers(p *payload.MessageCreated) (users []*janken.User) {
	for _, e := range p.Message.Embedded {
		if e.Type == "user" && e.ID != botUuid {
			users = append(users, &janken.User{
				// e.Raw example: "@takashi_trap"
				Name: e.Raw[1:],
				ID:   e.ID,
			})
		}
	}
	return
}

// @sender {content}のように送られたチャンネルへ返信する
func respond(p *payload.MessageCreated, content string) (err error) {
	_, err = api.PostMessage(p.Message.ChannelID, fmt.Sprintf("@%s %s", p.Message.User.Name, content))
	return
}
