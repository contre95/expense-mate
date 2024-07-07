package telegram

import (
	"expenses-app/pkg/app/health"
	"expenses-app/pkg/app/managing"
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
/unknown - Categorize unknown expenses. /done and continue in another moment.
/new - Creates a new expense. /fix if you made made a mistake.
/ping - Checks bot availability and health.
/help - Displays this menu.
`

const NOT_ALLOWED_MSG string = "I speak without a mouth and hear without ears. I have no body, but I come alive with wind. What am I?"

func isAllowed(chatID string, allowedUsernames *[]string, mu *sync.Mutex) bool {
	mu.Lock()
	defer mu.Unlock()
	for _, authorizedID := range *allowedUsernames {
		if chatID == authorizedID {
			return true
		}
	}
	return false
}

// Run start the Telegram expense bot
func Run(tbot *tgbotapi.BotAPI, commands chan string, botStatus *int32, h *health.Service, t *tracking.Service, q *querying.Service, m *managing.Service) {
	tbot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := tbot.GetUpdatesChan(u)
	var mu sync.Mutex
	allowedUsernames := []string{}
	err := updateAllowedUsers(&allowedUsernames, m)
	if err != nil {
		fmt.Println("Couldn't get allowed users")
	}
	for {
		select {
		case update := <-updates:
			if atomic.LoadInt32(botStatus) == 0 {
				continue
			}
			if isAllowed(update.Message.Chat.UserName, &allowedUsernames, &mu) {
				switch update.Message.Text {
				case "/categories":
					listCategories(tbot, &update, q)
				case "/new":
					createExpense(tbot, &update, &updates, t, q, m)
				case "/help":
					tbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, HELP_MSG))
				case "/summary":
					tbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "When implemented, I will send you a report of the past month"))
				case "/unknown":
					categorizeUnknowns(tbot, &update, &updates, t, q)
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
			case "updateAllowedUsers":
				mu.Lock()
				err = updateAllowedUsers(&allowedUsernames, m)
				mu.Unlock()
				if err != nil {
					fmt.Println(err)
					commands <- err.Error()
				}
				commands <- "Users updated."
			case "getAllowedUsers":
				mu.Lock()
				commands <- strings.Join(allowedUsernames, ",") // I'm passing a copy here
				mu.Unlock()
			default:
				continue
			}
		}
	}

}

func ping(tbot *tgbotapi.BotAPI, update tgbotapi.Update, h *health.Service) {
	tbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, h.Ping()))
}

func updateAllowedUsers(allowedUsernames *[]string, m *managing.Service) error {
	resp, err := m.UserManager.List()
	if err != nil {
		return err
	}
	*allowedUsernames = []string{}
	for _, u := range resp.Users {
		*allowedUsernames = append(*allowedUsernames, u.TelegramUsername)
	}
	fmt.Println("Usernames updated -----------------------------------------> ", resp.Users)
	return nil
}

func getKeybaordMarkup(items []string, rowsCant int) tgbotapi.ReplyKeyboardMarkup {
	matrix := [][]tgbotapi.KeyboardButton{}
	row := []tgbotapi.KeyboardButton{}
	for _, i := range items {
		newButton := tgbotapi.NewKeyboardButton(i)
		row = append(row, newButton)
		if len(row) == rowsCant || len(items)-len(matrix)*rowsCant == len(row) {
			matrix = append(matrix, row)
			row = []tgbotapi.KeyboardButton{}
		}
	}
	return tgbotapi.NewOneTimeReplyKeyboard(matrix...)
}

var NumberKeyboard = tgbotapi.NewOneTimeReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("1"), tgbotapi.NewKeyboardButton("2"), tgbotapi.NewKeyboardButton("3"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("4"), tgbotapi.NewKeyboardButton("5"), tgbotapi.NewKeyboardButton("6"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("7"), tgbotapi.NewKeyboardButton("8"), tgbotapi.NewKeyboardButton("9"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("0"), tgbotapi.NewKeyboardButton("."), tgbotapi.NewKeyboardButton("00"),
	),
)
