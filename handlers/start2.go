package handlers


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
		models.ModifyUser(db,message.Chat.ID,user.Username,message.Text)
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
			bot.Send(msg)
		}
	} else {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Invalid password. Please try again.")
		bot.Send(msg)
	}
}