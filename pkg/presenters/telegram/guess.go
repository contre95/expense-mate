// guess.go
package telegram

import (
	"context"
	"expenses-app/pkg/app/managing"
	"expenses-app/pkg/app/querying"
	"expenses-app/pkg/app/tracking"
	"expenses-app/pkg/gateways/ollama"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func guessExpense(tbot *tgbotapi.BotAPI, u *tgbotapi.Update, uc *tgbotapi.UpdatesChannel, t *tracking.Service, q *querying.Service, m *managing.Service, o *ollama.OllamaAPI, username string) {
	chatID := u.Message.Chat.ID
	var msg tgbotapi.MessageConfig
	ctx := context.TODO() // Add context to function
	octx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	running, ollamaErr := o.IsRunning(octx)
	if !running || ollamaErr != nil {
		fmt.Println(ollamaErr)
		msg = tgbotapi.NewMessage(chatID, "⚠️ Failed to reach 🦙 Ollama.")
		tbot.Send(msg)
		return
	}
	// Handle initial request
	msg = tgbotapi.NewMessage(chatID, fmt.Sprintf("ℹ️ Info: Timeout set to ⏱️%.2f minutes", o.TimeOut.Minutes()))
	tbot.Send(msg)
	msg = tgbotapi.NewMessage(chatID, "Send a receipt photo 📸 or write 💬 the expense.")
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
	var guesses []ollama.ExpenseGuess
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

		guesses, err = o.GuessFromImage(imageData)
	} else {
		// Process text
		guesses, err = o.GuessFromText(u.Message.Text)
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

		// First prompt: Save, Skip, or Change Date
		msg := tgbotapi.NewMessage(chatID, expenseText)
		msg.ParseMode = tgbotapi.ModeHTML
		msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Skip"),
				tgbotapi.NewKeyboardButton("Change Date"),
				tgbotapi.NewKeyboardButton("Save"),
			),
		)
		tbot.Send(msg)

		// Wait for user response
		var response string
		var newDate time.Time
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
			case "Save", "Skip", "Change Date":
				response = update.Message.Text
			default:
				tbot.Send(tgbotapi.NewMessage(chatID, "⚠️ Please choose Skip, Change Date, or Save"))
			}
		}

		// Handle Change Date
		if response == "Change Date" {
			msg = tgbotapi.NewMessage(chatID, "📅 Enter the new date (YYYY-MM-DD, 'today', or 'yesterday'):")
			tbot.Send(msg)

			for {
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

				input := strings.ToLower(update.Message.Text)
				if input == "today" || input == "now" {
					newDate = time.Now()
					break
				} else if input == "y" || input == "yesterday" {
					newDate = time.Now().Add(-24 * time.Hour)
					break
				} else {
					var err error
					newDate, err = time.Parse("2006-01-02", update.Message.Text)
					if err == nil {
						break
					}
					tbot.Send(tgbotapi.NewMessage(chatID, "⚠️ Invalid date format. Use YYYY-MM-DD, 'today', or 'yesterday'"))
				}
			}
			response = "Save" // Proceed to save after changing date
		}

		// Handle Save (with or without date change)
		if response == "Save" {
			// Fetch available categories
			categoryResp, err := q.CategoryQuerier.Query()
			if err != nil {
				tbot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("❌ Failed to fetch categories: %v", err)))
				return
			}

			// Prepare category selection
			var categoryNames []string
			var reverseMap = map[string]string{}
			for id, name := range categoryResp.Categories {
				reverseMap[name] = id
				categoryNames = append(categoryNames, name)
			}

			// Prompt for category
			msg = tgbotapi.NewMessage(chatID, "📂 Choose a category:")
			msg.ReplyMarkup = getKeyboardMarkup(categoryNames, 3)
			tbot.Send(msg)

			var categoryID string
			for {
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

				if id, exists := reverseMap[update.Message.Text]; exists {
					categoryID = id
					break
				}
				tbot.Send(tgbotapi.NewMessage(chatID, "⚠️ Invalid category. Please choose from the list."))
			}

			// Prepare the expense request
			dateToUse := guess.Date
			if !newDate.IsZero() {
				dateToUse = newDate
			}

			req := tracking.CreateExpenseReq{
				Amount:     guess.Amount,
				CategoryID: categoryID,
				Date:       dateToUse,
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
