package handlers

import (
	"fmt"
	"log"
	"main/models"
	"main/states"
	"main/utils"
	"math"
	"regexp"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// var (
// 	userState *states.UserState
// )
var userOffset = 0
const Marzban_Url = "http://localhost:8000"
const paddingChar = '\u00A0' // Non-breaking space

// Centralized error handler (logs and sends message to user)
func HandleCallbackQuery(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
    // Acknowledge the callback query
    answer := tgbotapi.NewCallback(callback.ID, "")
    if _, err := bot.Request(answer); err != nil {
        log.Printf("Failed to answer callback query: %v", err)
    }

    switch {
    case callback.Data == "next":
        userOffset++ // Ensure userOffset is managed correctly
        HandleButtonPress(bot, callback, userOffset)
        return
    case callback.Data == "back":
        if userOffset > 0 {
            userOffset--
            HandleButtonPress(bot, callback, userOffset)
        }
        return
    case strings.HasPrefix(callback.Data, "back_to_users_"):
        HandleButtonPress(bot, callback, userOffset)
        return
    case strings.HasPrefix(callback.Data, "modify_"):
        // Handle modify logic here
        return
    case strings.HasPrefix(callback.Data, "delete_"):
        // Handle delete logic here
        return
    case strings.HasPrefix(callback.Data, "change_subscription_link_"):
        username := strings.TrimPrefix(callback.Data, "change_subscription_link_")
        // Call a function to handle changing the subscription link
        handleRevokeSubscription(bot, callback.Message, username)
        return
    case callback.Data == "exit":
        editConfig := tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, "Hello there")
        if _, err := bot.Send(editConfig); err != nil {
            log.Printf("Failed to edit message: %v", err)
        }
        return
    default:
        re := regexp.MustCompile(`^user_(.+)$`)
        matches := re.FindStringSubmatch(callback.Data)
        if len(matches) > 1 {
            username := matches[1]
            handleUser(bot, callback.Message, username)
            return
        }
        HandleButtonPress(bot, callback, userOffset)
    }
}

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

	// Check if the message is a callback query
		if message.Text == "📋 Show Users" {
			handleShowUsers(bot, message)
			return
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

func handleShowUsers(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	user, err := models.ReadUser(message.Chat.ID)
	if err != nil || user.Token == "" {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Session expired or not logged in. Please use /start.")
		bot.Send(msg)
		return
	}

	// Call GetUsers with the user's token
	usersResponse, err := utils.GetUsers(user.Token, Marzban_Url, 0, 3, "") // Adjust the parameters as needed
	if err != nil {
		handleError(bot, message.Chat.ID, err)
		return
	}

	// Assuming usersResponse is a map with a "users" key containing the list of users
	usersInterface := usersResponse["users"].([]interface{}) // Type assert to slice of interfaces

	// Convert []interface{} to []map[string]interface{}
	var users []map[string]interface{}
	for _, userInterface := range usersInterface {
		userMap, ok := userInterface.(map[string]interface{})
		if ok {
			users = append(users, userMap)
		}
	}

	// Format the user list and send it as a message
	userList, keyboard := formatUserList(users,message.Chat.ID,userOffset,len(users))
	msg := tgbotapi.NewMessage(message.Chat.ID, userList)
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func HandleButtonPress(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery, offset int) {
	// Fetch user info based on callback
	user, err := models.ReadUser(callback.From.ID)
	if err != nil || user.Token == "" {
		msg := tgbotapi.NewMessage(callback.From.ID, "Session expired or not logged in. Please use /start.")
		bot.Send(msg)
		return
	}

	// Fetch users with updated offset
	usersResponse, err := utils.GetUsers(user.Token, Marzban_Url, offset, 3, "")
	if err != nil {
		handleError(bot, callback.Message.Chat.ID, err)
		return
	}

	// Assuming usersResponse is a map with a "users" key containing the list of users
	usersInterface, ok := usersResponse["users"].([]interface{})
	if !ok {
		msg := tgbotapi.NewMessage(callback.From.ID, "Failed to retrieve user list.")
		bot.Send(msg)
		return
	}

	// Convert []interface{} to []map[string]interface{}
	var users []map[string]interface{}
	for _, userInterface := range usersInterface {
		if userMap, ok := userInterface.(map[string]interface{}); ok {
			users = append(users, userMap)
		}
	}

	// Format the user list and create the inline keyboard
	// getallusers, err := utils.GetAllUsers(user.Token,Marzban_Url)
	userList, inlineKeyboard := formatUserList(users,callback.From.ID,userOffset,len(users))

	// Edit the message with the new user list and inline keyboard
	editMsg := tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, userList)
	editMsg.ReplyMarkup = &inlineKeyboard // Use a pointer to inlineKeyboard
	bot.Send(editMsg)
}

func formatUsername(username interface{}) string {
	// Convert to string and ensure it's 9 characters
	strUsername := fmt.Sprintf("%v", username)
	if len(strUsername) > 9 {
		return fmt.Sprintf("%s...", strUsername[:6])
	} else if len(strUsername) < 9 {
		return strUsername+strings.Repeat(" ", 9-len(strUsername))
	}
	return strUsername
}

func formatUserList(users []map[string]interface{},ChatID int64 ,offset int, totalUsers int) (string, tgbotapi.InlineKeyboardMarkup) {
	if len(users) == 0 {
		return "No users found.", tgbotapi.InlineKeyboardMarkup{}
	}

	var inlineKeyboardButtons [][]tgbotapi.InlineKeyboardButton

	// Add top layer with headers
	topRow := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("👤 Username", "header_username"),
		tgbotapi.NewInlineKeyboardButtonData("📅 Expire Date", "header_expire"),
		tgbotapi.NewInlineKeyboardButtonData("💾 Data Limit", "header_data_limit"),
	)

	inlineKeyboardButtons = append(inlineKeyboardButtons, topRow)

	// Emoji mapping for user statuses
	statusEmojis := map[string]string{
		"active":   "✅",
		"disabled": "❌",
		"limited":  "⚠️",
		"expired":  "⏳",
		"on_hold":  "⏸️",
	}
	for _, user := range users {
		// Extract the necessary fields with safety checks
		username, _ := user["username"].(string) // Default to empty string if not present
		expire, _ := user["expire"].(float64)     // Get as float64
		dataLimit, _ := user["data_limit"].(float64) // Get as float64
		userStatus, _ := user["status"].(string) // Default to empty string if not present
		// Handle the case where userStatus is empty
		if userStatus == "" {
			userStatus = "unknown"
		}

		// Get the corresponding emoji for the user status
		emoji, exists := statusEmojis[userStatus]
		if !exists {
			emoji = "❓" // Default emoji if status is unknown
		}

		// Calculate remaining days until expiration if it's not nil
		var remainingDays string
		if expire != 0 {
			expiryDate := time.Unix(int64(expire), 0)
			remaining := time.Until(expiryDate).Hours() / 24 // Calculate remaining days
			if remaining > 0 {
				remainingDays = fmt.Sprintf("%d days", int(remaining))
			} else {
				remainingDays = "فاقد انقضا" // Expired
			}
		} else {
			remainingDays = "فاقد انقضا" // Expired if no expire date
		}

		// Convert data limit to GB or show unlimited sign if it's not set
		var dataLimitStr string
		if dataLimit > 0 {
			dataLimitGB := int64(dataLimit) / (1024 * 1024 * 1024)
			dataLimitStr = fmt.Sprintf("%d GB", dataLimitGB)
		} else {
			dataLimitStr = "∞" // Unlimited
		}

		// Format the username with the status emoji
		formattedUsername := fmt.Sprintf("%s %s", emoji, formatUsername(username))

		// Append user details as buttons in the keyboard
		inlineKeyboardButtons = append(inlineKeyboardButtons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(formattedUsername, fmt.Sprintf("user_%s", username)),
			tgbotapi.NewInlineKeyboardButtonData(remainingDays, fmt.Sprintf("user_%s", username)),
			tgbotapi.NewInlineKeyboardButtonData(dataLimitStr, fmt.Sprintf("user_%s", username)),
		))
	}

	// Add navigation buttons in a single row
	navigationRow := tgbotapi.NewInlineKeyboardRow()
	if offset > 0 {
		navigationRow = append(navigationRow, tgbotapi.NewInlineKeyboardButtonData("◀️ Back", "back"))
	}

	// Check if there's a next page
	if offset+1 < (totalUsers/1) { // If there are more users than the current offset
		navigationRow = append(navigationRow, tgbotapi.NewInlineKeyboardButtonData("➡️ Next", "next"))
	}

	// Only add the navigation row if it contains buttons
	if len(navigationRow) > 0 {
		inlineKeyboardButtons = append(inlineKeyboardButtons, navigationRow)
	}

	navigationRow2 := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("❌ Exit", "exit"),
	)
	inlineKeyboardButtons = append(inlineKeyboardButtons, navigationRow2)

	// Create the inline keyboard markup
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(inlineKeyboardButtons...)

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalUsers) / 1.0)) // Assuming a limit of 1 user per page
	currentPage := (offset / 1) + 1 // Current page (1-based index)
	user, err := models.ReadUser(ChatID)
	if err != nil || user.Token == "" {
		inlineKeyboard2 := tgbotapi.NewInlineKeyboardMarkup()
		return "Session expired or not logged in. Please use /start.",inlineKeyboard2
	}
	getallusers, err := utils.GetAllUsers(user.Token,Marzban_Url,userOffset,3,"")
	// Prepare the user list message
	
	userList := fmt.Sprintf("Total Users: %d\nPage %d of %d", getallusers, currentPage, totalPages)
	return userList, inlineKeyboard
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



func handleUser(bot *tgbotapi.BotAPI, message *tgbotapi.Message, username string) {
    user, err := models.ReadUser(message.Chat.ID)
    if err != nil || user.Token == "" {
        msg := tgbotapi.NewMessage(message.Chat.ID, "Session expired or not logged in. Please use /start.")
        bot.Send(msg)
        return
    }

    // Call GetUserInfo with the user's token
    userInfoResponse, err := utils.GetUserInfo(user.Token, Marzban_Url, username)
    if err != nil {
        handleError(bot, message.Chat.ID, err)
        return
    }

    // Format user information for display
    userInfoMessage := formatUserInfo(userInfoResponse)

    // Create buttons for the keyboard with specific user context
    keyboard := tgbotapi.NewInlineKeyboardMarkup(
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("🔙 بازگشت به لیست", fmt.Sprintf("back_to_users_%s", username)),
        ),
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("🔗 تغییر لینک اشتراک", fmt.Sprintf("change_subscription_link_%s", username)),
            tgbotapi.NewInlineKeyboardButtonData("♻️ بازنشانی مصرف", fmt.Sprintf("reset_consumption_%s", username)),
        ),
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("⚙️ تنظیم", fmt.Sprintf("settings_%s", username)),
        ),
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("❌ غیرفعال کردن", fmt.Sprintf("deactivate_%s", username)),
        ),
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("🔗 دریافت لینک", fmt.Sprintf("get_link_%s", username)),
        ),
    )
    
    // Add padding to the user info message if needed
    userInfoMessage = fmt.Sprintf("%s%s", userInfoMessage, string(paddingChar))

    editMsg := tgbotapi.NewEditMessageText(message.Chat.ID, message.MessageID, userInfoMessage)
    editMsg.ReplyMarkup = &keyboard
    bot.Send(editMsg)
}

// Function to format user information into a string for display
func formatUserInfo(userInfo map[string]interface{}) string {
    // Current time
    now := time.Now()

    // Extracting necessary fields from the userInfo map
    username := "0"
    if val, ok := userInfo["username"]; ok && val != nil {
        username = val.(string)
    }

    // Status mapping with emojis
    statusEmoji := map[string]string{
        "active":   "🟢",
        "disabled": "❌",
        "limited":  "⚠️",
        "expired":  "⏳",
        "on_hold":  "⏸️",
    }

    status := "0"
    emoji := "❓" // Default emoji for unknown status
    if val, ok := userInfo["status"]; ok && val != nil {
        status = val.(string)
        if e, exists := statusEmoji[status]; exists {
            emoji = e
        }
    }

    // Calculate remaining days until expiration
    expireTimestamp := int64(0)
    if val, ok := userInfo["expire"]; ok && val != nil {
        expireTimestamp = int64(val.(float64)) // Assuming expire is a timestamp
    }

    expirationTime := time.Unix(expireTimestamp, 0)
    var daysUntilExpiry string
    if expireTimestamp == 0 || expirationTime.Before(now) {
        daysUntilExpiry = "بدون انقضا" // No expiration
    } else {
        remainingDays := int(expirationTime.Sub(now).Hours() / 24)
        daysUntilExpiry = fmt.Sprintf("%d روز", remainingDays)
    }

    // Handle storage limit
    storageLimit := "نامحدود"
    if val, ok := userInfo["data_limit"]; ok && val != nil {
        limit := val.(float64) / (1024 * 1024 * 1024) // Convert bytes to GB
        if limit > 0 {
            storageLimit = fmt.Sprintf("%.1f GB", limit)
        }
    }

    usedStorage := 0.0
    if val, ok := userInfo["used_traffic"]; ok && val != nil {
        usedStorage = val.(float64) / (1024 * 1024) // Convert bytes to MB
    }

    lastOnline := "0"
    if val, ok := userInfo["online_at"]; ok && val != nil {
        lastOnline = val.(string)
    }

    note := "0"
    if val, ok := userInfo["note"]; ok && val != nil {
        note = val.(string)
    }

    // Format the output string
    return fmt.Sprintf(
        "👤 نام کاربری: %s\n"+
        "📊 وضعیت: %s %s\n"+
        "⏳ انقضا در: %s\n"+
        "💾 محدودیت حجم: %s\n"+
        "📈 حجم مصرفی: %.1f MB\n"+
        "🕒 آخرین آنلاینی: %s\n"+
        "📝 یادداشت: %s",
        username,
        status,
        emoji,
        daysUntilExpiry,
        storageLimit,
        usedStorage,
        lastOnline,
        note,
    )
}

func handleRevokeSubscription(bot *tgbotapi.BotAPI, message *tgbotapi.Message, username string) {
	user, err := models.ReadUser(message.Chat.ID)
    if err != nil || user.Token == "" {
        msg := tgbotapi.NewMessage(message.Chat.ID, "Session expired or not logged in. Please use /start.")
        bot.Send(msg)
        return
    }

    link, err := utils.RevokeSubscription(user.Token, Marzban_Url, username)
    if err != nil {
        msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Error: %s", err))
        bot.Send(msg)
        return
    }
	link = link + "/" +Marzban_Url
    // Create the response message with the subscription link
    responseMessage := fmt.Sprintf("اشتراک با موفقیت لغو شد.\nلینک اشتراک جدید: %s", link)

    // Create a back button
    keyboard := tgbotapi.NewInlineKeyboardMarkup(
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("🔙 بازگشت به لیست", fmt.Sprintf("back_to_users_%s", username)),
        ),
    )

    // Send the response message
    msg := tgbotapi.NewMessage(message.Chat.ID, responseMessage)
    msg.ReplyMarkup = &keyboard
    bot.Send(msg)
}
