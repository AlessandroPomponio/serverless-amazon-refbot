// Copyright 2019 Alessandro Pomponio. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utility

import (
	"fmt"
	"testing"
	"time"

	"github.com/AlessandroPomponio/serverless-amazon-refbot/structs"
)

func TestFormatDate(t *testing.T) {

	tests := []struct {
		name string
		date time.Time
		want string
	}{
		{
			name: "Epoch",
			date: time.Unix(0, 0),
			want: "Thu 1 Jan 1970 01:00:00",
		},
		{
			name: "Jan 1, 2019",
			date: time.Unix(1546300800, 0),
			want: "Tue 1 Jan 2019 01:00:00",
		},
		{
			name: "Dec 31, 2019",
			date: time.Unix(1577750400, 0),
			want: "Tue 31 Dec 2019 01:00:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatDate(tt.date); got != tt.want {
				t.Errorf("FormatDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatRequests(t *testing.T) {

	referenceTime := time.Unix(0, 0)
	formattedTime := FormatDate(referenceTime)

	tests := []struct {
		name     string
		requests []structs.Request
		want     string
	}{
		{
			name: "One request",
			requests: []structs.Request{
				{
					TelegramID: 777000,
					URL:        "https://amzn.to/",
					Time:       referenceTime,
				},
			},
			want: fmt.Sprintf("➡️ <a href=\"tg://user?id=777000\">777000</a> requested https://amzn.to/ on %s\n\n", formattedTime),
		},
		{
			name: "Three requests",
			requests: []structs.Request{
				{
					TelegramID: 777000,
					URL:        "https://amzn.to/",
					Time:       referenceTime,
				},
				{
					TelegramID: 777000,
					URL:        "https://amzn.to/",
					Time:       referenceTime,
				},
				{
					TelegramID: 777000,
					URL:        "https://amzn.to/",
					Time:       referenceTime,
				},
			},
			want: fmt.Sprintf("➡️ <a href=\"tg://user?id=777000\">777000</a> requested https://amzn.to/ on %s\n\n➡️ <a href=\"tg://user?id=777000\">777000</a> requested https://amzn.to/ on %s\n\n➡️ <a href=\"tg://user?id=777000\">777000</a> requested https://amzn.to/ on %s\n\n", referenceTime, referenceTime, referenceTime),
		},
		{
			name:     "No requests",
			requests: []structs.Request{},
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatRequests(tt.requests); got != tt.want {
				t.Errorf("FormatRequests() = %v, want %v", got, tt.want)
			}
		})
	}
}
