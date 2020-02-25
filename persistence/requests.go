// Copyright 2019 Alessandro Pomponio. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package persistence contains methods to store and retrieve
// content from Amazon DynamoDB.
package persistence

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/pkg/errors"
	"github.com/rs/xid"

	"github.com/AlessandroPomponio/serverless-amazon-refbot/structs"
)

// GetRequestsSince returns the request that happened after a certain threshold.
// The threshold is given in Unix timestamp format.
func GetRequestsSince(threshold int64, client *dynamodb.DynamoDB) (requests []structs.Request, err error) {

	filter := expression.Name("UnixTime").GreaterThan(expression.Value(threshold))
	projection := expression.NamesList(expression.Name("TelegramID"), expression.Name("Time"), expression.Name("URL"))

	expr, err := expression.NewBuilder().WithFilter(filter).WithProjection(projection).Build()
	if err != nil {
		err = errors.Errorf("GetRequestsSince: error while building the expression: %s", err)
		return
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(structs.Request{}.Table()),
	}

	result, err := client.Scan(params)
	if err != nil {
		err = errors.Errorf("GetRequestsSince: error while querying the database: %s", err)
		return
	}

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &requests)
	if err != nil {
		err = errors.Errorf("GetRequestsSince: error while unmarshaling the results: %s", err)
	}

	return

}

// PutRequest saves a request on DynamoDB.
func PutRequest(userID int, url string, client *dynamodb.DynamoDB) error {

	request := structs.Request{
		XID:        xid.New().String(),
		TelegramID: userID,
		URL:        url,
		Time:       time.Now(),
		UnixTime:   time.Now().Unix(),
	}

	marshalledRequest, err := dynamodbattribute.MarshalMap(request)
	if err != nil {
		return errors.Errorf("PutRequest: error while marshaling request: %v", err)
	}

	_, err = client.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(request.Table()),
		Item:      marshalledRequest,
	})

	return err

}
