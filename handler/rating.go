package handler

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/traPtitech/traq-ws-bot/payload"

	"github.com/motoki317/bot-extreme/api"
	"github.com/motoki317/bot-extreme/repository"
)

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func handleShowRating(repo repository.Repository, p *payload.MessageCreated) {
	ratings, err := repo.GetAllRatings()
	if err != nil {
		err = respond(p, "レーティング一覧を取得中にエラーが発生しました...")
		if err != nil {
			log.Println(err)
		}
		return
	}

	if len(ratings) == 0 {
		err = respond(p, "一つもレーティングが存在しないようです。")
		if err != nil {
			log.Println(err)
		}
		return
	}

	users, err := api.GetUsers(true)
	if err != nil {
		err = respond(p, "API通信中にエラーが発生しました...")
		return
	}

	// descending sort
	sort.Slice(ratings, func(i, j int) bool {
		return ratings[i].Rating > ratings[j].Rating
	})

	// trim
	ratings = ratings[:min(100, len(ratings))]

	message := make([]string, 0, len(ratings)+3)
	message = append(message, "レーティング一覧")
	message = append(message, "| | ユーザー | レーティング |")
	message = append(message, "| - | - | - |")

	for i, r := range ratings {
		var name string
		if user, ok := users[r.ID]; ok {
			name = user.Name
		} else {
			name = r.ID
		}
		message = append(message, fmt.Sprintf("| %v. | :@%s: | %.2f |", i+1, name, r.Rating))
	}

	err = respond(p, strings.Join(message, "\n"))
	if err != nil {
		log.Println(err)
	}
}
