package api

import (
	"github.com/antihax/optional"
	openapi "github.com/sapphi-red/go-traq"
)

// POST /channels/:channelId/messages メッセージをチャンネルに投稿、embedを自動変換（embed=1）
func PostMessage(channelId, content string) (*openapi.Message, error) {
	message, _, err := client.MessageApi.PostMessage(
		auth,
		channelId,
		openapi.SendMessage{Text: content},
		&openapi.PostMessageOpts{
			Embed: optional.NewInt32(1),
		},
	)
	return &message, err
}
