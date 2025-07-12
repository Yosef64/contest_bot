package handlers

import (
	"app/app/models"
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotHandler struct {
	bot *tgbotapi.BotAPI
}

func NewBotHandler(bot *tgbotapi.BotAPI) *BotHandler {
	return &BotHandler{bot: bot}
}

func (h *BotHandler) HandleUpdate(update tgbotapi.Update) {
	// Handle callback queries (inline keyboard buttons)
	if update.CallbackQuery != nil {
		h.handleCallbackQuery(update.CallbackQuery)
		return
	}

	if update.Message == nil {
		return
	}

	message := update.Message
	chatID := message.Chat.ID
	userID := message.From.ID
	text := message.Text

	// Handle /start command
	if text == "/start" {
		h.handleStartCommand(chatID, userID)
		return
	}

	// Handle registration command
	if text == "/register" {
		h.handleRegisterCommand(chatID, userID)
		return
	}

	// Check if user is registered for other commands
	user, err := models.GetUserByTelegramID(userID)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		h.sendMessage(chatID, "Sorry, there was an error. Please try again.")
		return
	}

	if user == nil {
		h.sendMessage(chatID, "Please register first by clicking the Register button or sending /register")
		return
	}

	// Handle other commands for registered users
	h.handleRegisteredUserMessage(chatID, text)
}

func (h *BotHandler) handleCallbackQuery(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	data := callback.Data

	// Answer the callback query to remove loading state
	callbackAnswer := tgbotapi.NewCallback(callback.ID, "")
	h.bot.Send(callbackAnswer)

	switch data {

	default:
		h.sendMessage(chatID, "Unknown button action")
	}
}

func (h *BotHandler) handleStartCommand(chatID, userID int64) {
	// Check if user exists in backend
	dbUser, err := models.GetUserByTelegramID(userID)
	if err != nil {
		log.Printf("Error checking user: %v", err)
		h.sendMessage(chatID, "Sorry, there was an error. Please try again.")
		return
	}
	fmt.Print("dbUser: ", dbUser)

	if dbUser == nil || *dbUser == (models.User{}) {
		// User doesn't exist in backend
		welcomeMsg := `Welcome! 👋 

I'm your Telegram bot. To get started, you need to register first.

Please click the Register button below to complete your registration.`

		// Create inline keyboard with Register button that opens registration page
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonWebApp("Open Mini App", tgbotapi.WebAppInfo{URL: "https://7wwb0knl-5173.euw.devtunnels.ms/"}),
			),
		)

		h.sendMessageWithKeyboard(chatID, welcomeMsg, keyboard)
		return
	}

	// User exists in backend
	welcomeMsg := `Welcome back! 🎉

You're all set up and ready to use the bot.

Available commands:
- /start - Show this message
- /register - Re-register (if needed)

Feel free to send me any message!`

	h.sendMessage(chatID, welcomeMsg)
}

func (h *BotHandler) handleRegisterCommand(chatID, userID int64) {
	// Check if user exists in backend
	dbUser, err := models.GetUserByTelegramID(userID)
	if err != nil {
		log.Printf("Error checking user: %v", err)
		h.sendMessage(chatID, "Sorry, there was an error. Please try again.")
		return
	}

	if dbUser == nil {
		// User doesn't exist in backend
		errorMsg := `❌ Registration failed!

It seems you are not registered in our system. Please contact support to get registered first.`

		h.sendMessage(chatID, errorMsg)
		return
	}

	// User exists in backend
	successMsg := `✅ Registration successful!

Welcome to the bot! You're now registered and can use all features.

Send /start to see available commands.`

	h.sendMessage(chatID, successMsg)
}

func (h *BotHandler) handleRegisteredUserMessage(chatID int64, text string) {
	// Handle messages from registered users
	response := "You said: " + text

	// Add some basic command handling
	if strings.HasPrefix(text, "/") {
		switch text {
		case "/help":
			response = `Available commands:
- /start - Show welcome message
- /register - Re-register (if needed)
- /help - Show this help message

You can also send me any message and I'll respond!`
		default:
			response = "Unknown command. Send /help for available commands."
		}
	}

	h.sendMessage(chatID, response)
}

func (h *BotHandler) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"

	_, err := h.bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func (h *BotHandler) sendMessageWithKeyboard(chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = keyboard

	_, err := h.bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message with keyboard: %v", err)
	}
}
