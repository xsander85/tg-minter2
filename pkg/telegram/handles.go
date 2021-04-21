package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/xsander85/tg-minter2/pkg/minter"
)

func (b *Bot) handleMessage(message *tgbotapi.Message, v *minter.Valid) {
	b.handleMessageAll(message.Text, message.Chat.ID, v)
}

func (b *Bot) handleCommands(message *tgbotapi.Message, v *minter.Valid) {
	b.handleMessageAll(message.Command(), message.Chat.ID, v)
}

func (b *Bot) handleMessageAll(message string, chatId int64, v *minter.Valid) {
	switch message {

	case "start":
		b.Send(chatId, "message.Chat.ID"+fmt.Sprint(chatId))

	case "stop_validator", "Выключить Валидатор (тест)":
		res := v.SendTransactionOff(0)
		b.Send(chatId, res)
	case "start_validator", "Включить Валидатор (тест)":
		res := v.SendTransactionOn(0)
		b.Send(chatId, res)
	case "status_validator", "Данные по валидатору (тест)":
		res := v.StatusValidator(
			v.GetValidatorData(0))
		b.Send(chatId, res)

	default:
		b.Send(chatId, "Команда не найдена:"+message)
	}
}
