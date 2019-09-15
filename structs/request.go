package structs

import (
	"os"
	"time"
)

const (
	requestTableKey = "REQUEST_TABLE_NAME"
)

// Request represents a request.
// It contains an unique identifier, the Telegram
// ID of the user, the URL that was returned to
// the user, the timestamp of the request and an
// Unix representation of the time to be used for
// the table TTL, if needed.
type Request struct {
	XID        string
	TelegramID int
	URL        string
	Time       time.Time
	UnixTime   int64
}

// Table returns the name of the Request table
// reading it from the environment variables.
func (Request) Table() string {
	return os.Getenv(requestTableKey)
}
