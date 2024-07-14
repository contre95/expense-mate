package main

import (
	"log"
	"os"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Bot struct {
	bot       *tgbotapi.BotAPI
	token     string
	wg        sync.WaitGroup
	userChats map[int64]chan bool
	mux       sync.Mutex
}

func NewBot(token string) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		bot:       bot,
		token:     token,
		userChats: make(map[int64]chan bool),
	}, nil
}

func (b *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.Message != nil {
			go b.handleMessage(update.Message)
		}
	}
}

func (b *Bot) handleMessage(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID

	b.mux.Lock()
	defer b.mux.Unlock()

	// Check if there's already a goroutine running for this chat
	if _, ok := b.userChats[chatID]; !ok {
		b.userChats[chatID] = make(chan bool)
		go b.userRoutine(chatID)
	}

	// Send a response
	reply := tgbotapi.NewMessage(chatID, "Hello! I am processing your request.")
	b.bot.Send(reply)
}

func (b *Bot) userRoutine(chatID int64) {
	defer b.wg.Done()

	for {
		select {
		case <-b.userChats[chatID]:
			// Received a stop signal, exit the routine
			return
		}
	}
}

func main() {
	token := os.Getenv("TELEGRAM_APITOKEN")

	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN env var not set")
	}

	bot, err := NewBot(token)
	if err != nil {
		log.Fatalf("Error initializing bot: %s", err)
	}

	log.Printf("Authorized on account %s", bot.bot.Self.UserName)

	bot.Run()
}
