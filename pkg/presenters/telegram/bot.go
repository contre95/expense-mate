package telegram

import (
	"expenses-app/pkg/app/authenticating"
	"expenses-app/pkg/app/health"
	"expenses-app/pkg/app/managing"
	"expenses-app/pkg/app/querying"
	"expenses-app/pkg/app/tracking"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const HELP_MSG = "I understand the following commands \n /ListCategories \n /Status \n /N26MonthTracking"

// Register the following commands in the botfather

// categories - Get a list of categories
// n26importer - Import expenses from N26 csv file export
// ping - Check if the bot is working
func Run(tbot *tgbotapi.BotAPI, h *health.Service, m *managing.Service, t *tracking.Service, a *authenticating.Service, q *querying.Service) {
	tgbotapi.NewRemoveKeyboard(true)
	tbot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := tbot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message.IsCommand() {
			handleCommands(update, &updates, tbot, q, t, h)
		} else {
			tbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, HELP_MSG))
		}
	}
}

func handleCommands(update tgbotapi.Update, updates *tgbotapi.UpdatesChannel, tbot *tgbotapi.BotAPI, q *querying.Service, t *tracking.Service, h *health.Service) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	switch update.Message.Command() {
	case "categories":
		listCategories(tbot, &update, q)
	case "n26importer":
		n26MonthTracking(tbot, &update, updates, t, q)
	case "ping":
		ping(tbot, update, h)
	default:
		tbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, HELP_MSG))
	}
	return msg
}

func ping(tbot *tgbotapi.BotAPI, update tgbotapi.Update, h *health.Service) {
	tbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, h.Ping()))
}
