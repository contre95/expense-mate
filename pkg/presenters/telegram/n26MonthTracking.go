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

const N26_EXIT_TRACKER string = "Closing N26 Report Importer 📄\n --------"
const SKIP_EXP_1 string = "Skip"
const SKIP_EXP_2 string = "⏭"
const EDIT_PRICE_1 string = "✏️ 💶"
const EDIT_PERSON_1 string = "✏️ 👥" // This is different from EDIT_PRICE_1

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

func skipRow(row []string) bool {
	if strings.Contains(row[4], "Round-up") {
		return true
	}
	if row[3] == "Income" {
		return true
	}
	return false
}

func formatMessage(record []string, expNum int, chatID int64) tgbotapi.MessageConfig {
	// https://core.telegram.org/bots/api#formatting-options
	msgText := fmt.Sprintf(` %d ) Expense 💶:
<code>
<b>Type:</b>         %s
<b>Shop:</b>        %s
<b>Price:</b>        `+"€ %s"+`
<b>Date:</b>         %s
<b>Reference:</b>    %s
</code>

What category does it belong ?`,
		expNum,
		record[3],
		record[1],
		strings.Replace(record[5], ".", ",", -1),
		record[0], // Using the original date string not to fomratted again
		record[4], // Using the original date string not to fomratted again
	)
	msg := tgbotapi.NewMessage(chatID, msgText)
	msg.ParseMode = tgbotapi.ModeHTML
	return msg
}

func newCreateRequest(record []string, userName string) (*tracking.CreateExpenseReq, error) {
	// Create Expense Request for use case
	date, err := time.Parse("2006-01-02", record[0])
	if err != nil {
		return nil, err
	}
	record[0] = date.Format("2006-01-02") // Replace original record with parsed value
	price, err := strconv.ParseFloat(record[5], 64)
	if err != nil {
		return nil, err
	}
	record[5] = fmt.Sprintf("%.2f", float64(price*-1)) // Replace original record with parsed value no to be using the request model
	return &tracking.CreateExpenseReq{
		Price:    float64(price * -1),
		Currency: "Euro",
		Shop:     record[1],
		City:     "Barcelona",
		Date:     date,
		People:   globalBotConfig.PeopleUsers[userName],
	}, nil
}

func importN26Expenses(tbot *tgbotapi.BotAPI, updates *tgbotapi.UpdatesChannel, chatID int64, userName string, categories []string, file *tgbotapi.File, t *tracking.Service) error {
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
	// Iterate over the imported rows of the csv
	tbot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("There are %d expenses to process", len(csvData))))
	for i, record := range csvData[0:] {
		// Skip row
		userSkip := false
		if skipRow(record) {
			fmt.Println("Skipping", record)
			continue
		}
		// create CreateExpenseReq
		createReq, err := newCreateRequest(record, userName)
		if err != nil {
			return err
		}
		// Send message with row already using the response
		allkeys := append(categories, SKIP_EXP_2) // Add skip button aside from categories
		allkeys = append(allkeys, EDIT_PRICE_1)   // Add edit buttong
		allkeys = append(allkeys, EDIT_PERSON_1)  // Add edit buttong
		msg := formatMessage(record, i+1, chatID)
		msg.ReplyMarkup = setOneTimeKeyBoardMap(allkeys, 4)
		_, err = tbot.Send(msg)
		if err != nil {
			return err
		}

		// Iterate over updates
		// TODO: Maybe handle all these if/else in different function or make a switch case
		for productUpdate := range *updates {
			if contains(categories, productUpdate.Message.Text) { // The user has send a category
				createReq.Category = productUpdate.Message.Text
				tbot.Send(tgbotapi.NewMessage(chatID, "Category set ✅"))
				tbot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Now, please type the name of the product for %s:", record[1])))
				for productUpdate := range *updates {
					if productUpdate.Message == nil {
						continue
					} else {
						createReq.Product = productUpdate.Message.Text
						break
					}
				}
				break
			} else if contains([]string{SKIP_EXP_1, SKIP_EXP_2}, productUpdate.Message.Text) { // The user has skipped the Expense
				tbot.Send(tgbotapi.NewMessage(chatID, "Exepense skipped ⏭"))
				userSkip = true
				break
			} else if productUpdate.Message.Text == EDIT_PERSON_1 { // The user want's to edit the price
				peopleMsg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Who has made this Expense ? \n%s ", record[1]))
				peopleMsg.ReplyMarkup = setOneTimeKeyBoardMap(globalBotConfig.People, 4)
				tbot.Send(peopleMsg)
				for peopleUpdate := range *updates {
					if peopleUpdate.Message == nil {
						continue
					}
					if !contains(globalBotConfig.People, peopleUpdate.Message.Text) {
						tbot.Send(tgbotapi.NewMessage(chatID, "Invalid person, please try again."))
						tbot.Send(peopleMsg)
						continue
					}
					createReq.People = peopleUpdate.Message.Text
					peopleDoneMsg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s now marked as the expender", createReq.People))
					peopleDoneMsg.ReplyMarkup = setOneTimeKeyBoardMap(allkeys, 4)
					tbot.Send(peopleDoneMsg)
					tbot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Now, I'll need a the category for %s please.", createReq.Product)))
					break
				}
			} else if productUpdate.Message.Text == EDIT_PRICE_1 { // The user want's to edit the price
				tbot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Please give me the new price for %s", record[1])))
				for priceUpdate := range *updates {
					if priceUpdate.Message == nil {
						continue
					}
					p, err := strconv.ParseFloat(priceUpdate.Message.Text, 64)
					if err != nil {
						tbot.Send(tgbotapi.NewMessage(chatID, "Invalid price, please send me a float"))
						continue
					}
					createReq.Price = p
					priceDoneMsg := tgbotapi.NewMessage(chatID, fmt.Sprintf("New price € %.2f set for %s", createReq.Price, createReq.Product))
					priceDoneMsg.ReplyMarkup = setOneTimeKeyBoardMap(allkeys, 4)
					tbot.Send(priceDoneMsg)
					tbot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Now, I'll need a the category for %s please.", createReq.Product)))
					break
				}
			} else { // The user has set a wrong update
				tbot.Send(tgbotapi.NewMessage(chatID, "Please pick a proper category or send \"Skip\" in order to skip this expense."))
				continue
			}
		}

		//  Save category here
		if !userSkip {
			createResp, createErr := t.ExpenseCreator.Create(*createReq)
			if createErr != nil {
				tbot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Could not save Expense: %v", createErr)))
			} else {
				tbot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Expense saved 💾\n%s", createResp.ID)))
			}
		}

		// Check if it has finished importing csv
		if i >= len(csvData)-1 { // Will never be > than len(csvData)
			tbot.Send(tgbotapi.NewMessage(chatID, "Then end 📉 :)"))
			msg := tgbotapi.NewMessage(chatID, N26_EXIT_TRACKER)
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			tbot.Send(msg)
		} else {
			tbot.Send(tgbotapi.NewMessage(chatID, "Next ⬇️"))
		}
	}
	return nil
}

func n26MonthTracking(tbot *tgbotapi.BotAPI, update *tgbotapi.Update, updates *tgbotapi.UpdatesChannel, t *tracking.Service, q *querying.Service) {
	tbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Starting N26 Report 📄\n --------"))
	n26HelpMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please send me an N26 CSV file export. You can pick the date range and download it from <a href='https://app.n26.com/downloads'>here </a>")
	n26HelpMsg.ParseMode = tgbotapi.ModeHTML
	n26HelpMsg.DisableWebPagePreview = true
	tbot.Send(n26HelpMsg)
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
			err = importN26Expenses(tbot, updates, chatID, newUpdate.Message.Chat.UserName, categories, file, t)
			if err != nil {
				tbot.Send(tgbotapi.NewMessage(chatID, err.Error()))
				break
			}
		} else if newUpdate.Message != nil && newUpdate.Message.Text == "Exit" {
			msg := tgbotapi.NewMessage(chatID, N26_EXIT_TRACKER)
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			tbot.Send(msg)
			break
		} else {
			tbot.Send(tgbotapi.NewMessage(newUpdate.Message.Chat.ID, "Please send me an N26 CSV file export. Otherwise send \"Exit\" "))
			continue // This can be removed
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
