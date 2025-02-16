// guess.go
package telegram

import (
	"expenses-app/pkg/app/managing"
	"expenses-app/pkg/app/tracking"
	"expenses-app/pkg/domain/expense"
	"expenses-app/pkg/gateways/ai"
	"fmt"
	"io"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func guessExpense(tbot *tgbotapi.BotAPI, u *tgbotapi.Update, uc *tgbotapi.UpdatesChannel, t *tracking.Service, m *managing.Service, aiGuesser *ai.Guesser, username string) {
	chatID := u.Message.Chat.ID
	var msg tgbotapi.MessageConfig

	// Handle initial request
	msg = tgbotapi.NewMessage(chatID, "📸 Send a receipt photo or paste transaction text")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	tbot.Send(msg)

	// Wait for response
	for {
	waitForInput:
		update := <-*uc
		if update.Message == nil {
			continue
		}

		if update.Message.Chat.UserName != username {
			tbot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				fmt.Sprintf("Please wait for @%s to finish their current operation", username),
			))
			goto waitForInput
		}

		if update.Message.Text == "/cancel" {
			msg = tgbotapi.NewMessage(chatID, "🚫 Expense guessing canceled")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			tbot.Send(msg)
			return
		}

		if update.Message.Photo == nil && update.Message.Text == "" {
			msg = tgbotapi.NewMessage(chatID, "🖼️ Please send a receipt photo/text or /cancel to abort")
			tbot.Send(msg)
			goto waitForInput
		}

		u = &update
		break
	}

	// Process input
	var guesses []ai.ExpenseGuess
	var err error

	if u.Message.Photo != nil {
		// Process image
		photo := u.Message.Photo[len(u.Message.Photo)-1]
		fileURL, err := tbot.GetFileDirectURL(photo.FileID)
		if err != nil {
			tbot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("❌ Failed to get image: %v", err)))
			return
		}

		resp, err := http.Get(fileURL)
		if err != nil {
			tbot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("❌ Failed to download image: %v", err)))
			return
		}
		defer resp.Body.Close()

		imageData, err := io.ReadAll(resp.Body)
		if err != nil {
			tbot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("❌ Failed to read image: %v", err)))
			return
		}

		guesses, err = aiGuesser.GuessFromImage(imageData)
	} else {
		// Process text
		guesses, err = aiGuesser.GuessFromText(u.Message.Text)
	}

	if err != nil {
		tbot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("🤖 AI processing failed: %v", err)))
		return
	}

	if len(guesses) == 0 {
		tbot.Send(tgbotapi.NewMessage(chatID, "🤖 No expenses detected in the input"))
		return
	}

	// Find user ID
	userID := ""
	usersResp, err := m.UserManager.List()
	if err == nil {
		for id, u := range usersResp.Users {
			if u.TelegramUsername == username {
				userID = id
				break
			}
		}
	}
	if userID == "" {
		tbot.Send(tgbotapi.NewMessage(chatID, "❌ User not found in system"))
		return
	}

	// Process each guess
	for _, guess := range guesses {
		// Create confirmation message
		expenseText := fmt.Sprintf(`📷 AI Guessed Expense:
<code>
Shop:    %s
Amount:  € %.2f
Date:    %s
Product: %s
</code>
Save this expense?`,
			guess.Shop, guess.Amount, guess.Date.Format("2006-01-02"), guess.Product)

		msg := tgbotapi.NewMessage(chatID, expenseText)
		msg.ParseMode = tgbotapi.ModeHTML
		msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Skip"),
				tgbotapi.NewKeyboardButton("Save"),
			),
		)
		tbot.Send(msg)

		// Wait for user response
		var response string
		for response == "" {
			update := <-*uc
			if update.Message == nil || update.Message.Chat.UserName != username {
				continue
			}

			if update.Message.Text == "/cancel" {
				msg = tgbotapi.NewMessage(chatID, "🚫 Expense guessing canceled")
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				tbot.Send(msg)
				return
			}

			switch update.Message.Text {
			case "Save", "Skip":
				response = update.Message.Text
			default:
				tbot.Send(tgbotapi.NewMessage(chatID, "⚠️ Please choose Skip or Save"))
			}
		}

		// Handle response
		if response == "Save" {
			req := tracking.CreateExpenseReq{
				Amount:     guess.Amount,
				CategoryID: expense.UnknownCategoryID,
				Date:       guess.Date,
				Product:    guess.Product,
				Shop:       guess.Shop,
				UsersID:    []string{userID},
			}

			if _, err := t.ExpenseCreator.Create(req); err != nil {
				tbot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("❌ Failed to save: %v", err)))
			} else {
				tbot.Send(tgbotapi.NewMessage(chatID, "✅ Expense saved with 'Unknown' category"))
			}
		} else {
			tbot.Send(tgbotapi.NewMessage(chatID, "⏭️ Expense skipped"))
		}

		// Clear keyboard
		clearMsg := tgbotapi.NewMessage(chatID, "Processing next expense...")
		clearMsg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		tbot.Send(clearMsg)
	}

	// Final message
	finalMsg := tgbotapi.NewMessage(chatID, "🏁 All expenses processed")
	finalMsg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	tbot.Send(finalMsg)
}
