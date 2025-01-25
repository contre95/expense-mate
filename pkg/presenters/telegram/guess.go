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

	// First check if we're receiving the image directly
	if u.Message.Photo == nil {
		// Send image request
		msg = tgbotapi.NewMessage(chatID, "📸 Please send a receipt photo for analysis")
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		tbot.Send(msg)

		// Wait for image response
		for {
		waitForImage:
			update := <-*uc
			if update.Message == nil {
				continue
			}

			// Check if message is from the same user
			if update.Message.Chat.UserName != username {
				tbot.Send(tgbotapi.NewMessage(
					update.Message.Chat.ID,
					fmt.Sprintf("Please wait for @%s to finish their current operation", username),
				))
				goto waitForImage
			}

			// Check if user wants to cancel
			if update.Message.Text == "/cancel" {
				msg = tgbotapi.NewMessage(chatID, "🚫 Expense guessing canceled")
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				tbot.Send(msg)
				return
			}

			// Check if received image
			if update.Message.Photo == nil {
				msg = tgbotapi.NewMessage(chatID, "🖼️ Please send a receipt photo or /cancel to abort")
				tbot.Send(msg)
				goto waitForImage
			}

			// Update the update reference with the image message
			u = &update
			break
		}
	}

	// Rest of your image processing logic...
	photo := u.Message.Photo[len(u.Message.Photo)-1]
	fileURL, err := tbot.GetFileDirectURL(photo.FileID)
	if err != nil {
		msg = tgbotapi.NewMessage(chatID, fmt.Sprintf("❌ Failed to get image: %v", err))
		tbot.Send(msg)
		return
	}

	// Download image
	resp, err := http.Get(fileURL)
	if err != nil {
		msg = tgbotapi.NewMessage(chatID, fmt.Sprintf("❌ Failed to download image: %v", err))
		tbot.Send(msg)
		return
	}
	defer resp.Body.Close()

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		msg = tgbotapi.NewMessage(chatID, fmt.Sprintf("❌ Failed to read image: %v", err))
		tbot.Send(msg)
		return
	}

	// Get AI guess
	shop, amount, date, product, err := aiGuesser.GuessExpense(imageData)
	if err != nil {
		msg = tgbotapi.NewMessage(chatID, fmt.Sprintf("🤖 AI processing failed: %v", err))
		tbot.Send(msg)
		return
	}
	// Find user ID
	userID := ""
	usersResp, err := m.UserManager.List()
	if err == nil {
		for id, u := range usersResp.Users {
			if username == u.TelegramUsername {
				userID = id
				break
			}
		}
	}

	// Create confirmation message
	expenseText := fmt.Sprintf(`📷 AI Guessed Expense:
<code>
Shop:    %s
Amount:  € %.2f
Date:    %s
Product: %s
</code>
Save this expense?`, shop, amount, date.Format("2006-01-02"), product)

	msg = tgbotapi.NewMessage(chatID, expenseText)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Skip"),
			tgbotapi.NewKeyboardButton("Save"),
		),
	)
	tbot.Send(msg)

	// Wait for user response
waitForResponse:
	update := <-*uc
	if update.Message.Chat.UserName != username {
		tbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
			fmt.Sprintf("Please wait for @%s to respond", username)))
		goto waitForResponse
	}

	switch update.Message.Text {
	case "Save":
		req := tracking.CreateExpenseReq{
			Amount:     amount,
			CategoryID: expense.UnkownCategoryID,
			Date:       date,
			Product:    product,
			Shop:       shop,
			UsersID:    []string{userID},
		}
		_, err := t.ExpenseCreator.Create(req)
		if err != nil {
			msg = tgbotapi.NewMessage(chatID, fmt.Sprintf("Failed to save: %v", err))
		} else {
			msg = tgbotapi.NewMessage(chatID, "✅ Expense saved with 'Unknown' category")
		}
	case "Skip":
		msg = tgbotapi.NewMessage(chatID, "⏭️ Expense skipped")
	default:
		msg = tgbotapi.NewMessage(chatID, "⚠️ Please choose Skip or Save")
		tbot.Send(msg)
		goto waitForResponse
	}

	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	tbot.Send(msg)
}
