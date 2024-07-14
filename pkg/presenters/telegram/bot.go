package telegram

import (
	"expenses-app/pkg/app/health"
	"expenses-app/pkg/app/managing"
	"expenses-app/pkg/app/querying"
	"expenses-app/pkg/app/tracking"
	"fmt"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const HELP_MSG string = `
Check the menu for available commands, please.
/categories - Sends you all the categories available.
/summary - Sends summar of last month's expenses.
/unknown - Categorize imported expenses. /done and continue in another moment.
/new - Creates a new expense. /fix if you made made a mistake.
/ping - Checks bot availability and health.
/help - Displays this menu.
`

const NOT_ALLOWED_MSG string = "I speak without a mouth and hear without ears. I have no body, but I come alive with wind. What am I?"

type ControlSignal int

const (
	Continue ControlSignal = iota
	Stop
)

// Run starts the Telegram expense bot
func Run(tbot *tgbotapi.BotAPI, receives, sends chan string, h *health.Service, t *tracking.Service, q *querying.Service, m *managing.Service) {
	tbot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := tbot.GetUpdatesChan(u)
	var mu sync.Mutex
	allowedUsernames := []string{}
	err := updateAllowedUsers(&allowedUsernames, m)
	if err != nil {
		fmt.Println("Couldn't get allowed users:", err)
		return
	}

	running := false
	done := make(chan bool)
	for command := range receives {
		switch command {
		case "status":
			if running {
				sends <- "running"
			} else {
				sends <- "stopped"
			}
		case "start":
			fmt.Println("Starting new go routine")
			if !running {
				go checkUpdates(done, updates, tbot, h, t, q, m, &allowedUsernames, &mu)
				running = true
			}
		case "stop":
			fmt.Println("Attempting to stop go routine")
			if running {
				done <- true
				running = false
			}
			fmt.Println("Go routine stopped")
		case "updateAllowedUsers":
			err := updateAllowedUsers(&allowedUsernames, m)
			if err != nil {
				fmt.Println("Couldn't update allowed users:", err)
			} else {
				fmt.Println("Allowed users updated.")
			}
		case "getAllowedUsers":
			sends <- strings.Join(allowedUsernames, ", ")
		default:
			fmt.Println("Unknown command:", command)
		}
	}

}

func checkUpdates(ImDone chan bool, updates tgbotapi.UpdatesChannel, tbot *tgbotapi.BotAPI, h *health.Service, t *tracking.Service, q *querying.Service, m *managing.Service, allowedUsernames *[]string, mu *sync.Mutex) {
	fmt.Println("Go routine started")
	for {
		select {
		case <-ImDone:
			fmt.Println("Go routine stopped")
			return
		case update := <-updates:
			if update.Message == nil {
				continue
			}
			if !isAllowed(update.Message.Chat.UserName, allowedUsernames, mu) {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, NOT_ALLOWED_MSG)
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				tbot.Send(msg)
				continue
			}
			switch update.Message.Text {
			case "/categories":
				listCategories(tbot, &update, q)
			case "/new":
				createExpense(tbot, &update, &updates, t, q, m)
			case "/help":
				tbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, HELP_MSG))
			case "/summary":
				lastMonthSummary(tbot, &update, &updates, q)
			case "/unknown":
				categorizeUnknowns(tbot, &update, &updates, t, q, m, update.Message.Chat.UserName)
			case "/ping":
				ping(tbot, update, h)
			default:
				tbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, HELP_MSG))
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
