package handlers

import (
	"app/app/models"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotHandler struct {
	bot         *tgbotapi.BotAPI
	providerToken string
}

func NewBotHandler(bot *tgbotapi.BotAPI) *BotHandler {
	providerToken := os.Getenv("TELEGRAM_PAYMENT_PROVIDER_TOKEN")
	if providerToken == "" {
		providerToken = "6141645565:TEST:9ytidmVoJ0Kq9nFttiQO" // Default for testing
	}

	return &BotHandler{
		bot:         bot,
		providerToken: providerToken,
	}
}

func (h *BotHandler) HandleUpdate(update tgbotapi.Update) {
	// Handle pre-checkout queries
	log.Printf("Received update: %+v", update)
	if update.PreCheckoutQuery != nil {
		h.handlePreCheckoutQuery(update.PreCheckoutQuery)
		return
	}

	// Handle successful payments
	if update.Message != nil && update.Message.SuccessfulPayment != nil {
		h.handleSuccessfulPayment(update.Message)
		return
	}

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

	// Handle payment commands
	if strings.HasPrefix(text, "/payment") {
		h.handlePaymentCommand(chatID, userID, text)
		return
	}

	// Handle payment amount input

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



func (h *BotHandler) handlePaymentCommand(chatID, userID int64, text string) {
	// Check if user is registered
	// user, err := models.GetUserByTelegramID(userID)
	// if err != nil {
	// 	log.Printf("Error getting user: %v", err)
	// 	h.sendMessage(chatID, "Sorry, there was an error. Please try again.")
	// 	return
	// }

	// if user == nil {
	// 	h.sendMessage(chatID, "Please register first by clicking the Register button or sending /register")
	// 	return
	// }

	// Parse payment command
	// parts := strings.Fields(text)
	// if len(parts) == 1 {
	// 	// Show payment menu
	// 	h.showPaymentMenu(chatID)
	// 	return
	// }

	// Handle specific payment subcommands
	// switch parts[1] {
	// case "start":
	// err := h.sendInvoice(chatID, userID, 1000, "Telegram Bot Payment")
	// if err != nil {
	// 	log.Printf("Error sending invoice: %v", err)
	// 	h.sendMessage(chatID, "❌ Failed to create payment. Please try again later.")
	// 	return
	// }
	h.showPaymentMenu(chatID)

// 	case "help":
// 		h.sendMessage(chatID, `💳 <b>Payment Help</b>

// Available payment commands:
// • /payment start - Start a new payment
// • /payment help - Show this help message

// Supported currencies: USD, EUR, NGN, GHS, KES, UGX, TZS, ZAR

// The payment will be processed securely through Telegram Payments.`)
// 	default:
// 		h.sendMessage(chatID, "Unknown payment command. Use /payment help for available options.")
	// }
}

// func (h *BotHandler) handlePaymentAmountInput(chatID, userID int64, text string) {
// 	// Parse amount
// 	amount, err := strconv.ParseFloat(text, 64)
// 	if err != nil {
// 		h.sendMessage(chatID, "❌ Invalid amount. Please enter a valid number (e.g., 1000):")
// 		return
// 	}

// 	if amount <= 0 {
// 		h.sendMessage(chatID, "❌ Amount must be greater than 0. Please try again:")
// 		return
// 	}

// 	// Create and send invoice
// 	err = h.sendInvoice(chatID, userID, amount, "Telegram Bot Payment")
// 	if err != nil {
// 		log.Printf("Error sending invoice: %v", err)
// 		h.sendMessage(chatID, "❌ Failed to create payment. Please try again later.")
// 		return
// 	}
// }

func (h *BotHandler) sendInvoice(chatID, userID int64, amount float64, description string) error {
	// Convert amount to cents (Telegram expects amounts in cents)
	amountCents := int(amount * 100)

	// Create invoice item
	item := tgbotapi.LabeledPrice{
		Label:  description,
		Amount: amountCents,
	}

	// Create invoice config
	invoice := tgbotapi.NewInvoice(chatID, description, description, 
		fmt.Sprintf("payment_%d_%d", userID, time.Now().Unix()), 
		h.providerToken, "USD", "ETB", []tgbotapi.LabeledPrice{item})

	// Set optional parameters
	invoice.StartParameter = "payment"
	invoice.NeedName = true
	invoice.NeedEmail = true
	invoice.NeedPhoneNumber = true
	invoice.NeedShippingAddress = false
	invoice.IsFlexible = false
	invoice.SuggestedTipAmounts = []int{100, 200, 500} // Optional tip amounts in cents
	invoice.Payload = fmt.Sprintf("user_%d", userID) // Custom payload to identify the user
	invoice.MaxTipAmount = 600

	// Send the invoice
	_, err := h.bot.Send(invoice)
	return err
}

func (h *BotHandler) handlePreCheckoutQuery(query *tgbotapi.PreCheckoutQuery) {
	// Log the pre-checkout query
	log.Printf("Pre-checkout query received: %+v", query)

	// For now, we'll approve all pre-checkout queries
	// In a real application, you might want to validate the payment details
	answer := tgbotapi.CallbackConfig{
		CallbackQueryID: query.ID,
		Text:            "",
		ShowAlert:       false,
	}
	
	_, err := h.bot.Send(answer)
	if err != nil {
		log.Printf("Error answering pre-checkout query: %v", err)
	}
}

func (h *BotHandler) handleSuccessfulPayment(message *tgbotapi.Message) {
	payment := message.SuccessfulPayment
	chatID := message.Chat.ID
	userID := message.From.ID

	// Log the successful payment
	log.Printf("Successful payment received: %+v", payment)

	// Create payment record
	paymentRecord := &models.Payment{
		UserID:      userID,
		TelegramID:  userID,
		Amount:      float64(payment.TotalAmount) / 100.0, // Convert from cents
		Currency:    payment.Currency,
		Description: payment.InvoicePayload,
		Status:      "completed",
		TelegramPaymentChargeID: payment.TelegramPaymentChargeID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Here you would typically save the payment record to your database
	log.Printf("Payment record created: %+v", paymentRecord)

	// Send success message to user
	successMsg := fmt.Sprintf(`✅ <b>Payment Successful!</b>

Amount: $%.2f %s
Description: %s
Transaction ID: %s

Thank you for your payment! Your order has been processed successfully.`, 
		paymentRecord.Amount, paymentRecord.Currency, 
		paymentRecord.Description, paymentRecord.TelegramPaymentChargeID)

	h.sendMessage(chatID, successMsg)
}

func (h *BotHandler) showPaymentMenu(chatID int64) {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💳 Start Payment", "payment_start"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❓ Payment Help", "payment_help"),
		),
	)

	msg := `💳 <b>Payment Menu</b>

Welcome to the payment system! You can make secure payments using Telegram's built-in payment system.

Click "Start Payment" to begin a new payment transaction.`

	h.sendMessageWithKeyboard(chatID, msg, keyboard)
}

func (h *BotHandler) handleCallbackQuery(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	data := callback.Data

	// Answer the callback query to remove loading state
	callbackAnswer := tgbotapi.NewCallback(callback.ID, "")
	h.bot.Send(callbackAnswer)

	switch data {
	case "payment_start":
		h.sendMessage(chatID, "Please enter the amount you want to pay (in USD):")
	case "payment_help":
		h.sendMessage(chatID, `💳 <b>Payment Help</b>

Available payment commands:
• /payment start - Start a new payment
• /payment help - Show this help message

Supported currencies: USD, EUR, NGN, GHS, KES, UGX, TZS, ZAR

The payment will be processed securely through Telegram Payments.`)

	case "payment_cancel":
		h.sendMessage(chatID, "❌ Payment cancelled.")

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
				tgbotapi.NewInlineKeyboardButtonURL("Open Mini App", "https://7wwb0knl-5173.euw.devtunnels.ms/"),
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
