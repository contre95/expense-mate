package telegram

import (
	"expenses-app/pkg/app/querying"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func listCategories(tbot *tgbotapi.BotAPI, update *tgbotapi.Update, q *querying.Service) {
	cg := q.CategoryQuerier
	resp, err := cg.Query()
	if err != nil {
		tbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
	}
	for _, name := range resp.Categories {
		tbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, name))
	}
}
