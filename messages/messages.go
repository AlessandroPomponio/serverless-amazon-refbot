// Copyright 2019 Alessandro Pomponio. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package messages contains the functions to handle
// messages and their content.
package messages

import (
	"fmt"
	"log"
	"strings"
	"unicode/utf16"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"github.com/zpnk/go-bitly"

	"github.com/AlessandroPomponio/serverless-amazon-refbot/repository"
	"github.com/AlessandroPomponio/serverless-amazon-refbot/urlwork"
	"github.com/AlessandroPomponio/serverless-amazon-refbot/utility"
)

// HandleMessage handles messages and returns referral URLs.
func HandleMessage(msg *tgbotapi.Message, bot *tgbotapi.BotAPI) (returnedURLs []string, err error) {

	text := utility.GetMessageText(msg)
	entities := utility.GetMessageEntities(msg)
	tUTF16 := utf16.Encode([]rune(text))
	urls := GetURLs(tUTF16, entities)

	if len(urls) == 0 {
		_, _ = bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "No URLs found 😢"))
		err = errors.Errorf("No URLs found.\nText was: %s", text)
		return
	}

	bitlyClient := bitly.New(repository.BitlyAPIKey)
	for _, url := range urls {

		refurl, err := urlwork.GetRefURL(url, repository.ReferralID, bitlyClient)
		if err != nil {
			log.Println(err)
			continue
		}

		returnedURLs = append(returnedURLs, refurl)

	}

	if len(returnedURLs) == 0 {
		_, _ = bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "No matching URLs found 😢"))
		err = errors.Errorf("No matching URLs found.\nText was: %s", text)
		return
	}

	builder := strings.Builder{}
	for _, uri := range returnedURLs {
		builder.WriteString(fmt.Sprintf("➡️ %s\n\n", uri))
	}
	_, _ = bot.Send(tgbotapi.NewMessage(msg.Chat.ID, builder.String()))

	return

}

// GetURLs returns the urls and the markdown links in a message.
func GetURLs(tUTF16 []uint16, entities []tgbotapi.MessageEntity) (urls []string) {

	if len(entities) == 0 {
		return urls
	}

	for _, entity := range entities {
		switch entity.Type {
		case "url":
			urls = append(urls, string(utf16.Decode(tUTF16[entity.Offset:entity.Offset+entity.Length])))
		case "text_link":
			urls = append(urls, entity.URL)
		}
	}

	return urls

}
