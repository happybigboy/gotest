package handlers

import (
	"fmt"
	"log"
	"main/models"
	"main/states"
	"main/utils"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	userState *states.UserState
)

const Marzban_Url = "http://localhost:8000"

// Centralized error handler (logs and sends message to user)
func handleError(bot *tgbotapi.BotAPI, chatID int64, err error) {
	log.Printf("Error: %v", err)
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("An error occurred: %v", err))
	bot.Send(msg)
}

func HandleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, userState *states.UserState) {
	if message.IsCommand() {
		switch message.Command() {
		case "start":
			handleStart(bot, message, userState)

		case "help":
			handleHelp(bot, message)

		case "menu":
			handleMenu(bot, message)

		default:
			bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Unknown command. Use /help for available commands."))
		}
		return // Stop further processing after handling commands
	}

	handleNormalMessage(bot, message, userState) // Handle normal messages
}

func handleStart(bot *tgbotapi.BotAPI, message *tgbotapi.Message, userState *states.UserState) {
	userState.SetState(message.Chat.ID, "awaiting_username")
	msg := tgbotapi.NewMessage(message.Chat.ID, "Welcome! Please enter your username (only letters and numbers are allowed):")
	bot.Send(msg) // No keyboard here
}

func handleNormalMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, userState *states.UserState) {
	state := userState.GetState(message.Chat.ID)

	if state == "awaiting_username" {
		handleUsernameInput(bot, message, userState)
	} else if state == "awaiting_password" {
		handlePasswordInput(bot, message, userState)
	} else {
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "I'm not sure what you mean. Use /help for available commands."))
	}
}

func handleUsernameInput(bot *tgbotapi.BotAPI, message *tgbotapi.Message, userState *states.UserState) {
	username := message.Text
	if isValidUsername(username) {
		userState.SetState(message.Chat.ID, "awaiting_password") // Update state to await password
		models.CreateUser(message.Chat.ID, username, "", "")
		msg := tgbotapi.NewMessage(message.Chat.ID, "Username valid. Please enter a password:")
		bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Invalid username. Please use only letters and numbers.")
		bot.Send(msg)
	}
}

func handlePasswordInput(bot *tgbotapi.BotAPI, message *tgbotapi.Message, userState *states.UserState) {
	password := message.Text
	if isValidPassword(password) {
		user, err := models.ReadUser(message.Chat.ID)
		if err != nil {
			log.Panic("NO user")
		}
		models.ModifyUser(message.Chat.ID, user.Username, password, user.Token)

		accessToken, err := utils.GetAccessToken(Marzban_Url, user.Username, password)
		if err != nil {
			handleError(bot, message.Chat.ID, err)
			return
		}

		userState.ResetState(message.Chat.ID)
		models.ModifyUser(message.Chat.ID, user.Username, password, accessToken)
		handleMenu(bot, message) // Redirect to menu after successful login
	} else {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Invalid password. Please try again.")
		bot.Send(msg)
	}
}

func handleMenu(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	user, err := models.ReadUser(message.Chat.ID)
	if err != nil || user.Token == "" {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Session expired or not logged in. Please use /start.")
		bot.Send(msg)
		return
	}

	// Send menu options if user is authenticated
	menuText := "Menu options:\n🔍 Get User Info\n📋 Show Users\n➕ Add User"
	msg := tgbotapi.NewMessage(message.Chat.ID, menuText)
	msg.ReplyMarkup = getMainMenuKeyboard() // Attach the inline keyboard
	bot.Send(msg)
}

func handleHelp(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	helpText := `Available commands:
	/start - Begin setup or reset your username
	/help - Show this help message
	/menu - Access the menu`
	msg := tgbotapi.NewMessage(message.Chat.ID, helpText)
	bot.Send(msg)
}

func getMainMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🔍 Get User Info"),
			tgbotapi.NewKeyboardButton("📋 Show Users"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("➕ Add User"),
		),
	)
	return keyboard
}

func isValidUsername(username string) bool {
	match, _ := regexp.MatchString("^[a-zA-Z0-9]+$", username)
	return match
}

func isValidPassword(password string) bool {
	return len(password) >= 1 // Adjust as needed
}
