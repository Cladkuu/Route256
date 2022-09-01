package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gitlab.ozon.dev/astoyakin/route256/internal/bot/commandProcessor"
	"gitlab.ozon.dev/astoyakin/route256/internal/config"
	"log"
	"os"
)

type TelegramBot struct {
	botApi           *tgbotapi.BotAPI
	channel          tgbotapi.UpdatesChannel
	commandProcessor commandProcessor.ICommandProcessor
}

func NewTelegramBot(CommandProcessor commandProcessor.ICommandProcessor) (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv(config.TelegramApiKey))
	if err != nil {
		return nil, err
	}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		return nil, err
	}

	return &TelegramBot{botApi: bot,
		channel:          updates,
		commandProcessor: CommandProcessor}, nil

}

func (tb *TelegramBot) StartWorking() {
	var msg tgbotapi.MessageConfig
	for update := range tb.channel {
		ctx := context.Background()
		if update.Message.IsCommand() {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, tb.commandProcessor.ProcessCommand(ctx, update.Message.Command(), update.Message.CommandArguments()))
		} else {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "bot supported only commands\nUse /help to find out about commands")
		}
		if _, err := tb.botApi.Send(msg); err != nil {
			log.Print(err)
		}
	}
}
