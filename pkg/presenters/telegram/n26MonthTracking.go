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

func setOneTimeKeyBoardMap(items []string, rowsCant int) tgbotapi.ReplyKeyboardMarkup {
	matrix := [][]tgbotapi.KeyboardButton{}
	row := []tgbotapi.KeyboardButton{}
	for _, category := range items {
		newButton := tgbotapi.NewKeyboardButton(category)
		row = append(row, newButton)
		if len(row) == rowsCant || len(items)-len(matrix)*rowsCant == len(row) {
			matrix = append(matrix, row)
			row = []tgbotapi.KeyboardButton{}
		}
	}
	return tgbotapi.NewOneTimeReplyKeyboard(matrix...)
}

func n26MonthTracking(tbot *tgbotapi.BotAPI, update *tgbotapi.Update, updates *tgbotapi.UpdatesChannel, t *tracking.Service, q *querying.Service) {
	tbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Starting N26 Report 📄"))
	tbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Please send me an N26 csv file export"))
	for newUpdate := range *updates {
		if newUpdate.Message.Document != nil {
			fileID := newUpdate.Message.Document.FileID
			file, err := tbot.GetFile(tgbotapi.FileConfig{FileID: fileID})
			if err != nil {
				tbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
				break
			}
			fmt.Printf("Received file '%s' with ID %s", file.FilePath, fileID)
			url := "https://api.telegram.org/file/bot" + tbot.Token + "/" + file.FilePath
			tbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("You can download the file here: %s", url)))
			// TODO: Move this into a new function for better readability.
			cg := q.CategoryQuerier
			resp, err := cg.Query()
			if err != nil {
				tbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
			}
			category_keys := []string{}
			for _, name := range resp.Categories {
				category_keys = append(category_keys, name)
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Pick a category")
			msg.ReplyMarkup = setOneTimeKeyBoardMap(category_keys, 4)
            tbot.Send(msg)
			// END TODO
			break
		} else {
			tbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Closing N26 Report Importer 📄"))
			tbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Please send me an N26 csv file export. N26 Monlty tracker stopped."))
			break
		}
	}
}
