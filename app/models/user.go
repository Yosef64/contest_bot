package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID        int64     `json:"id"`
	TelegramID int64     `json:"telegram_id"`
	Username  string    `json:"name"`
	FirstName string    `json:"first_name"`
	Grade	string    `json:"grade"`
	PhoneNumber string    `json:"phoneNumber"`
	IsRegistered bool   `json:"is_registered"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}



func GetUserByTelegramID(telegramID int64) (*User, error) {
	user := &User{}
	
	url := fmt.Sprintf("https://victory-contest-backend.vercel.app/api/student/%d", telegramID)
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user from backend: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK response code: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(user)
	if err != nil {
		return nil, fmt.Errorf("failed to decode user: %v", err)
	}

	return user, nil
}


