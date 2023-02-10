package telegram

import (
	"encoding/csv"
	"expenses-app/pkg/app/querying"
	"expenses-app/pkg/app/tracking"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DocumentConfig struct {
	ChannelUsername string
	ChatID          int64
	MessageID       int
}

const EXIT_COMMAND = "N26TrackingClose"

func getCategories(tbot *tgbotapi.BotAPI, chatID int64, q *querying.Service) ([]string, error) {
	cg := q.CategoryQuerier
	resp, err := cg.Query()
	if err != nil {
		return nil, err
	}
	categories := []string{}
	for _, name := range resp.Categories {
		categories = append(categories, name)
	}
	// msg := tgbotapi.NewMessage(chatID, "Pick a category")
	// tbot.Send(msg)
	return categories, nil
}

func categorizeAndCreate(tbot *tgbotapi.BotAPI, updates *tgbotapi.UpdatesChannel, chatID int64, userName string, categories []string, file *tgbotapi.File, t *tracking.Service) error {
	// Download the file and read it
	response, err := http.Get(file.Link(tbot.Token))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	csvData, err := csv.NewReader(response.Body).ReadAll()
	if err != nil {
		return nil
	}
	// Fetch the categories
	tbot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("There are %d expenses to process", len(csvData))))
	for _, record := range csvData[1:] {
		date, err := time.Parse("2006-01-02", record[0])
		if err != nil {
			return err
		}
		price, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			return err
		}
		resp := tracking.CreateExpenseReq{Product: "Product", Price: float64(price * -1), Currency: "Euro", Place: record[1], City: "Barcelona", Date: date, People: userName, Category: "Categ"}
		// https://core.telegram.org/bots/api#formatting-options
		msgText := fmt.Sprintf("<b>Place:</b>          %s\n<b>Price:</b>           <code>€ %s</code>\n<b>People:</b>       %s\n<b>Date:</b>           %s\n\n What category does it belong ?",
			strings.Replace(fmt.Sprintf("%.2f", resp.Price), ".", ",", -1),
			resp.Place,
			record[0], // Using the original date string not to fomratted again
			resp.People,
		)
		msg := tgbotapi.NewMessage(chatID, msgText)
		msg.ParseMode = tgbotapi.ModeHTML
		msg.ReplyMarkup = setOneTimeKeyBoardMap(categories, 4)
		_, err = tbot.Send(msg)
		for update := range *updates {
          // check if in category
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func n26MonthTracking(tbot *tgbotapi.BotAPI, update *tgbotapi.Update, updates *tgbotapi.UpdatesChannel, t *tracking.Service, q *querying.Service) {
	tbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Starting N26 Report 📄"))
	tbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Please send me an N26 csv file export"))
	// Iterate over the new newUpdates
	for newUpdate := range *updates {
		chatID := newUpdate.Message.Chat.ID
		if newUpdate.Message.Document != nil && (newUpdate.Message.Document.MimeType == "text/comma-separated-values" || newUpdate.Message.Document.MimeType == "text/csv") {
			file, err := readSentFile(tbot, chatID, newUpdate)
			if err != nil {
				tbot.Send(tgbotapi.NewMessage(chatID, err.Error()))
				break
			}
			categories, errGet := getCategories(tbot, chatID, q)
			if errGet != nil {
				tbot.Send(tgbotapi.NewMessage(chatID, err.Error()))
				break
			}
			err = categorizeAndCreate(tbot, updates, chatID, newUpdate.Message.Chat.UserName, categories, file, t)
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
