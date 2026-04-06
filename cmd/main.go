package main

import (
	"blockchain-bot/internal/config"
	"context"
	"log"
	"log/slog"
	"os/signal"
	"strings"
	"syscall"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	cfg := config.Config{}
	tgBotCfg := cfg.Load(ctx)

	bot, err := api.NewBotAPI(tgBotCfg.TgToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = tgBotCfg.DebugMode

	slog.Info("connection established", "username", bot.Self.UserName)

	u := api.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for u := range updates {
		if u.Message == nil || !u.Message.IsCommand() {
			continue
		}

		if u.Message.Command() == "post" {
			text := strings.TrimSpace(u.Message.CommandArguments())
			if text == "" {
				reply(bot, u.Message.Chat.ID, "Usage: `/post your text`")
				continue
			}

			msg := api.NewMessage(tgBotCfg.TgChanId, text)
			msg.ParseMode = "" // Simply plain text.

			_, err := bot.Send(msg)
			if err != nil {
				log.Printf("Error in sending to channel: %v", err)
				reply(bot, u.Message.Chat.ID, "Error in publishing post")
			} else {
				reply(bot, u.Message.Chat.ID, "Post successfully published")
			}
		}
	}
}

func reply(bot *api.BotAPI, chatID int64, text string) {
	msg := api.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	bot.Send(msg)
}