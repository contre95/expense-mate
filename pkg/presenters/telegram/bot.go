package telegram

import (
	"expenses-app/pkg/app/analyzing"
	"expenses-app/pkg/app/health"
	"expenses-app/pkg/app/managing"
	"expenses-app/pkg/app/querying"
	"expenses-app/pkg/app/tracking"
	"expenses-app/pkg/gateways/ai"
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
/guess - Analyzes an image and prompts you the expenses for you to save.
/new - Creates a new expense. /fix if you made made a mistake.
/ping - Checks bot availability and health.
/help - Displays this menu.
`

const NOT_ALLOWED_MSG string = "I have keys but open no locks, I have space but no room. You can enter but not go outside, What am I?"
const ANSWER string = "keyboard"

type ControlSignal int

const (
	Continue ControlSignal = iota
	Stop
)

type UserSession struct {
	State string
	// Add more fields as needed to manage the conversation state
}

type Bot struct {
	API          *tgbotapi.BotAPI
	AllowedUsers []string
}

// Run starts the Telegram expense bot
func (b *Bot) Run(tbot *tgbotapi.BotAPI, receives, sends chan string, h *health.Service, t *tracking.Service, q *querying.Service, m *managing.Service, a *analyzing.Service, ai *ai.Guesser) {
	tbot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := tbot.GetUpdatesChan(u)
	var mu sync.Mutex
	err := b.updateAllowedUsers(m)
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
				go b.checkUpdates(done, updates, tbot, h, t, q, a, m, ai, &b.AllowedUsers, &mu)
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
			err := b.updateAllowedUsers(m)
			if err != nil {
				fmt.Println("Couldn't update allowed users:", err)
			} else {
				fmt.Println("Allowed users updated.")
			}
		case "getAllowedUsers":
			sends <- strings.Join(b.AllowedUsers, ", ")
		default:
			fmt.Println("Unknown command:", command)
		}
	}

}

func (b *Bot) checkUpdates(ImDone chan bool, updates tgbotapi.UpdatesChannel, tbot *tgbotapi.BotAPI, h *health.Service, t *tracking.Service, q *querying.Service, a *analyzing.Service, m *managing.Service, ai *ai.Guesser, allowedUsernames *[]string, mu *sync.Mutex) {
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
				if strings.Contains(update.Message.Text, ANSWER) {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Yeah.. it's a '%s', %s. But I'm still not letting you in.", ANSWER, update.Message.Chat.UserName))
				}
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				tbot.Send(msg)
				continue
			}
			switch update.Message.Text {
			case "/categories":
				listCategories(tbot, &update, q)
			case "/guess":
				guessExpense(tbot, &update, &updates, t, m, ai, update.Message.Chat.UserName)
			case "/new":
				createExpense(tbot, &update, &updates, t, q, m)
			case "/help":
				tbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, HELP_MSG))
			case "/summary":
				lastMonthSummary(tbot, &update, a)
			case "/unknown":
				categorizeUnknowns(tbot, &update, &updates, t, q, m, update.Message.Chat.UserName)
			case "/ping":
				b.ping(update, h)
			default:
				tbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, HELP_MSG))
			}
		}
	}
}

func (b *Bot) ping(update tgbotapi.Update, h *health.Service) {
	b.API.Send(tgbotapi.NewMessage(update.Message.Chat.ID, h.Ping()))
}

func (b *Bot) updateAllowedUsers(m *managing.Service) error {
	resp, err := m.UserManager.List()
	if err != nil {
		return err
	}
	b.AllowedUsers = []string{}
	for _, u := range resp.Users {
		b.AllowedUsers = append(b.AllowedUsers, u.TelegramUsername)
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
