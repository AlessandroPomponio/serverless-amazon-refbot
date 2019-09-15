package structs

import (
	"os"
)

const (
	userTableKey = "USER_TABLE_NAME"
)

// User represents a telegram user.
// It contains the Telegram ID of the user,
// a boolean field indicating whether it's
// an admin and a boolean flag that tells
// if the user blocked the bot.
type User struct {
	TelegramID    int
	IsAdmin       bool
	HasBlockedBot bool
}

// Table returns the name of the User table
// reading it from the environment variables.
func (User) Table() string {
	return os.Getenv(userTableKey)
}
