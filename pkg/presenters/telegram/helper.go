package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func setOneTimeKeyBoardMap(items []string, rowsCant int) tgbotapi.ReplyKeyboardMarkup {
	matrix := [][]tgbotapi.KeyboardButton{}
	row := []tgbotapi.KeyboardButton{}
	for _, category := range items {
		newButton := tgbotapi.NewKeyboardButton(category)
		row = append(row, newButton)
		if len(row) == rowsCant || len(items)-len(matrix)*rowsCant == len(row) {
			matrix = append(matrix, row)
			row = []tgbotapi.KeyboardButton{}
		}
	}
	return tgbotapi.NewOneTimeReplyKeyboard(matrix...)
}

func getKeyboardMarkup(items []string, rowsCount int) tgbotapi.ReplyKeyboardMarkup {
	matrix := [][]tgbotapi.KeyboardButton{}
	row := []tgbotapi.KeyboardButton{}
	for _, i := range items {
		newButton := tgbotapi.NewKeyboardButton(i)
		row = append(row, newButton)
		if len(row) == rowsCount || len(items)-len(matrix)*rowsCount == len(row) {
			matrix = append(matrix, row)
			row = []tgbotapi.KeyboardButton{}
		}
	}
	return tgbotapi.NewOneTimeReplyKeyboard(matrix...)
}
