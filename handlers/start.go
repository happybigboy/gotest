package handlers

import (
	"log"
	"main/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Handler for the /start command
func HandleStart(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Welcome! Send me a message and I will make a request for you.")
	bot.Send(msg)
}

// Handler for general user messages
func HandleUserMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	// Log the message
	log.Printf("[%s] %s", message.From.UserName, message.Text)

	// Make a GET request to the local server and send the response back
	response, err := utils.MakeRequest()
	if err != nil {
		log.Printf("Error making request: %v", err)
		errorMsg := tgbotapi.NewMessage(message.Chat.ID, "Something went wrong, please try again later.")
		bot.Send(errorMsg)
		return
	}

	// Send the JSON response back to the user
	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	bot.Send(msg)
}
