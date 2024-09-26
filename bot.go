package main

import (
	"log"
	"main/handlers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// Replace with your bot token
	botToken := "7362762333:AAF0KMRRjtvea7KDeyzuiscbA-9_Z7i4IQo"

	// Initialize the bot
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Create the update channel for receiving messages
	updates := bot.GetUpdatesChan(u)

	// Route updates to respective handlers
	for update := range updates {
		if update.Message == nil { // ignore non-Message updates
			continue
		}

		// Handle different commands or messages
		switch update.Message.Text {
		case "/start":
			handlers.HandleStart(bot, update.Message)
		default:
			handlers.HandleUserMessage(bot, update.Message)
		}
	}
}
