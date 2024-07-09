package telegram

import (
	"expenses-app/pkg/app/querying"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func lastMonthSummary(tbot *tgbotapi.BotAPI, u *tgbotapi.Update, uc *tgbotapi.UpdatesChannel, q *querying.Service) {
	chatID := u.Message.Chat.ID

	// Define time ranges
	now := time.Now().Add(time.Hour * 24)
	startOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	startOfLastMonth := startOfThisMonth.AddDate(0, -1, 0)
	endOfLastMonth := startOfThisMonth.Add(-time.Nanosecond)

	pastMonth := [2]time.Time{startOfLastMonth, endOfLastMonth}
	thisMonth := [2]time.Time{startOfThisMonth, now}

	// Function to get expenses
	getExpenses := func(timeRange [2]time.Time) ([]querying.ExpensesBasics, error) {
		expensesReq := querying.ExpenseQuerierReq{
			Page:        0,
			MaxPageSize: 1000, // assuming a high number to get all expenses
			ExpenseFilter: querying.ExpenseQuerierFilter{
				ByCategoryID: []string{},
				ByUsers:      []string{},
				ByShop:       "",
				ByProduct:    "",
				ByAmount:     [2]uint{},
				ByTime:       timeRange,
			},
		}
		expensesResp, err := q.ExpenseQuerier.Query(expensesReq)
		if err != nil {
			return nil, err
		}
		return expensesResp.Expenses, nil
	}

	// Get expenses for past month
	pastMonthExpenses, err := getExpenses(pastMonth)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Failed to fetch past month expenses: %v", err))
		tbot.Send(msg)
		return
	}

	// Get expenses for this month
	thisMonthExpenses, err := getExpenses(thisMonth)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Failed to fetch this month expenses: %v", err))
		tbot.Send(msg)
		return
	}

	// Calculate totals
	calculateTotal := func(expenses []querying.ExpensesBasics) (total float64, categories map[string]float64) {
		categories = make(map[string]float64)
		for _, expense := range expenses {
			total += expense.Amount
			categories[expense.Category.Name] += expense.Amount
		}
		return
	}

	pastMonthTotal, pastMonthCategories := calculateTotal(pastMonthExpenses)
	thisMonthTotal, thisMonthCategories := calculateTotal(thisMonthExpenses)

	// Format summary message
	summary := fmt.Sprintf("📊 Expenses summary:\n")
	summary += fmt.Sprintf("\n📅 Past Month:\n<strong>\n")
	for category, amount := range pastMonthCategories {
		summary += fmt.Sprintf(" - %s: <code>$%.2f</code>\n", category, amount)
	}
	summary += fmt.Sprintf("\nTotal</strong>: <code>$%.2f</code>", pastMonthTotal)
	summary += fmt.Sprintf("\n---------------------\n")
	summary += fmt.Sprintf("\n📅 This Month (so far):\n<strong>\n")
	for category, amount := range thisMonthCategories {
		summary += fmt.Sprintf("%s: <code>$%.2f</code>\n", category, amount)
	}
	summary += fmt.Sprintf("\nTotal</strong>: <code>$%.2f</code>", thisMonthTotal)
	summary += fmt.Sprintf("\n---------------------\n")

	// Send message
	msg := tgbotapi.NewMessage(chatID, summary)
	msg.ParseMode = tgbotapi.ModeHTML
	tbot.Send(msg)
}
