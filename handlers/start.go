package handlers

// import (
// 	"log"
// 	"fmt"
// 	"main/models" // Replace with your actual models package
// 	"gorm.io/gorm"
// 	"gorm.io/gorm/clause"
// 	"main/utils"
// 	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
// )

// func HandleStart(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *gorm.DB) {
// 	msg := tgbotapi.NewMessage(message.Chat.ID, "👋 Welcome! Please enter your username:")
// 	bot.Send(msg)

// 	// Prepare state data for upsert
// 	state := models.StateModel{
// 		UserID: int(message.Chat.ID),
// 		ChatID: int(message.Chat.ID),
// 		State:  "waiting_for_username", // Initial state when user starts
// 	}

// 	// Upsert user state using OnConflict
// 	if err := db.Clauses(clause.OnConflict{
// 		Columns:   []clause.Column{{Name: "user_id"}, {Name: "chat_id"}},
// 		DoUpdates: clause.AssignmentColumns([]string{"state"}), // Update state field if conflict
// 	}).Create(&state).Error; err != nil {
// 		log.Println("Failed to update or insert user state:", err)
// 	}
// }

// func HandleUsername(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *gorm.DB) {
// 	// Update the user's state to the username they provided
// 	state := models.StateModel{
// 		UserID: int(message.Chat.ID),
// 		ChatID: int(message.Chat.ID),
// 		State:  message.Text, // Store the username in the state field
// 	}

// 	// Upsert the state (Insert or update if conflict)
// 	if err := db.Clauses(clause.OnConflict{
// 		Columns:   []clause.Column{{Name: "user_id"}, {Name: "chat_id"}},
// 		DoUpdates: clause.AssignmentColumns([]string{"state"}), // Update the state field on conflict
// 	}).Create(&state).Error; err != nil {
// 		log.Println("Failed to update username:", err)
// 		return
// 	}

// 	msg := tgbotapi.NewMessage(message.Chat.ID, "🔒 Please enter your password:")
// 	bot.Send(msg)
// }

// func HandlePassword(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *gorm.DB) {
// 	// Retrieve the state for the given user to get the stored username
// 	var state models.StateModel
// 	if err := db.Where("chat_id = ?", message.Chat.ID).First(&state).Error; err != nil {
// 		log.Println("Failed to get user state:", err)
// 		return
// 	}

// 	// Get the password and attempt to authenticate
// 	password := message.Text
// 	Marzban_Url := "https://de.speedur.site:2053"
// 	token, err := utils.GetAccessToken(state.State, password,Marzban_Url) // Assuming state.State holds the username
// 	if err != nil {
// 		errorMsg := tgbotapi.NewMessage(message.Chat.ID, "❌ Login failed: "+err.Error()+" Please try again.")
// 		bot.Send(errorMsg)
// 		return
// 	}

// 	// Update token in the state model
// 	state.State = token // Assuming the state is used to store tokens (or create a dedicated field)

// 	// Use OnConflict to upsert the token
// 	if err := db.Clauses(clause.OnConflict{
// 		Columns:   []clause.Column{{Name: "user_id"}, {Name: "chat_id"}},
// 		DoUpdates: clause.AssignmentColumns([]string{"state"}), // Update token on conflict
// 	}).Create(&state).Error; err != nil {
// 		log.Println("Failed to store token:", err)
// 		return
// 	}

// 	// Confirm successful login
// 	msg := tgbotapi.NewMessage(message.Chat.ID, "✅ Login successful!")
// 	bot.Send(msg)

// 	// Send the main menu options
// 	mainMenu := getMainMenuKeyboard()
// 	bot.Send(mainMenu)
// }

// func HandleMenu(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *gorm.DB) {
// 	// Retrieve token from the state (assuming state.State holds the token)
// 	var state models.StateModel
// 	if err := db.Where("chat_id = ?", message.Chat.ID).First(&state).Error; err != nil {
// 		log.Println("Failed to get user state:", err)
// 		msg := tgbotapi.NewMessage(message.Chat.ID, "Session expired. Please log in again.")
// 		bot.Send(msg)
// 		return
// 	}

// 	token := state.State

// 	switch message.Text {
// 	case "🔍 Get User Info":
// 		msg := tgbotapi.NewMessage(message.Chat.ID, "Please enter the username you want to get info for:")
// 		bot.Send(msg)

// 	case "📋 Show Users":
// 		if token != "" {
// 			// Fetch users using the stored token
// 			Marzban_Url := "https://de.speedur.site:2053"
// 			users, err := utils.GetUsers(token, Marzban_Url, 0, 10, "")
// 			if err != nil {
// 				errorMsg := tgbotapi.NewMessage(message.Chat.ID, "Failed to get users: "+err.Error())
// 				bot.Send(errorMsg)
// 				return
// 			}

// 			// Format users into a string for displaying
// 			userList, err := formatUserList(users)
// 			if err != nil {
// 				errorMsg := tgbotapi.NewMessage(message.Chat.ID, "Failed to format users: "+err.Error())
// 				bot.Send(errorMsg)
// 				return
// 			}

// 			// Send the formatted user list
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
// // Helper function to format the user list from the map
// func formatUserList(users map[string]interface{}) (string, error) {
// 	// Extract users from the map
// 	userList, ok := users["users"].([]interface{})
// 	if !ok {
// 		return "", fmt.Errorf("unexpected format for users data")
// 	}

// 	// Format each user into a readable string
// 	var result string
// 	for i, user := range userList {
// 		userData, ok := user.(map[string]interface{})
// 		if !ok {
// 			return "", fmt.Errorf("unexpected format for user")
// 		}
// 		// Assuming userData has "username" and other fields
// 		username, _ := userData["username"].(string) // Adjust field names according to the actual response structure
// 		result += fmt.Sprintf("%d. %s\n", i+1, username)
// 	}

// 	return result, nil
// }

// func GetCurrentState(db *gorm.DB, chatID int64) string {
// 	// Query the user's state based on chatID
// 	var state models.StateModel
// 	if err := db.Where("chat_id = ?", chatID).First(&state).Error; err != nil {
// 		return "none" // Default state if user is not found
// 	}
// 	return state.State
// }

// func HandleMenuFallback(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *gorm.DB) {
// 	msg := tgbotapi.NewMessage(message.Chat.ID, "Please use the provided options.")
// 	bot.Send(msg)
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