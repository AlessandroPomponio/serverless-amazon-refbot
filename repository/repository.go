// Copyright 2019 Alessandro Pomponio. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package repository contains variables used
// throughout the application.
package repository

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	tgBotTokenName = "TG_KEY"
	bAPIKeyName    = "BITLY_KEY"
	refIDKeyName   = "REF_ID"
	amazonDomain   = "AMAZON_DOMAIN"
)

var (
	//environment variables

	//TelegramBotToken is the Telegram bot token.
	TelegramBotToken string
	//BitlyAPIKey is the Bitly API key.
	BitlyAPIKey string
	//ReferralID is the Amazon referral ID.
	ReferralID string
	//AmazonDomain is the Amazon domain for which
	//the ReferralID is valid.
	AmazonDomain string

	//AWS-related variables

	//AWSSession is the current AWS session.
	AWSSession *session.Session
	//DynamoDBClient is the client to access DynamoDB.
	DynamoDBClient *dynamodb.DynamoDB
)

// LoadEnvVariables loads the environment variables and checks their values.
// If they're empty, the application will quit.
func LoadEnvVariables() {

	TelegramBotToken = os.Getenv(tgBotTokenName)
	if TelegramBotToken == "" {
		log.Fatalf("Missing Telegram bot token. Make sure you have it in your environment variables with the key %s", tgBotTokenName)
	}

	BitlyAPIKey = os.Getenv(bAPIKeyName)
	if BitlyAPIKey == "" {
		log.Fatalf("Missing Bitly API Key. Make sure you have it in your environment variables with the key %s", bAPIKeyName)
	}

	ReferralID = os.Getenv(refIDKeyName)
	if ReferralID == "" {
		log.Fatalf("Missing referral ID. Make sure you have it in your environment variables with the key %s", refIDKeyName)
	}

	AmazonDomain = os.Getenv(amazonDomain)
	if ReferralID == "" {
		log.Fatalf("Missing Amazon domain. Make sure you have it in your environment variables with the key %s", amazonDomain)
	}

}

// CreateAWSSession creates an AWS session.
func CreateAWSSession() (err error) {
	AWSSession, err = session.NewSession()
	return
}

// StartDynamoDBClient starts the DynamoDB client.
func StartDynamoDBClient() {
	DynamoDBClient = dynamodb.New(AWSSession)
}
