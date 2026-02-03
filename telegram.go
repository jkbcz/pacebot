package pacebot

type StatusMessage struct {
	CurrentStatus       float64
	SeasonTarget        float64
	AssistantPercentage int

	Currency string

	MilestoneTarget float64

	RegisterURL string
	DonateURL   string

	LogoutURL string

	ShowNotifyAll bool
}

type TelegramService interface {
	GetBotUrl() string
	SendWelcomeMessage(chatId int, loginUrl string) error
	SendStatusMessage(chatId int, status StatusMessage) error
	SendErrorMessage(chatId int, msg string) error
	EditStatusMessage(chatId int, messageId int, status StatusMessage) error
}
