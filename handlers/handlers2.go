package handlers

// "fmt"
// "log"
// "main/models"
// "main/states"
// "main/utils"
// "regexp"

// tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// const Marzban_Url = "http://localhost:8000"

// // Custom error types
// type AuthError struct {
// 	Message string
// }

// func (e *AuthError) Error() string {
// 	return e.Message
// }

// type NetworkError struct {
// 	Message string
// 	Err     error
// }

// func (e *NetworkError) Error() string {
// 	return fmt.Sprintf("network error: %s - %v", e.Message, e.Err)
// }

// type JSONParseError struct {
// 	Message string
// 	Err     error
// }

// func (e *JSONParseError) Error() string {
// 	return fmt.Sprintf("json parse error: %s - %v", e.Message, e.Err)
// }

// // Centralized error handler (logs and sends message to user)
// func handleError(bot *tgbotapi.BotAPI, chatID int64, err error) {
// 	// Log the error (could be sent to a logging service or file)
// 	log.Printf("Error: %v", err)

// 	// Inform the user about the error
// 	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("An error occurred: %v", err))
// 	bot.Send(msg)
// }

// func HandleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
// 	// Handle commands
// 	if message.IsCommand() {
// 		switch message.Command() {
// 		case "start":
// 			handleStart(bot, message, states.NewUserState())

// 		case "help":
// 			handleHelp(bot, message)

// 		default:
// 			bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Unknown command. Use /help for available commands."))
// 		}
// 		return // Stop further processing after handling commands
// 	}

// 	// Handle menu options
// 	user, err := models.ReadUser(message.Chat.ID)
// 	if err != nil {
// 		log.Println("No user found:", err)
// 		states.NewUserState().ResetState(message.Chat.ID)
// 		errorMsg := tgbotapi.NewMessage(message.Chat.ID, "Failed to get user: " + err.Error())
// 		bot.Send(errorMsg)
// 		return
// 	}

// 	token := user.Token
// 	switch message.Text {
// 	case "🔍 Get User Info":
// 		msg := tgbotapi.NewMessage(message.Chat.ID, "Please enter the username you want to get info for:")
// 		bot.Send(msg)

// 	case "📋 Show Users":
// 		if token != "" {
// 			users, err := utils.GetUsers(token, Marzban_Url, 0, 10, "")
// 			if err != nil {
// 				errorMsg := tgbotapi.NewMessage(message.Chat.ID, "Failed to get users: " + err.Error())
// 				bot.Send(errorMsg)
// 				return
// 			}

// 			userList, err := formatUserList(users)
// 			if err != nil {
// 				errorMsg := tgbotapi.NewMessage(message.Chat.ID, "Failed to format users: " + err.Error())
// 				bot.Send(errorMsg)
// 				return
// 			}

// 			bot.Send(tgbotapi.NewMessage(message.Chat.ID, userList))
// 		} else {
// 			msg := tgbotapi.NewMessage(message.Chat.ID, "Session expired. Please log in again.")
// 			bot.Send(msg)
// 		}

// 	case "➕ Add User":
// 		msg := tgbotapi.NewMessage(message.Chat.ID, "Please enter the username for the new user:")
// 		bot.Send(msg)

// 	default:
// 		msg := tgbotapi.NewMessage(message.Chat.ID, "Please use the provided options.")
// 		bot.Send(msg)
// 	}
// }

// func handleStart(bot *tgbotapi.BotAPI, message *tgbotapi.Message, userState *states.UserState) {
// 	userState.SetState(message.Chat.ID, "awaiting_username")

// 	msg := tgbotapi.NewMessage(message.Chat.ID, "Welcome! Please enter your username (only letters and numbers are allowed):")
// 	msg.ReplyMarkup = getMainMenuKeyboard() // Attach the inline keyboard
// 	bot.Send(msg)
// }

// func handleHelp(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
// 	helpText := `Available commands:
// 	/start - Begin setup or reset your username
// 	/help - Show this help message`
// 	msg := tgbotapi.NewMessage(message.Chat.ID, helpText)
// 	bot.Send(msg)
// }

// func handleNormalMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, userState *states.UserState) {
// 	state := userState.GetState(message.Chat.ID)
// 	if state == "awaiting_username" {
// 		handleUsernameInput(bot, message,states.NewUserState())
// 	} else if state == "awaiting_password" {
// 		handlePasswordInput(bot, message,states.NewUserState())
// 	} else {
// 		msg := tgbotapi.NewMessage(message.Chat.ID, "I'm not sure what you mean. Use /help for available commands.")
// 		bot.Send(msg)
// 	}
// }

// func handleUsernameInput(bot *tgbotapi.BotAPI, message *tgbotapi.Message, userState *states.UserState) {
// 	// Extract username from the message text
// 	username := message.Text

// 	// Validate username
// 	if isValidUsername(username) {
// 		userState.SetState(message.Chat.ID, "awaiting_password") // Update state to await password
// 		models.CreateUser(message.Chat.ID, message.Text, "", "")
// 		msg := tgbotapi.NewMessage(message.Chat.ID, "Username valid. Please enter a password:")
// 		bot.Send(msg)

// 	} else {
// 		// Inform the user about the invalid username
// 		msg := tgbotapi.NewMessage(message.Chat.ID, "Invalid username. Please use only letters and numbers. ")
// 		bot.Send(msg)
// 	}
// }

// func handlePasswordInput(bot *tgbotapi.BotAPI, message *tgbotapi.Message, userState *states.UserState) {
// 	password := message.Text
// 	if isValidPassword(password) {
// 		user, err := models.ReadUser(message.Chat.ID)
// 		models.ModifyUser(message.Chat.ID, user.Username, message.Text, user.Token)
// 		if err != nil {
// 			log.Panic("NO user")
// 		}
// 		// Try to get the access token with the constant URL
// 		accessToken, err := utils.GetAccessToken(Marzban_Url, user.Username, message.Text)
// 		if err != nil {
// 			switch err.(type) {
// 			case *AuthError:
// 				userState.SetState(message.Chat.ID, "awaiting_username")
// 				userState.ResetState(message.Chat.ID)
// 				msg := tgbotapi.NewMessage(message.Chat.ID, "Invalid username or password. Please start again with /start.")
// 				bot.Send(msg)
// 			default:
// 				handleError(bot, message.Chat.ID, err) // Use centralized error handler
// 			}
// 		} else {
// 			// Successfully got access token
// 			userState.ResetState(message.Chat.ID)
// 			msg := tgbotapi.NewMessage(message.Chat.ID, "Access token obtained successfully: "+accessToken)
// 			user, e := models.ReadUser(message.Chat.ID)
// 			if e != nil {
// 				log.Panic("No User")
// 			}
// 			models.ModifyUser(message.Chat.ID, user.Username, user.Password, accessToken)
// 			bot.Send(msg)
// 		}
// 	} else {
// 		msg := tgbotapi.NewMessage(message.Chat.ID, "Invalid password. Please try again.")
// 		bot.Send(msg)
// 	}
// }

// func getMainMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
// 	keyboard := tgbotapi.NewReplyKeyboard(
// 		tgbotapi.NewKeyboardButtonRow(
// 			tgbotapi.NewKeyboardButton("🔍 Get User Info"),
// 			tgbotapi.NewKeyboardButton("📋 Show Users"),
// 		),
// 		tgbotapi.NewKeyboardButtonRow(
// 			tgbotapi.NewKeyboardButton("➕ Add User"),
// 		),
// 	)
// 	return keyboard
// }

// func isValidUsername(username string) bool {
// 	match, _ := regexp.MatchString("^[a-zA-Z0-9]+$", username)
// 	return match
// }

// func isValidPassword(password string) bool {
// 	// Add password validation logic here, e.g., minimum length
// 	return len(password) >= 1
// }
