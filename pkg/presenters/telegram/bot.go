package telegram

import (
	"expenses-app/pkg/app/authenticating"
	"expenses-app/pkg/app/health"
	"expenses-app/pkg/app/managing"
	"expenses-app/pkg/app/querying"
	"expenses-app/pkg/app/tracking"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const HELP_MSG = "I understand the following commands \n /ListCategories \n /Status \n /N26MonthTracking"

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

func Run(tbot *tgbotapi.BotAPI, h *health.Service, m *managing.Service, t *tracking.Service, a *authenticating.Service, q *querying.Service) {
	tbot.Debug = true
	log.Printf("Authorized on account %s", tbot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := tbot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message.IsCommand() {
			tbot.Send(handleCommands(update, &updates, tbot, q, t, h))
			// } else if update.Message.Document != nil {
		} else {
			tbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, HELP_MSG))
		}
	}
}

func handleCommands(update tgbotapi.Update, updates *tgbotapi.UpdatesChannel, tbot *tgbotapi.BotAPI, q *querying.Service, t *tracking.Service, h *health.Service) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	switch update.Message.Command() {
	case "help":
		msg.Text = HELP_MSG
	case "ListCategories":
		listCategories(&msg, q)
	case "N26MonthTracking":
		n26MonthTracking(&msg, tbot, &update, updates, t)
	case "Status":
		ping(&msg, h)
	default:
		msg.Text = HELP_MSG
	}
	return msg
}

func ping(msg *tgbotapi.MessageConfig, h *health.Service) {
	msg.Text = h.Ping()
}
