// Copyright 2019 Alessandro Pomponio. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utility

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// GetMessageEntities returns an array of entities.
// It can be message.Entities, message.CaptionEntities or an empty slice.
func GetMessageEntities(msg *tgbotapi.Message) (entities []tgbotapi.MessageEntity) {

	if msg.Entities != nil {
		entities = msg.Entities
	} else if msg.CaptionEntities != nil {
		entities = msg.CaptionEntities
	}

	return

}

// GetMessageText returns a text for the message.
// It can be message.Text, message.Caption or an empty string.
func GetMessageText(msg *tgbotapi.Message) (text string) {

	if msg.Text != "" {
		text = msg.Text
	} else if msg.Caption != "" {
		text = msg.Caption
	}

	return

}
