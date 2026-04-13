package database

import (
	"log"
	"time"

	"socialmanager/backend/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	DB.AutoMigrate(&models.Post{}, &models.Task{}, &models.Log{}, &models.Settings{})
	seedData()
}

func seedData() {
	var count int64
	DB.Model(&models.Post{}).Count(&count)
	if count == 0 {
		// Seed Posts
		DB.Create(&models.Post{Title: "Greeting New followers", Content: "Hello!", Summary: "Welcome message", Status: "Draft", Notes: "-", PostAt: time.Now()})
		DB.Create(&models.Post{Title: "Weekly recap", Content: "It was a good week.", Summary: "Recap", Status: "Published", Notes: "-", PostAt: time.Now()})
		DB.Create(&models.Post{Title: "Promotion post", Content: "Check out our new tool.", Summary: "Promo", Status: "Scheduled", Notes: "-", PostAt: time.Now().Add(24 * time.Hour)})
		DB.Create(&models.Post{Title: "Maintenance Note", Content: "Offline for 1 hour.", Summary: "Maintenance", Status: "Draft", Notes: "-", PostAt: time.Now()})
		DB.Create(&models.Post{Title: "Question of the day", Content: "What do you think?", Summary: "Engagement", Status: "Draft", Notes: "-", PostAt: time.Now()})

		// Seed tasks
		now := time.Now()
		DB.Create(&models.Task{Name: "Daily sync", Description: "Sync local data with cloud", Type: "Sync", Schedule: "Daily", Enabled: true, LastRun: &now})
		DB.Create(&models.Task{Name: "Clean logs", Description: "Delete old logs", Type: "Maintenance", Schedule: "Weekly", Enabled: true, LastRun: &now})
		DB.Create(&models.Task{Name: "Check upcoming posts", Description: "Reminder for scheduled posts", Type: "Reminder", Schedule: "Hourly", Enabled: true, LastRun: &now})
		DB.Create(&models.Task{Name: "Back up DB", Description: "SQLite backup", Type: "Maintenance", Schedule: "Daily", Enabled: false})
		DB.Create(&models.Task{Name: "Ping official API", Description: "Keep token alive", Type: "Sync", Schedule: "Hourly", Enabled: true, LastRun: &now})

		// Seed Settings
		DB.Create(&models.Settings{
			Theme:            "dark",
			AppName:          "Social Desktop Manager",
			DatabasePath:     "./data.db",
			AutoStart:        false,
			LogLevel:         "info",
			Timezone:         "Local",
			ConnectionStatus: "connected",
			ConnectedAccount: "official_page_demo",
		})

		// Seed logs
		for i := 0; i < 10; i++ {
			DB.Create(&models.Log{
				Time:    time.Now().Add(-time.Duration(i) * time.Hour),
				Module:  "System",
				Action:  "Init",
				Status:  "Success",
				Details: "System operational.",
			})
		}
	}
}
