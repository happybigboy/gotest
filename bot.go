package main

import (
	"fmt"
	"log"
	"main/models"
	"main/utils"
	"regexp"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Marzban URL (constant)
const Marzban_URL = "http://localhost:8000"

// State management
type UserState struct {
	mu     sync.Mutex
	states map[int64]string
}

func NewUserState() *UserState {
	return &UserState{
		states: make(map[int64]string),
	}
}

func (us *UserState) SetState(chatID int64, state string) {
	us.mu.Lock()
	defer us.mu.Unlock()
	us.states[chatID] = state
}

func (us *UserState) GetState(chatID int64) string {
	us.mu.Lock()
	defer us.mu.Unlock()
	return us.states[chatID]
}

func (us *UserState) ResetState(chatID int64) {
	us.mu.Lock()
	defer us.mu.Unlock()
	delete(us.states, chatID)
}

var (
	db        *gorm.DB
	userState *UserState
)

// Custom error types
type AuthError struct {
	Message string
}

func (e *AuthError) Error() string {
	return e.Message
}

type NetworkError struct {
	Message string
	Err     error
}

func (e *NetworkError) Error() string {
	return fmt.Sprintf("network error: %s - %v", e.Message, e.Err)
}

type JSONParseError struct {
	Message string
	Err     error
}

func (e *JSONParseError) Error() string {
	return fmt.Sprintf("json parse error: %s - %v", e.Message, e.Err)
}

// Centralized error handler (logs and sends message to user)
func handleError(bot *tgbotapi.BotAPI, chatID int64, err error) {
	// Log the error (could be sent to a logging service or file)
	log.Printf("Error: %v", err)

	// Inform the user about the error
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("An error occurred: %v", err))
	bot.Send(msg)
}

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("user_data.db"), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}

	// Migrate the schema
	db.AutoMigrate(&models.User{})

	userState = NewUserState()

	bot, err := tgbotapi.NewBotAPI("7362762333:AAF0KMRRjtvea7KDeyzuiscbA-9_Z7i4IQo")
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
			handleMessage(bot, update.Message)
		}
	}
}

func handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	if message.IsCommand() {
		switch message.Command() {
		case "start":
			handleStart(bot, message)

		case "help":
			handleHelp(bot, message)
		default:
			bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Unknown command. Use /help for available commands."))
		}
	} else {
		handleNormalMessage(bot, message)
	}
}

func handleStart(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	userState.SetState(message.Chat.ID, "awaiting_username")
	msg := tgbotapi.NewMessage(message.Chat.ID, "Welcome! Please enter your username (only letters and numbers are allowed):")
	bot.Send(msg)
}

func handleHelp(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	helpText := `Available commands:
/start - Begin setup or reset your username
/help - Show this help message`
	msg := tgbotapi.NewMessage(message.Chat.ID, helpText)
	bot.Send(msg)
}

func handleNormalMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	state := userState.GetState(message.Chat.ID)
	if state == "awaiting_username" {
		handleUsernameInput(bot, message)
	} else if state == "awaiting_password" {
		handlePasswordInput(bot, message)
	} else {
		msg := tgbotapi.NewMessage(message.Chat.ID, "I'm not sure what you mean. Use /help for available commands.")
		bot.Send(msg)
	}
}

func handleUsernameInput(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	// Extract username from the message text
	username := message.Text

	// Validate username
	if isValidUsername(username) {
		userState.SetState(message.Chat.ID, "awaiting_password") // Update state to await password
		models.CreateUser(db, message.Chat.ID, message.Text, "","")
		msg := tgbotapi.NewMessage(message.Chat.ID, "Username valid. Please enter a password:")
		bot.Send(msg)

	} else {
		// Inform the user about the invalid username
		msg := tgbotapi.NewMessage(message.Chat.ID, "Invalid username. Please use only letters and numbers. ")
		bot.Send(msg)
	}
}

func handlePasswordInput(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	password := message.Text
	if isValidPassword(password) {
		user, err := models.ReadUser(db, message.Chat.ID)
		models.ModifyUser(db,message.Chat.ID,user.Username,message.Text,user.Token)
		if err != nil {
			log.Panic("NO user")
		}
		// Try to get the access token with the constant URL
		accessToken, err := utils.GetAccessToken(Marzban_URL,user.Username, message.Text)
		if err != nil {
			switch err.(type) {
			case *AuthError:
				userState.SetState(message.Chat.ID, "awaiting_username")
				userState.ResetState(message.Chat.ID)
				msg := tgbotapi.NewMessage(message.Chat.ID, "Invalid username or password. Please start again with /start.")
				bot.Send(msg)
			default:
				handleError(bot, message.Chat.ID, err) // Use centralized error handler
			}
		} else {
			// Successfully got access token
			userState.ResetState(message.Chat.ID)
			msg := tgbotapi.NewMessage(message.Chat.ID, "Access token obtained successfully: "+accessToken)
			user , e := models.ReadUser(db,message.Chat.ID)
			if e != nil {
				log.Panic("No User")
			}
			models.ModifyUser(db,message.Chat.ID,user.Username,user.Password,accessToken)
			bot.Send(msg)
		}
	} else {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Invalid password. Please try again.")
		bot.Send(msg)
	}
}

func isValidUsername(username string) bool {
	match, _ := regexp.MatchString("^[a-zA-Z0-9]+$", username)
	return match
}

func isValidPassword(password string) bool {
	// Add password validation logic here, e.g., minimum length
	return len(password) >= 1
}
