// Copyright 2019 Alessandro Pomponio. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package persistence

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/pkg/errors"

	"github.com/AlessandroPomponio/serverless-amazon-refbot/structs"
)

// GetAllUsers returns all the users that didn't block the bot.
func GetAllUsers(client *dynamodb.DynamoDB) (users []structs.User, err error) {

	filter := expression.Name("HasBlockedBot").Equal(expression.Value(false))
	projection := expression.NamesList(expression.Name("TelegramID"))

	expr, err := expression.NewBuilder().WithFilter(filter).WithProjection(projection).Build()
	if err != nil {
		err = errors.Errorf("GetAllUsers: error while building the expression: %s", err)
		return
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(structs.User{}.Table()),
	}

	// Make the DynamoDB Query API call
	result, err := client.Scan(params)
	if err != nil {
		err = errors.Errorf("GetAllUsers: error while querying the database: %s", err)
		return
	}

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		err = errors.Errorf("GetAllUsers: error while unmarshaling the users: %s", err)
		return
	}

	return

}

// PutUser saves a user on DynamoDB or updates it if they
// had previously blocked the bot.
func PutUser(userID int, client *dynamodb.DynamoDB) error {

	user := structs.User{
		TelegramID:    userID,
		HasBlockedBot: false,
	}

	marshalledUser, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return errors.Errorf("PutUser: error while marshaling user: %v", err)
	}

	//Add the user only if the telegramID does not already exist.
	_, err = client.PutItem(&dynamodb.PutItemInput{
		ConditionExpression: aws.String("attribute_not_exists(TelegramID)"),
		TableName:           aws.String(user.Table()),
		Item:                marshalledUser,
	})

	if err != nil {
		// Casting to the awserr.Error type allows us to inspect the error
		// code returned by the service.
		if aerr, ok := err.(awserr.Error); ok {

			// Primary key condition failed.
			// Update the user to make sure HasBlockedUser is not true.
			if aerr.Code() == dynamodb.ErrCodeConditionalCheckFailedException {
				_ = UpdateUserBlockStatus(userID, false, client)
			}

		}
	}

	return err

}

// UpdateUserBlockStatus updates the HasBlockedUser field according to the input flag.
func UpdateUserBlockStatus(userID int, hasBlockedBot bool, client *dynamodb.DynamoDB) (err error) {

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":b": {
				BOOL: aws.Bool(hasBlockedBot),
			},
		},
		TableName: aws.String(structs.User{}.Table()),
		Key: map[string]*dynamodb.AttributeValue{
			"TelegramID": {
				N: aws.String(strconv.Itoa(userID)),
			},
		},
		UpdateExpression: aws.String("set HasBlockedBot = :b"),
	}

	_, err = client.UpdateItem(input)
	if err != nil {
		err = errors.Errorf("UpdateUserBlockStatus: unable to update user status: %s", err)
	}

	return

}

// IsUserAdmin returns true if the user's IsAdmin field is true.
func IsUserAdmin(userID int, client *dynamodb.DynamoDB) (isAdmin bool, err error) {

	output, err := client.GetItem(&dynamodb.GetItemInput{
		AttributesToGet: aws.StringSlice([]string{"IsAdmin"}),
		Key: map[string]*dynamodb.AttributeValue{
			"TelegramID": {
				N: aws.String(strconv.Itoa(userID)),
			},
		},
		TableName: aws.String(structs.User{}.Table()),
	})

	if err != nil {
		return
	}

	var user structs.User
	err = dynamodbattribute.UnmarshalMap(output.Item, &user)
	if err != nil {
		err = errors.Errorf("IsUserAdmin: failed to unmarshal user, %v", err)
	}

	isAdmin = user.IsAdmin
	return

}
