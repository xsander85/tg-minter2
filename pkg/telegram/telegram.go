package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/xsander85/tg-minter2/pkg/config"
	"github.com/xsander85/tg-minter2/pkg/minter"
)

type Bot struct {
	bot    *tgbotapi.BotAPI
	ChatId []int64
}

func New(config *config.TelegramConf) *Bot {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = config.Debug

	return &Bot{bot: bot, ChatId: config.ChatId}
}

func (b *Bot) Run(v *minter.Valid) {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		message := update.Message
		// ignore any non-Message Updates
		if message == nil {
			continue
		}

		if status, _ := config.In_array(message.Chat.ID, b.ChatId); !status {
			b.Send(message.Chat.ID, "В данном чате нет возможности упралять данным ботом")
			continue
		}

		if message.IsCommand() {
			b.handleCommands(message, v)
			continue
		}
		b.handleMessage(message, v)
	}
}

func (b *Bot) Send(chatId int64, message string) bool {
	if chatId == 0 {
		chatId = b.ChatId[0]
	}
	msg := tgbotapi.NewMessage(chatId, message)
	msg.ParseMode = "markdown"
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Выключить Валидатор (тест)"),
		tgbotapi.NewKeyboardButton("Включить Валидатор (тест)"),
		tgbotapi.NewKeyboardButton("Данные по валидатору (тест)")))
	_, err := b.bot.Send(msg)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}
