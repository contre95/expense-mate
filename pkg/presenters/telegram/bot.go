package telegram

import (
	"expenses-app/pkg/app/health"
	"expenses-app/pkg/app/querying"
	"expenses-app/pkg/app/tracking"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const HELP_MSG string = `
Check the menu for available commands, please.
/categories - Sends you all the categories available.
/summary - Sends summar of last month's expenses.
/ping - Checks bot availability and health.
/help - Displays this menu.
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

func isAllowed(chatID string, authorizedUsers []string, mu *sync.Mutex) bool {
	mu.Lock()
	defer mu.Unlock()
	for _, authorizedID := range authorizedUsers {
		if chatID == authorizedID {
			return true
		}
	}
	return false
}

// Run start the Telegram expense bot
func Run(tbot *tgbotapi.BotAPI, allowedUsers []string, commands chan string, botStatus *int32, h *health.Service, t *tracking.Service, q *querying.Service) {
	tbot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := tbot.GetUpdatesChan(u)
	var mu sync.Mutex
	for {
		select {
		case update := <-updates:
			if atomic.LoadInt32(botStatus) == 0 {
				continue
			}
			if isAllowed(update.Message.Chat.UserName, allowedUsers, &mu) {
				switch update.Message.Text {
				case "/categories":
					listCategories(tbot, &update, q)
				// case "n26importer":
				// 	n26MonthTracking(tbot, &update, &updates, t, q)
				case "/help":
					tbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, HELP_MSG))
				case "/summary":
					tbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "When implemented, I will send you a report of the past month"))
				case "/ping":
					ping(tbot, update, h)
				default:
					tbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, HELP_MSG))
				}
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, NOT_ALLOWED_MSG)
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				tbot.Send(msg)
			}
		case command := <-commands:
			switch command {
			case "start":
				fmt.Println(command) 
				atomic.StoreInt32(botStatus, 1)
				commands <- string("Started")
			case "stop":
				fmt.Println(command)
				atomic.StoreInt32(botStatus, 0)
				commands <- string("Stopped")
			case "getAllowedUsers":
				mu.Lock()
				commands <- strings.Join(allowedUsers, ",")
				mu.Unlock()
			case "getHelpMessage":
				// commands <- HELP_MSG
				commands <- fmt.Sprintf("Hola\nComo\n")
			default:
				continue
			}
		}
	}

}

func ping(tbot *tgbotapi.BotAPI, update tgbotapi.Update, h *health.Service) {
	tbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, h.Ping()))
}
