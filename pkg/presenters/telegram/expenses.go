package telegram

import (
	"expenses-app/pkg/app/querying"
	"expenses-app/pkg/app/tracking"
	"fmt"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func createExpense(tbot *tgbotapi.BotAPI, u *tgbotapi.Update, uc *tgbotapi.UpdatesChannel, t *tracking.Service, q *querying.Service) {
	chatID := u.Message.Chat.ID

	var msg tgbotapi.MessageConfig
	var update tgbotapi.Update
	var userInput string

	// Helper function to wait for user input
	waitForResponse := func(prompt string, keyboard *tgbotapi.ReplyKeyboardMarkup) string {
		msg = tgbotapi.NewMessage(chatID, prompt)
		msg.ReplyMarkup = keyboard
		tbot.Send(msg)
		update = <-*uc
		if update.Message.Text == "/fix" {
			return "/fix"
		}
		return update.Message.Text
	}

	// Fetch available categories
	categoryResp, err := q.CategoryQuerier.Query()
	if err != nil {
		msg = tgbotapi.NewMessage(chatID, fmt.Sprintf("Failed to fetch categories: %v", err))
		tbot.Send(msg)
		return
	}

	// Create a keyboard for category selection
	var product, shop, categoryID string
	var amount float64
	var date time.Time

collectData:
	// Collect Product
	for {
		userInput = waitForResponse("Please enter the product name:", nil)
		if userInput != "/fix" {
			product = userInput
			break
		}
	}

	// Collect Amount
	for {
		userInput = waitForResponse("Please enter the amount:", nil)
		if userInput == "/fix" {
			goto collectData
		}
		amount, err = strconv.ParseFloat(userInput, 64)
		if err == nil {
			break
		}
		tbot.Send(tgbotapi.NewMessage(chatID, "Invalid amount. Please provide a valid number."))
	}

	// Collect Shop
	for {
		userInput = waitForResponse("Please enter the shop name:", nil)
		if userInput != "/fix" {
			shop = userInput
			break
		}
		goto collectData
	}

	// Collect Date
	for {
		userInput = waitForResponse("Please enter the date (YYYY-MM-DD):", nil)
		if userInput == "/fix" {
			goto collectData
		}
		date, err = time.Parse("2006-01-02", userInput)
		if err == nil {
			break
		}
		tbot.Send(tgbotapi.NewMessage(chatID, "Invalid date format. Please use YYYY-MM-DD."))
	}

	// Request category selection
	var reverseMap = map[string]string{}
	var categoryNames = []string{}
	for k, v := range categoryResp.Categories {
		reverseMap[v] = k
		categoryNames = append(categoryNames, v)
	}
	for {
		msg = tgbotapi.NewMessage(chatID, "Please choose a category:")
		msg.ReplyMarkup = getKeybaordMarkup(categoryNames, len(categoryNames)/2)
		tbot.Send(msg)
		update = <-*uc
		if update.Message.Text == "/fix" {
			goto collectData
		}
		categoryID = reverseMap[update.Message.Text]
		if _, exists := categoryResp.Categories[categoryID]; exists {
			break
		}
		tbot.Send(tgbotapi.NewMessage(chatID, "Invalid category ID. Please try again."))
	}

	// Create the request object
	req := tracking.CreateExpenseReq{
		Product:    product,
		Amount:     amount,
		Shop:       shop,
		Date:       date,
		CategoryID: categoryID,
	}

	// Create the expense
	resp, err := t.ExpenseCreator.Create(req)
	if err != nil {
		msg = tgbotapi.NewMessage(chatID, fmt.Sprintf("Failed to create expense: %v", err))
		tbot.Send(msg)
		return
	}

	// Send the success message to the user
	msg = tgbotapi.NewMessage(chatID, fmt.Sprintf("Expense created successfully with ID: %s", resp.ID))
	tbot.Send(msg)
}
