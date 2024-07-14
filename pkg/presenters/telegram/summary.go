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

	pastMonthSummary, err := a.ExpenseAnalyzer.Summarize(analyzing.ExpenseSummaryReq{
		TimeRange: pastMonth,
	})
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Failed to fetch last month expenses: %v", err))
		tbot.Send(msg)
		return
	}
	thisMonthSummary, err := a.ExpenseAnalyzer.Summarize(analyzing.ExpenseSummaryReq{
		TimeRange: thisMonth,
	})
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Failed to fetch this month expenses: %v", err))
		tbot.Send(msg)
		return
	}

	// Format summary message
	summary := fmt.Sprintf("📊 Expenses summary:\n")
	summary += fmt.Sprintf("\n📅 %s:\n<strong>\n", pastMonthSummary.Month.String())
	for _, category := range pastMonthSummary.Summaries {
		summary += fmt.Sprintf(" - %s: <code>$%.2f</code>\n", category.Category, category.Total)
	}
	summary += fmt.Sprintf("\nTotal</strong>: <code>$%.2f</code>", pastMonthSummary.Total)
	summary += fmt.Sprintf("\n---------------------\n")
	summary += fmt.Sprintf("\n📅 %s (so far):\n<strong>\n", thisMonthSummary.Month.String())
	for _, category := range thisMonthSummary.Summaries {
		summary += fmt.Sprintf(" - %s: <code>$%.2f</code>\n", category.Category, category.Total)
	}
	summary += fmt.Sprintf("\nTotal</strong>: <code>$%.2f</code>", pastMonthSummary.Total)
	summary += fmt.Sprintf("\n---------------------\n")

	// Send message
	msg := tgbotapi.NewMessage(chatID, summary)
	msg.ParseMode = tgbotapi.ModeHTML
	tbot.Send(msg)
}
