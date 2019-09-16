// Copyright 2019 Alessandro Pomponio. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utility

import (
	"encoding/json"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func unmarshalTestMessage(rawJSON string, t *testing.T) *tgbotapi.Message {

	var msg tgbotapi.Message
	err := json.Unmarshal([]byte(rawJSON), &msg)
	if err != nil {
		t.Errorf("unmarshalTestMessage: unable to unmarshal message: %s\n rawJSON: %s", err, rawJSON)
	}

	return &msg

}

func TestGetMessageEntities(t *testing.T) {

	threeMessageEntities := `{"text":"/command @telegram hello","entities":[{"offset":0,"length":8,"type":"bot_command"},{"offset":9,"length":9,"type":"mention"},{"offset":19,"length":5,"type":"italic"}]}`
	oneMessageEntity := `{"text":"Hello! This is a test 😀","entities":[{"offset":17,"length":4,"type":"bold"}]}`
	oneCaptionEntity := `{"caption":"Check this out @telegram!","caption_entities":[{"offset":15,"length":9,"type":"mention"}]}`
	noEntities := `{"text": "Whoops! I forgot all message entities \ud83d\ude44"}`

	tests := []struct {
		name         string
		msg          *tgbotapi.Message
		wantEntities []tgbotapi.MessageEntity
	}{
		{
			name: "Three message entities",
			msg:  unmarshalTestMessage(threeMessageEntities, t),
			wantEntities: []tgbotapi.MessageEntity{
				{
					Type:   "bot_command",
					Offset: 0,
					Length: 8,
				},
				{
					Type:   "mention",
					Offset: 9,
					Length: 9,
				},
				{
					Type:   "italic",
					Offset: 19,
					Length: 5,
				},
			},
		},
		{
			name: "One message entity",
			msg:  unmarshalTestMessage(oneMessageEntity, t),

			wantEntities: []tgbotapi.MessageEntity{
				{
					Type:   "bold",
					Offset: 17,
					Length: 4,
				},
			},
		},
		{
			name: "Caption entity",
			msg:  unmarshalTestMessage(oneCaptionEntity, t),
			wantEntities: []tgbotapi.MessageEntity{
				{
					Type:   "mention",
					Offset: 15,
					Length: 9,
				},
			},
		},
		{
			name:         "No entities",
			msg:          unmarshalTestMessage(noEntities, t),
			wantEntities: []tgbotapi.MessageEntity{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotEntities := GetMessageEntities(tt.msg)
			if len(gotEntities) != len(tt.wantEntities) {
				t.Errorf("Different number of entities in %s.\nGetMessageEntities() = %v, want %v", tt.name, gotEntities, tt.wantEntities)
			}

			for i := 0; i < len(gotEntities); i++ {

				if gotEntities[i].Type != tt.wantEntities[i].Type {
					t.Errorf("Entities of different types in %s.\nGetMessageEntities() = %v, want %v", tt.name, gotEntities, tt.wantEntities)
				}

				if gotEntities[i].Offset != tt.wantEntities[i].Offset {
					t.Errorf("Entities at different offsets in %s.\nGetMessageEntities() = %v, want %v", tt.name, gotEntities, tt.wantEntities)
				}

				if gotEntities[i].Length != tt.wantEntities[i].Length {
					t.Errorf("Entities with different lengths in %s.\nGetMessageEntities() = %v, want %v", tt.name, gotEntities, tt.wantEntities)
				}

			}

		})
	}

}

func TestGetMessageText(t *testing.T) {

	messageText := `{"text":"Hello! This is a test 😀","caption":""}`
	captionText := `{"text":"","caption":"Hello! This is a test 😀"}`
	noText := `{"text":"","caption":""}`

	tests := []struct {
		name     string
		msg      *tgbotapi.Message
		wantText string
	}{
		{
			name:     "Message text",
			msg:      unmarshalTestMessage(messageText, t),
			wantText: "Hello! This is a test 😀",
		},
		{
			name:     "Caption",
			msg:      unmarshalTestMessage(captionText, t),
			wantText: "Hello! This is a test 😀",
		},
		{
			name:     "No text",
			msg:      unmarshalTestMessage(noText, t),
			wantText: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotText := GetMessageText(tt.msg); gotText != tt.wantText {
				t.Errorf("GetMessageText() = %v, want %v", gotText, tt.wantText)
			}
		})
	}
}
