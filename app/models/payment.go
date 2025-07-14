package models

import (
	"time"
)

// Payment represents a payment transaction
type Payment struct {
	ID            string    `json:"id"`
	UserID        int64     `json:"user_id"`
	TelegramID    int64     `json:"telegram_id"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	Description   string    `json:"description"`
	Status        string    `json:"status"`
	TelegramPaymentChargeID string `json:"telegram_payment_charge_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// InvoiceItem represents an item in the invoice
type InvoiceItem struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Amount      int64   `json:"amount"` // Amount in cents
	Quantity    int     `json:"quantity"`
}

// InvoiceRequest represents a request to create an invoice
type InvoiceRequest struct {
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Payload     string        `json:"payload"`
	ProviderToken string      `json:"provider_token"`
	Currency    string        `json:"currency"`
	Prices      []InvoiceItem `json:"prices"`
	StartParameter string     `json:"start_parameter"`
}

// PreCheckoutQuery represents a pre-checkout query from Telegram
type PreCheckoutQuery struct {
	ID          string `json:"id"`
	From        struct {
		ID        int64  `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Username  string `json:"username"`
	} `json:"from"`
	Currency    string `json:"currency"`
	TotalAmount int64  `json:"total_amount"`
	InvoicePayload string `json:"invoice_payload"`
}

// SuccessfulPayment represents a successful payment from Telegram
type SuccessfulPayment struct {
	Currency                string `json:"currency"`
	TotalAmount            int64  `json:"total_amount"`
	InvoicePayload         string `json:"invoice_payload"`
	TelegramPaymentChargeID string `json:"telegram_payment_charge_id"`
	ProviderPaymentChargeID string `json:"provider_payment_charge_id"`
}

// PaymentRequest represents a payment request from the bot
type PaymentRequest struct {
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Description string  `json:"description"`
	UserID      int64   `json:"user_id"`
	Title       string  `json:"title"`
} 