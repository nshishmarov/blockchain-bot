package main

import (
	"context"
	"log"
	"log/slog"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"blockchain-bot/internal/config"
	"blockchain-bot/internal/generator"
	"blockchain-bot/internal/store"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	slog.Info("started loading configuration")

	cfg := config.Config{}
	tgBotCfg := cfg.Load(ctx)

	store, err := store.NewFileStore(tgBotCfg.FileName, tgBotCfg.Path)
	if err != nil {
		log.Fatal(err)
	}

	g := generator.Generator{}

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

			data, err := store.GetLastData()
			if err != nil {
				log.Fatal(err)
			}

			s, newData, err := g.GeneratePost(*data, time.Now(), text)
			if err != nil {
				log.Fatal(err)
			}

			err = store.SetLastData(newData.PrevHash, newData.BlockNumber)
			if err != nil {
				log.Fatal(err)
			}

			msg := api.NewMessage(tgBotCfg.TgChanId, s)
			msg.ParseMode = "" // Simply plain text.

			_, err = bot.Send(msg)
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
