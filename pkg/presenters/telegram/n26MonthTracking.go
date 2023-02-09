package telegram

import (
	"encoding/csv"
	"expenses-app/pkg/app/querying"
	"expenses-app/pkg/app/tracking"
	"fmt"
	"net/http"

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
	response, err := http.Get(file.Link(tbot.Token))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	csvData, err := csv.NewReader(response.Body).ReadAll()
	if err != nil {
		return nil
	}
	newMsg := tgbotapi.NewMessage(chatID, "")
	for i, record := range csvData {
		newMsg.Text += fmt.Sprintf("%d - %s - %s - %s\n", i, record[0], record[1], record[4])
	}
	tbot.Send(newMsg)
	return nil
}

func n26MonthTracking(tbot *tgbotapi.BotAPI, update *tgbotapi.Update, updates *tgbotapi.UpdatesChannel, t *tracking.Service, q *querying.Service) {
	tbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Starting N26 Report 📄"))
	tbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Please send me an N26 csv file export"))
	// Iterate over the new newUpdates
	for newUpdate := range *updates {
		chatID := newUpdate.Message.Chat.ID
		if newUpdate.Message.Document != nil && newUpdate.Message.Document.MimeType == "text/comma-separated-values" {
			file, err := readSentFile(tbot, chatID, newUpdate)
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
			msg := tgbotapi.NewMessage(newUpdate.Message.Chat.ID, "Closing N26 Report Importer 📄")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			tbot.Send(msg)
			tbot.Send(tgbotapi.NewMessage(newUpdate.Message.Chat.ID, "Please send me an N26 csv file export. N26 Monlty tracker stopped."))
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
	msg.Text = fmt.Sprint("File received 👍🏽")
	tbot.Send(msg)
	return &file, nil
}
