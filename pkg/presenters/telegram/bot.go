package telegram

import (
	"expenses-app/pkg/app/authenticating"
	"expenses-app/pkg/app/health"
	"expenses-app/pkg/app/managing"
	"expenses-app/pkg/app/querying"
	"expenses-app/pkg/app/tracking"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const HELP_MSG string = `
Check the menu for available commands, please.

/n26importer - Start the N26 tracking importer.
/categories - Sends you all the categories available
/ping - Checks bot availability and health 
/help - Displays this menu
`
const DEFAULT_MSG string = "I don't get it. I'm no ChatGPT 🤖\n"
const WELCOME_MSG string = `
"Hello! I'm your personal expense manager bot 🤖
I'll help you keep track of your spending and budget, so you can take control of your finances. 
Let's get started by adding your first expense. Need help with anything? Just type /help to see what I can do for you."
`

// Register the following commands in the botfather

// categories - Get a list of categories
// n26importer - Start the N26 tracking importer.
// ping - Check if the bot is working
// help - Display the help message

// Run start the Telegram expense bot
func Run(tbot *tgbotapi.BotAPI, h *health.Service, m *managing.Service, t *tracking.Service, a *authenticating.Service, q *querying.Service) {
	tbot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := tbot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message.IsCommand() {
			handleCommands(update, &updates, tbot, q, t, h)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, DEFAULT_MSG)
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			tbot.Send(msg)
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
	case "start":
		tbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, WELCOME_MSG))
	case "help":
		tbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, HELP_MSG))
	case "ping":
		ping(tbot, update, h)
	default:
	}
	return msg
}

func ping(tbot *tgbotapi.BotAPI, update tgbotapi.Update, h *health.Service) {
	tbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, h.Ping()))
}
