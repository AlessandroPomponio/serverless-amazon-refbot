package utility

import (
	"fmt"
	"strings"
	"time"

	"github.com/AlessandroPomponio/serverless-amazon-refbot/structs"
)

// FormatRequests formats a slice of requests.
func FormatRequests(requests []structs.Request) string {

	builder := strings.Builder{}

	for _, request := range requests {

		builder.WriteString(
			fmt.Sprintf("➡️ <a href=\"tg://user?id=%d\">%d</a> requested %s on %s\n\n",
				request.TelegramID, request.TelegramID, request.URL, FormatDate(request.Time)))

	}

	return builder.String()

}

// FormatDate formats a Time using the format
// Mon 2 Jan 2006 15:04:05.
func FormatDate(date time.Time) string {
	return date.Format("Mon 2 Jan 2006 15:04:05")
}
