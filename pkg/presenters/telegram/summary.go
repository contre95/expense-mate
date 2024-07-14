package telegram

import (
	"expenses-app/pkg/app/analyzing"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func lastMonthSummary(tbot *tgbotapi.BotAPI, u *tgbotapi.Update, a *analyzing.Service) {
	chatID := u.Message.Chat.ID

	// Define time ranges
	now := time.Now().Add(time.Hour * 24)
	startOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	startOfLastMonth := startOfThisMonth.AddDate(0, -1, 0)
	endOfLastMonth := startOfThisMonth.Add(-time.Nanosecond)

	pastMonth := [2]time.Time{startOfLastMonth, endOfLastMonth}
	thisMonth := [2]time.Time{startOfThisMonth, now}

	// Function to summarize expenses
	summarizeExpenses := func(timeRange [2]time.Time) ([]analyzing.ExpensesSummary, error) {
		summaryReq := analyzing.ExpenseSummaryReq{
			TimeRange: timeRange,
		}
		summaryResp, err := a.ExpenseAnalyzer.Summarize(summaryReq)
		if err != nil {
			return nil, err
		}
		return summaryResp.Summaries, nil
	}

	// Get summarized expenses for past month
	pastMonthSummary, err := summarizeExpenses(pastMonth)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Failed to fetch past month expenses: %v", err))
		tbot.Send(msg)
		return
	}

	// Get summarized expenses for this month
	thisMonthSummary, err := summarizeExpenses(thisMonth)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Failed to fetch this month expenses: %v", err))
		tbot.Send(msg)
		return
	}

	// Calculate totals
	calculateTotal := func(summaries []analyzing.ExpensesSummary) (total float64, categories map[string]float64) {
		categories = make(map[string]float64)
		for _, summary := range summaries {
			total += summary.Total
			categories[summary.Category] += summary.Total
		}
		return
	}

	pastMonthTotal, pastMonthCategories := calculateTotal(pastMonthSummary)
	thisMonthTotal, thisMonthCategories := calculateTotal(thisMonthSummary)

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
