package handler

import (
	"github.com/motoki317/bot-extreme/api"
	"github.com/motoki317/bot-extreme/evaluate"
	"github.com/motoki317/bot-extreme/repository"
	"log"
	"time"
)

const (
	processOlderThan = time.Hour * 3
)

type updater struct {
	repo repository.Repository
}

// channelIDのメッセージ全てを更新する
func (u *updater) updateRatings(channelID string) error {
	u.repo.ChannelLock()
	defer u.repo.ChannelUnlock()

	seen, err := u.repo.GetSeenChannel(channelID)
	if err != nil {
		return err
	}
	if seen == nil {
		seen = &repository.SeenChannel{
			ID:                   channelID,
			LastProcessedMessage: time.Unix(0, 0),
		}
	}

	// LastProcessedMessageかそれより前に送られたメッセージににたどり着くか、
	// チャンネルの最後まで読み込む
	processUntil := time.Now().Add(-processOlderThan)
	limit := api.DefaultMessageLimit
	offset := 0
	for {
		log.Printf("Processing for channel %s, limit %v, offset %v...\n", channelID, limit, offset)
		messages, hasMore, err := api.GetChannelMessages(channelID, limit, offset)
		if err != nil {
			return err
		}

		processedAll := false
		for _, m := range messages {
			if m.CreatedAt.After(processUntil) {
				continue
			}
			if m.CreatedAt.Equal(seen.LastProcessedMessage) || m.CreatedAt.Before(seen.LastProcessedMessage) {
				processedAll = true
				break
			}

			message := m
			err := evaluate.ProcessMessage(u.repo, evaluate.ParseMessage(&message))
			if err != nil {
				return err
			}
		}

		if processedAll || !hasMore {
			break
		}

		offset += limit

		// 5秒間隔でメッセージを取りに行く
		<-time.NewTimer(time.Second * 1).C
	}

	log.Printf("Processed for channel %s, from %v to %v", channelID, seen.LastProcessedMessage, processUntil)

	// どこまで見たかを保存する
	seen.LastProcessedMessage = processUntil
	return u.repo.UpdateSeenChannel(seen)
}
