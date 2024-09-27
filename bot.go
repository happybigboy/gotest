package main

import (
	"log"
	"main/handlers"
	"main/models"
	"main/states"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	models.Main()
	states.NewUserState()

	bot, err := tgbotapi.NewBotAPI(os.Getenv("Telegram_Bot_Token"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			handlers.HandleMessage(bot, update.Message)
		}
	}
}
