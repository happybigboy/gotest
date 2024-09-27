package main

import (
	"log"
	"main/handlers"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"main/models" // Replace with the correct path to your models package
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// InitializeDatabase sets up the SQLite database and runs migrations
func InitializeDatabase() (*gorm.DB, error) {
	// Open the SQLite database
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}

	// Migrate the models
	err = models.MigrateModels(db)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate the database: %v", err)
	}

	fmt.Println("Database migration completed successfully")
	return db, nil
}

func main() {
	// Initialize the database
	db, err := InitializeDatabase()
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}

	// Here, you can start using the database with GORM functions
	// For example, querying or inserting data
	fmt.Println("Database initialized:", db)
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
			handlers.HandleStart(bot, update.Message,db)
		default:
			currentState := handlers.GetCurrentState(db, update.Message.Chat.ID)
			switch currentState {
			case "waiting_for_username":
				handlers.HandleUsername(bot, update.Message, db)
			//waiting_for_password
			case "waiting_for_password":
				handlers.HandlePassword(bot, update.Message, db)
			
			case "waiting_for_menu":
				handlers.HandleMenu(bot, update.Message, db)
			default:
				handlers.HandleMenuFallback(bot, update.Message, db) // Default to menu if no valid state is set
			}
		}
	}
}
