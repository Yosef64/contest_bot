package main

import (
	"app/app/handlers"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var botHandler *handlers.BotHandler

func init() {
	// Initialize bot with your bot token
	bot, err := tgbotapi.NewBotAPI("7521099565:AAGDQx5aMWUOdidp5Vr8_eFHZH0dQOp3bYU")
	if err != nil {
		log.Fatalf("Failed to initialize bot: %v", err)
	}

	botHandler = handlers.NewBotHandler(bot)

}

func main() {
	router := gin.Default()

	// GET /
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Bot is running")
	})

	// POST /webhook
	router.POST("/webhook", func(c *gin.Context) {
		var update tgbotapi.Update
		if err := c.ShouldBindJSON(&update); err != nil {
			log.Printf("Failed to parse update: %v", err)
			c.String(http.StatusBadRequest, "Invalid request")
			return
		}

		botHandler.HandleUpdate(update)
		c.String(http.StatusOK, "OK")
	})

	// Run server on port 8000 or from ENV
	router.Run(":8000")
}