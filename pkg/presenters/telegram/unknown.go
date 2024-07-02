package telegram

import (
	"expenses-app/pkg/app/querying"
	"expenses-app/pkg/app/tracking"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// 	msgText := fmt.Sprintf(` %d ) Expense 💶:
// <code>
// <b>Type:</b>         %s
// <b>Shop:</b>        %s
// <b>Price:</b>        `+"€ %s"+`
// <b>Date:</b>         %s
// <b>Reference:</b>    %s
// </code>
// What category does it belong ?`,

func categorizeUnknowns(tbot *tgbotapi.BotAPI, u *tgbotapi.Update, uc *tgbotapi.UpdatesChannel, t *tracking.Service, q *querying.Service) {
	chatID := u.Message.Chat.ID
	var msg tgbotapi.MessageConfig
	var update tgbotapi.Update
	var categoryID string

	// Fetch available categories
	categoryResp, err := q.CategoryQuerier.Query()
	if err != nil {
		msg = tgbotapi.NewMessage(chatID, fmt.Sprintf("Failed to fetch categories: %v", err))
		tbot.Send(msg)
		return
	}

	expensesReq := querying.ExpenseQuerierReq{
		Page:          0,
		MaxPageSize:   0,
		ExpenseFilter: querying.ExpenseQuerierFilter{ByCategoryID: []string{"unknown"}, ByShop: "", ByProduct: "", ByAmount: [2]uint{}, ByTime: [2]time.Time{}},
	}
	expensesResp, err := q.ExpenseQuerier.Query(expensesReq)
	if err != nil {
		msg = tgbotapi.NewMessage(chatID, fmt.Sprintf("Failed to fetch unknown expenses: %v", err))
		tbot.Send(msg)
		return
	}

	// Request category selection
	var reverseMap = map[string]string{}
	var categoryNames = []string{}
	for k, v := range categoryResp.Categories {
		reverseMap[v] = k
		categoryNames = append(categoryNames, v)
	}

	if len(expensesResp.Expenses) == 0 {
		tbot.Send(tgbotapi.NewMessage(chatID, "No expenses to categorize."))
		return
	}
	count := 0
	tbot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("You have a total of %d expenes.\n Send /done when tired of categorizing and continue in another moment.", len(expensesResp.Expenses))))
	for _, e := range expensesResp.Expenses {
		updateReq := tracking.UpdateExpenseReq{
			Amount:    e.Amount,
			Date:      e.Date,
			ExpenseID: e.ID,
			Product:   e.Product,
			Shop:      e.Shop,
		}
		for {
			expenseText := fmt.Sprintf(` %d ) Expense 💶:
<code>
Shop: %s
Amount: `+"€ %.2f"+`
Date: %s
Product: %s
</code>

What category does it belong ?
                `, count, e.Shop, e.Amount, e.Date.Format("2006-01-02"), e.Product)

			msg = tgbotapi.NewMessage(chatID, expenseText)
			msg.ParseMode = tgbotapi.ModeHTML
			msg.ReplyMarkup = getKeybaordMarkup(categoryNames, 2)
			tbot.Send(msg)
			update = <-*uc
			if update.Message.Text == "/done" {
				tbot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("We are done for now. You categorized %d/%d expenses.", count, len(expensesResp.Expenses))))
				return
			}
			categoryID = reverseMap[update.Message.Text]
			if _, exists := categoryResp.Categories[categoryID]; exists {
				updateReq.CategoryID = categoryID
				break
			}
			tbot.Send(tgbotapi.NewMessage(chatID, "Invalid category ID. Please try again."))
		}
		t.ExpenseUpdater.Update(updateReq)
		count++
	}

}
