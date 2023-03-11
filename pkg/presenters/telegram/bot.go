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
const NOT_ALLOWED_MSG string = "I speak without a mouth and hear without ears. I have no body, but I come alive with wind. What am I?"
const WELCOME_MSG string = `
"Hello! I'm your personal expense manager bot 🤖
I'll help you keep track of your spending and budget, so you can take control of your finances. 
Let's get started by adding your first expense. Need help with anything? Just type /help to see what I can do for you."
`

// const PEOPLE []string = []string{"Contre", "Anoux / Contre", "Anoux"}

// Register the following commands in the botfather

// categories - Get a list of categories
// n26importer - Start the N26 tracking importer.
// ping - Check if the bot is working
// help - Display the help message

type BotConfig struct {
	AllowedUsers []string
	People       []string
	PeopleUsers  map[string]string
	AuthUsers    []int64
}

func validConfig(c BotConfig) bool {
	for _, au := range c.AllowedUsers {
		if pe, _ := c.PeopleUsers[au]; !contains(c.People, pe) {
			return false
		}
	}
	return true
}

var globalBotConfig BotConfig // TODO: Probably this is not the best way. Mayby pass the condig all around or use some ctx

func isAuthorized(chatID int64, authorizedUsers []int64) bool {
	for _, authorizedID := range authorizedUsers {
		if chatID == authorizedID {
			return true
		}
	}
	return false
}

// Run start the Telegram expense bot
func Run(tbot *tgbotapi.BotAPI, botConfig BotConfig, h *health.Service, m *managing.Service, t *tracking.Service, a *authenticating.Service, q *querying.Service) {
	if !validConfig(botConfig) {
		panic("Wrong Telegram configuration")
	} else {
		globalBotConfig = botConfig
	}
	tbot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := tbot.GetUpdatesChan(u)
	for update := range updates {
		if isAuthorized(update.Message.Chat.ID, botConfig.AuthUsers) {
			if update.Message.IsCommand() {
				handleCommands(update, &updates, tbot, q, t, h)
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, DEFAULT_MSG)
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				tbot.Send(msg)
			}
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, NOT_ALLOWED_MSG)
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
