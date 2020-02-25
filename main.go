// Copyright 2019 Alessandro Pomponio. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/AlessandroPomponio/serverless-amazon-refbot/commands"
	"github.com/AlessandroPomponio/serverless-amazon-refbot/messages"
	"github.com/AlessandroPomponio/serverless-amazon-refbot/persistence"
	"github.com/AlessandroPomponio/serverless-amazon-refbot/repository"
)

// main is the "entrance" to the program and the function that
// must be registered in the lambda function's page in the
// "handler" field.
func main() {

	// If any of the environment variables are missing,
	// the application will exit.
	repository.LoadEnvVariables()

	// The AWS Session creation may fail.
	// In this case, we will just try to handle the message.
	err := repository.CreateAWSSession()
	if err != nil {
		log.Println("main: unable to create AWS session:", err)
	} else {
		repository.StartDynamoDBClient()
	}

	lambda.Start(HandleUpdate)

}

//HandleUpdate handles a Telegram Update.
func HandleUpdate(update tgbotapi.Update) {

	// In some cases, like if the Lambda Proxy is active,
	// the message may not be correctly deserialized, ending
	// with a nil pointer dereference down the execution.
	msg := update.Message
	if msg == nil {
		log.Fatal("HandleUpdate: the message was nil. Make sure the Lambda Proxy integration is not active")
	}

	bot, err := tgbotapi.NewBotAPI(repository.TelegramBotToken)
	if err != nil {
		log.Fatalf("HandleUpdate: unable to create the bot instance: %s", err)
	}

	_, _ = bot.Send(tgbotapi.NewChatAction(msg.Chat.ID, "typing"))

	// We don't want to record users that just use the commands
	// in our DynamoDB instance, so we will return after handling
	// the command.
	if msg.IsCommand() {
		commands.HandleCommand(msg, bot)
		return
	}

	urls, err := messages.HandleMessage(msg, bot)
	if err != nil || repository.DynamoDBClient == nil {
		log.Println("HandleUpdate: error while handling message or nil DynamoDB client:", err, repository.DynamoDBClient)
		return
	}

	//Persist results
	err = persistence.PutUser(msg.From.ID, repository.DynamoDBClient)
	if err != nil {
		log.Println(err)
	}

	for _, uri := range urls {
		err = persistence.PutRequest(msg.From.ID, uri, repository.DynamoDBClient)
		if err != nil {
			log.Println(err)
		}
	}

}
