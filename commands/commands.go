// Copyright 2019 Alessandro Pomponio. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package commands contains handlers for the various commands
// the bot may receive.
package commands

import (
	"log"
	"time"

	"github.com/pkg/errors"

	"github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/AlessandroPomponio/serverless-amazon-refbot/persistence"
	"github.com/AlessandroPomponio/serverless-amazon-refbot/repository"
	"github.com/AlessandroPomponio/serverless-amazon-refbot/utility"
)

//HandleCommand handles and performs commands.
func HandleCommand(msg *tgbotapi.Message, bot *tgbotapi.BotAPI) {

	// Use a default reply to mask all the errors to the user,
	// while logging them for the developer on CloudWatch.
	reply := "Unable to handle this command 😫"
	var err error

	switch msg.Command() {
	case "start", "help":
		reply = "Welcome!\nSend me an Amazon link and I'll send you the referral version, if the region is supported."
	case "list":
		reply, err = retrieveLatestRequest(msg.From.ID)
	case "broadcast":
		reply, err = performBroadcast(msg.From.ID, msg.CommandArguments(), bot)
	}

	if err != nil {
		log.Println("HandleCommand: error:", err, "for command", msg.Command(), "userID:", msg.From.ID)
	}

	message := tgbotapi.NewMessage(msg.Chat.ID, reply)
	message.ParseMode = "HTML"
	_, _ = bot.Send(message)
	return

}

// retrieveLatestRequest returns the list of the requests that took
// place in the last 7 days.
func retrieveLatestRequest(userID int) (reply string, err error) {

	err = authorizeUser(userID)
	if err != nil {
		return
	}

	// 168 hours in a week, add with a minus to go back in time.
	lastWeek := time.Now().Add(-168 * time.Hour).Unix()
	requests, err := persistence.GetRequestsSince(lastWeek, repository.DynamoDBClient)
	if err != nil {
		err = errors.New("Unable to retrieve requests")
		return
	}

	reply = utility.FormatRequests(requests)
	return

}

// performBroadcast sends a message to all the users in the database that
// didn't block the bot. It will only send a message every 50ms in order
// not to get limited by Telegram.
func performBroadcast(userID int, messageToSend string, bot *tgbotapi.BotAPI) (reply string, err error) {

	err = authorizeUser(userID)
	if err != nil {
		return
	}

	users, err := persistence.GetAllUsers(repository.DynamoDBClient)
	if err != nil {
		err = errors.New("performBroadcast: unable to retrieve users")
		return
	}

	for _, user := range users {

		// Send a message every 50 milliseconds due to Telegram's restrictions.
		// Bots can send up to 30 messages per seconds.
		response, _ := bot.Request(tgbotapi.NewMessage(int64(user.TelegramID), messageToSend))
		if response.ErrorCode == 403 {
			_ = persistence.UpdateUserBlockStatus(user.TelegramID, true, repository.DynamoDBClient)
		}

		time.Sleep(50 * time.Millisecond)

	}

	return "Completed!", nil

}
