package handler

import (
	"fmt"
	"github.com/motoki317/bot-extreme/api"
	"github.com/motoki317/bot-extreme/janken"
	bot "github.com/traPtitech/traq-bot"
	"log"
	"os"
)

var (
	botUuid = os.Getenv("BOT_UUID")
)

func MessageReceived(processor *janken.Processor) func(payload *bot.MessageCreatedPayload) {
	return func(payload *bot.MessageCreatedPayload) {
		err := messageReceived(processor, payload)
		if err != nil {
			log.Println(err)
		}
	}
}

func messageReceived(processor *janken.Processor, payload *bot.MessageCreatedPayload) error {
	log.Println(fmt.Sprintf("[%s]: %s", payload.Message.User.DisplayName, payload.Message.PlainText))

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

	// TODO: handle message contents and change stamp ratings
	return nil
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
