package telegram

import (
	"expenses-app/pkg/app/querying"
	"expenses-app/pkg/app/tracking"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DocumentConfig struct {
	ChannelUsername string
	ChatID          int64
	MessageID       int
}

const EXIT_COMMAND = "N26TrackingClose"

func displayCategories(tbot *tgbotapi.BotAPI, chatID int64, q *querying.Service) error {
	cg := q.CategoryQuerier
	resp, err := cg.Query()
	if err != nil {
		return err
	}
	categories := []string{}
	for _, name := range resp.Categories {
		categories = append(categories, name)
	}
	msg := tgbotapi.NewMessage(chatID, "Pick a category")
	msg.ReplyMarkup = setOneTimeKeyBoardMap(categories, 4)
	tbot.Send(msg)
	return nil
}

func categorizeAndCreate(tbot *tgbotapi.BotAPI, chatID int64, newUpdate tgbotapi.Update, file *tgbotapi.File, t *tracking.Service) error {
	return nil
}

func n26MonthTracking(tbot *tgbotapi.BotAPI, update *tgbotapi.Update, updates *tgbotapi.UpdatesChannel, t *tracking.Service, q *querying.Service) {
	tbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Starting N26 Report 📄"))
	tbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Please send me an N26 csv file export"))
	// Iterate over the new updates
	for newUpdate := range *updates {
		chatID := update.Message.Chat.ID
		if newUpdate.Message.Document != nil {
			file, err := readSentFile(tbot, chatID, *update)
			if err != nil {
				tbot.Send(tgbotapi.NewMessage(chatID, err.Error()))
				break
			}
			err = displayCategories(tbot, chatID, q)
			if err != nil {
				tbot.Send(tgbotapi.NewMessage(chatID, err.Error()))
				break
			}
			err = categorizeAndCreate(tbot, chatID, newUpdate, file, t)
			if err != nil {
				tbot.Send(tgbotapi.NewMessage(chatID, err.Error()))
				break
			}
		} else {
			tbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Closing N26 Report Importer 📄"))
			tbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Please send me an N26 csv file export. N26 Monlty tracker stopped."))
			break
		}
	}
}

func readSentFile(tbot *tgbotapi.BotAPI, chatID int64, newUpdate tgbotapi.Update) (*tgbotapi.File, error) {
	msg := tgbotapi.NewMessage(chatID, "")
	fileID := newUpdate.Message.Document.FileID
	file, err := tbot.GetFile(tgbotapi.FileConfig{FileID: fileID})
	if err != nil {
		return nil, err
	}
	fmt.Printf("Received file '%s' with ID %s", file.FilePath, fileID)
	url := "https://api.telegram.org/file/bot" + tbot.Token + "/" + file.FilePath
	msg.Text = fmt.Sprintf("You can download the file here: \n %s\n %s", file.Link(tbot.Token), url)
	tbot.Send(msg)
	return &file, nil
}
