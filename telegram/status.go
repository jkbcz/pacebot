package telegram

import (
	"fmt"
	"time"

	"github.com/JakubC-projects/pacebot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"
)

func (s *Service) SendStatusMessage(chatId int, content pacebot.StatusMessage) error {
	text, buttons := s.getStatusMessage(content)

	msg := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:      int64(chatId),
			ReplyMarkup: buttons,
		},
		Text:      text,
		ParseMode: tgbotapi.ModeHTML,
	}

	_, err := s.bot.Send(msg)

	return err
}

func (s *Service) EditStatusMessage(chatId int, messageId int, content pacebot.StatusMessage) error {
	text, buttons := s.getStatusMessage(content)

	msg := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      int64(chatId),
			MessageID:   messageId,
			ReplyMarkup: &buttons,
		},
		Text:      text,
		ParseMode: tgbotapi.ModeHTML,
	}

	_, err := s.bot.Send(msg)

	return err
}

func (s *Service) getStatusMessage(content pacebot.StatusMessage) (string, tgbotapi.InlineKeyboardMarkup) {
	userPercent := content.CurrentStatus / content.SeasonTarget * 100

	missingAmount := (content.MilestoneTarget - userPercent) * content.SeasonTarget / 100

	statusEmoji := "🟢"
	statusMessage := ""

	if missingAmount > 0 {
		statusEmoji = "🔴"
		statusMessage += fmt.Sprintf("\nBrakuje ci: <b>%.2f</b> %s\n", missingAmount, content.Currency)
	}

	text := fmt.Sprintf(`
Hej, tu Bob!

Dzisiaj mam już <b>%d%%</b> i jestem na dobrej drodze, żeby w tym miesiącu być On Track.

Cel na ten miesiąc: <b>%.2f%%</b>

A Ty? 😏

Twój Status: <b>%.2f%%</b> %s (%.2f / %.2f %s)
%s
<a href="%s">Zapisz się na Dugnad!</a>
<a href="%s">Wpłać na MyShare!</a>

Dane z: %s
`,
		content.AssistantPercentage,
		content.MilestoneTarget,
		userPercent,
		statusEmoji,
		content.CurrentStatus,
		content.SeasonTarget,
		content.Currency,
		statusMessage,
		content.RegisterURL,
		content.DonateURL,
		time.Now().Format("2006-01-02 15:04:05"),
	)

	var keyboard [][]tgbotapi.InlineKeyboardButton

	keyboard = append(keyboard, []tgbotapi.InlineKeyboardButton{
		{Text: "Odśwież", CallbackData: lo.ToPtr("show-status")},
	})

	if content.ShowNotifyAll {
		keyboard = append(keyboard, []tgbotapi.InlineKeyboardButton{
			{Text: "Powiadom Wszystkich", CallbackData: lo.ToPtr("notify-all")},
		})
	}

	keyboard = append(keyboard, []tgbotapi.InlineKeyboardButton{
		{Text: "Wyloguj się", URL: lo.ToPtr(content.LogoutURL)},
	})

	return text, tgbotapi.InlineKeyboardMarkup{InlineKeyboard: keyboard}
}
