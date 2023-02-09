package telegram

import (
	"expenses-app/pkg/app/tracking"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DocumentConfig struct {
	ChannelUsername string
	ChatID          int64
	MessageID       int
}

const EXIT_COMMAND = "N26TrackingClose"

func n26MonthTracking(msg *tgbotapi.MessageConfig, tbot *tgbotapi.BotAPI, update *tgbotapi.Update, updates *tgbotapi.UpdatesChannel, t *tracking.Service) {
	tbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Please send me an N26 csv file export"))
	for newUpdate := range *updates {
		if newUpdate.Message.Document != nil {
			fmt.Println("HOLAAAA2")
			fileID := newUpdate.Message.Document.FileID
			file, err := tbot.GetFile(tgbotapi.FileConfig{FileID: fileID})
			if err != nil {
				msg.Text = err.Error()
			}
			fmt.Printf("Received file '%s' with ID %s", file.FilePath, fileID)
			// You can download the file using the file path and URL:
			url := "https://api.telegram.org/file/bot" + tbot.Token + "/" + file.FilePath
			msg.Text = fmt.Sprintf("You can download the file here: %s", url)
            break
		} else {
			msg.Text = fmt.Sprintf("Please send me an N26 csv file export. N26 Monlty tracker stopped.")
			break
		}
	}
}
