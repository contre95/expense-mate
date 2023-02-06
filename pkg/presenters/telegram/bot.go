package telegram

import (
	"expenses-app/pkg/app/authenticating"
	"expenses-app/pkg/app/health"
	"expenses-app/pkg/app/importing"
	"expenses-app/pkg/app/managing"
	"expenses-app/pkg/app/querying"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetKeyBoardMap(items []string, rowsCant int) tgbotapi.ReplyKeyboardMarkup {
	matrix := [][]tgbotapi.KeyboardButton{}
	row := []tgbotapi.KeyboardButton{}
	for _, category := range items {
		newButton := tgbotapi.NewKeyboardButton(category)
		row = append(row, newButton)
		if len(row) == rowsCant || len(items)-len(matrix)*rowsCant == len(row) {
			matrix = append(matrix, row)
			row = []tgbotapi.KeyboardButton{}
		}
	}
	return tgbotapi.NewReplyKeyboard(matrix...)
}

func Run(tbot *tgbotapi.BotAPI, h *health.Service, m *managing.Service, i *importing.Service, a *authenticating.Service, q *querying.Service) {
	tbot.Debug = true
	log.Printf("Authorized on account %s", tbot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := tbot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}
		if update.Message.IsCommand() { // ignore any non-command Messages
			tbot.Send(handleCommands(update, q))
		}
	}
}

func handleCommands(update tgbotapi.Update, q *querying.Service) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	switch update.Message.Command() {
	case "help":
		msg.Text = "I understand /categories and /status."
	case "categories":
		categories, err := getCategories(q.CategoryQuerier)
		if err != nil {
			msg.Text = err.Error()
		}
		msg.Text = "What category ?"
		msg.ReplyMarkup = GetKeyBoardMap(categories, 4)
	case "status":
		msg.Text = "I'm ok."
	default:
		msg.Text = "I don't know that command"

	}
	return msg
}
