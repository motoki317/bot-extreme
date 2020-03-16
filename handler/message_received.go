package handler

import (
	"fmt"
	"github.com/motoki317/bot-extreme/api"
	"github.com/motoki317/bot-extreme/janken"
	"github.com/motoki317/bot-extreme/repository"
	bot "github.com/motoki317/traq-bot"
	"log"
	"os"
)

var (
	botUuid = os.Getenv("BOT_UUID")
)

func MessageReceived(repo repository.Repository) func(payload *bot.MessageCreatedPayload) {
	processor := janken.NewProcessor(repo)
	updater := &updater{
		repo: repo,
	}

	return func(payload *bot.MessageCreatedPayload) {
		log.Println(fmt.Sprintf("[%s]: %s", payload.Message.User.DisplayName, payload.Message.PlainText))

		handleJanken(processor, payload)

		// より古いメッセージを処理しスタンプのレーティングを更新する
		go func() {
			err := updater.updateRatings(payload.Message.ChannelID)
			if err != nil {
				log.Println(err)
			}
		}()
	}
}

func handleJanken(processor *janken.Processor, payload *bot.MessageCreatedPayload) {
	if isMentioned(payload) {
		sender := &janken.User{
			DisplayName: payload.Message.User.DisplayName,
			ID:          payload.Message.User.ID,
		}
		plainText := payload.Message.PlainText

		processor.Handle(sender, plainText, getMentionedUsers(payload), func(content string) {
			err := respond(payload, content)
			if err != nil {
				log.Println(err)
			}
		})
	}
}

// メッセージ内でBotがメンションされたかを判定
func isMentioned(payload *bot.MessageCreatedPayload) bool {
	for _, e := range payload.Message.Embedded {
		if e.Type == "user" && e.ID == botUuid {
			return true
		}
	}
	return false
}

// Bot自身を除くメンションされたユーザー一覧を取得
func getMentionedUsers(payload *bot.MessageCreatedPayload) (users []*janken.User) {
	for _, e := range payload.Message.Embedded {
		if e.Type == "user" && e.ID != botUuid {
			users = append(users, &janken.User{
				// e.Raw example: "@takashi_trap"
				DisplayName: e.Raw[1:],
				ID:          e.ID,
			})
		}
	}
	return
}

// @sender {content}のように送られたチャンネルへ返信する
func respond(payload *bot.MessageCreatedPayload, content string) (err error) {
	_, err = api.PostMessage(payload.Message.ChannelID, fmt.Sprintf("@%s %s", payload.Message.User.Name, content))
	return
}
