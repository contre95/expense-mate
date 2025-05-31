package telegram

import (
	"expenses-app/pkg/app/analyzing"
	"expenses-app/pkg/app/health"
	"expenses-app/pkg/app/managing"
	"expenses-app/pkg/app/querying"
	"expenses-app/pkg/app/tracking"
	"expenses-app/pkg/config"
	"expenses-app/pkg/gateways/ollama"
	"fmt"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const HELP_MSG string = `
Check the menu for available commands, please.
/categories - Sends you all the categories available.
/summary - Sends summary of last month's expenses.
/unknown - Categorize imported expenses. /done and continue in another moment.
/ai - Analyze image/text for expenses. Send /cancel to quit.
/new - Creates a new expense. /fix if you made a mistake.
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
	Config       *config.Config
}

type BotConfig struct {
	BotAPI       *tgbotapi.BotAPI
	Health       *health.Service
	Tracking     *tracking.Service
	Querying     *querying.Service
	Managing     *managing.Service
	Analyzing    *analyzing.Service
	AI           *ollama.OllamaAPI
	AllowedUsers *[]string
	Mu           *sync.Mutex
}

// Run starts the Telegram expense bot
func (b *Bot) Run(tbot *tgbotapi.BotAPI, receives, sends chan string, cfg BotConfig) {
	tbot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := tbot.GetUpdatesChan(u)

	err := b.updateAllowedUsers(cfg.Managing)
	if err != nil {
		fmt.Println("Couldn't get allowed users:", err)
		return
	}

	done := make(chan bool)
	go b.checkUpdates(done, updates, cfg)

	running := true
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
			err := b.updateAllowedUsers(cfg.Managing)
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

func (b *Bot) checkUpdates(done chan bool, updates tgbotapi.UpdatesChannel, botUseCases BotConfig) {
	fmt.Println("Go routine started")
	for {
		select {
		case <-done:
			fmt.Println("Go routine stopped")
			return
		case update := <-updates:
			if update.Message == nil {
				continue
			}
			if !isAllowed(update.Message.Chat.UserName, botUseCases.AllowedUsers, botUseCases.Mu) {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, NOT_ALLOWED_MSG)
				if strings.Contains(update.Message.Text, ANSWER) {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Yeah.. it's a '%s', %s. But I'm still not letting you in.", ANSWER, update.Message.Chat.UserName))
				}
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				botUseCases.BotAPI.Send(msg)
				continue
			}
			fmt.Println(b.Config.OllamaEnabled())
			switch update.Message.Text {
			case "/categories":
				listCategories(botUseCases.BotAPI, &update, botUseCases.Querying)
			case "/ai":
				if b.Config.OllamaEnabled() {
					guessExpense(botUseCases.BotAPI, &update, &updates, botUseCases.Tracking, botUseCases.Querying, botUseCases.Managing, botUseCases.AI, update.Message.Chat.UserName)
				} else {
					fmt.Println("🦙 Ollama API currently not enabled.")
					omsg := tgbotapi.NewMessage(update.Message.Chat.ID, "🦙 Ollama API currently not enabled. Please refer to the [docs](https://chat.deepseek.com) to enable it.")
					omsg.ParseMode = tgbotapi.ModeMarkdown
					omsg.DisableWebPagePreview = true
					botUseCases.BotAPI.Send(omsg)
				}
			case "/new":
				createExpense(botUseCases.BotAPI, &update, &updates, botUseCases.Tracking, botUseCases.Querying, botUseCases.Managing)
			case "/help":
				botUseCases.BotAPI.Send(tgbotapi.NewMessage(update.Message.Chat.ID, HELP_MSG))
			case "/summary":
				lastMonthSummary(botUseCases.BotAPI, &update, botUseCases.Analyzing)
			case "/unknown":
				categorizeUnknowns(botUseCases.BotAPI, &update, &updates, botUseCases.Tracking, botUseCases.Querying, botUseCases.Managing, update.Message.Chat.UserName)
			case "/ping":
				b.ping(update, botUseCases.Health)
			default:
				botUseCases.BotAPI.Send(tgbotapi.NewMessage(update.Message.Chat.ID, HELP_MSG))
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
