package telegram

import (
	"expenses-app/pkg/app/querying"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func listCategories(msg *tgbotapi.MessageConfig, q *querying.Service) {
	cg := q.CategoryQuerier
	resp, err := cg.Query()
	if err != nil {
		msg.Text = err.Error()
	}
	categories := []string{}
	for _, name := range resp.Categories {
		categories = append(categories, name)
	}
	msg.Text = "What category ?"
	msg.ReplyMarkup = GetKeyBoardMap(categories, 4)
}
