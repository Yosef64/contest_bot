package main

import (
	"app/app/handlers"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
    botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
    webhookURL := os.Getenv("WEBHOOK_URL") // e.g., https://yourdomain.com/<token>

    if botToken == "" || webhookURL == "" {
        log.Fatal("TELEGRAM_BOT_TOKEN or WEBHOOK_URL not set")
    }

    bot, err := tgbotapi.NewBotAPI(botToken)
    if err != nil {
        log.Fatal(err)
    }

    // Create bot handler
    botHandler := handlers.NewBotHandler(bot)

    router := gin.Default()

    router.POST("/webhook", func(c *gin.Context) {
        var update tgbotapi.Update
        if err := c.ShouldBindJSON(&update); err != nil {
            log.Println("Invalid update received:", err)
            c.Status(http.StatusBadRequest)
            return
        }

        // Handle the update with our bot handler
        botHandler.HandleUpdate(update)

        c.Status(http.StatusOK)
    })
    
    router.GET("/", func(c *gin.Context) {
        c.JSON(http.StatusOK, "Bot is running")
    })

    log.Println("Bot is listening on :8000")
    router.Run(":8000")
}
