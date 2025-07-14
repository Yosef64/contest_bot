package config

import (
	"os"
)

// PaymentConfig holds configuration for payment settings
type PaymentConfig struct {
	TelegramPaymentProviderToken string
	DefaultCurrency string
	SupportedCurrencies []string
}

// GetPaymentConfig returns the payment configuration
func GetPaymentConfig() *PaymentConfig {
	return &PaymentConfig{
		TelegramPaymentProviderToken: getEnv("TELEGRAM_PAYMENT_PROVIDER_TOKEN", "TEST_PAYMENT_PROVIDER_TOKEN"),
		DefaultCurrency: "USD",
		SupportedCurrencies: []string{
			"USD", "EUR", "NGN", "GHS", "KES", "UGX", "TZS", "ZAR",
		},
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// IsSupportedCurrency checks if a currency is supported
func (pc *PaymentConfig) IsSupportedCurrency(currency string) bool {
	for _, supported := range pc.SupportedCurrencies {
		if supported == currency {
			return true
		}
	}
	return false
} 