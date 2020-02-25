# Go Serverless Amazon Refbot

[![Build Status](https://travis-ci.org/AlessandroPomponio/serverless-amazon-refbot.svg?branch=master)](https://travis-ci.org/AlessandroPomponio/serverless-amazon-refbot)
[![Go Report Card](https://goreportcard.com/badge/github.com/AlessandroPomponio/serverless-amazon-refbot)](https://goreportcard.com/report/github.com/AlessandroPomponio/serverless-amazon-refbot)
[![GoDoc](https://godoc.org/github.com/AlessandroPomponio/serverless-amazon-refbot?status.svg)](https://godoc.org/github.com/AlessandroPomponio/serverless-amazon-refbot)

This project allows you to run a bot on the Amazon Lambda serverless platform to generate referral links.

## Configuration

### Getting a Bitly API key

1. Go to Bitly's website [bitly.com](https://bitly.com) and login or create an account.
2. Click on your account name in the top right and choose `Group settings`.
3. Next up, choose `Advanced settings` and click on `API support`.
4. In the second paragraph you'll find a link to `User settings`, click on it.
5. Input your password and click on `Generate token` at the bottom.
6. Save this token, we'll need it later.

### DynamoDB configuration

1. Go to DynamoDB's web page: [https://console.aws.amazon.com/dynamodb/](https://console.aws.amazon.com/dynamodb/)
2. Press the "Create Table" button in the top of the page.
3. Create the `Users` table and use `TelegramID` as the partition key, `Number` type (leave "use default settings" ticked).
4. Create the `Requests` table and use `XID` as the partition key, `String` type (again, leave "use default settings" ticked).
5. Depending on your use, you may want to turn off the provisioning for the tables.

### IAM configuration

1. Go to IAM's web page: [https://console.aws.amazon.com/iam](https://console.aws.amazon.com/iam)
2. Press on `Roles` in the menu on the left.
3. Press the `Create role` button.
4. Choose `Lambda` from the list and press the `next` button.
5. Choose the `AWSLambdaBasicExecutionRole` role from the list.
6. Press the `Next` button.
7. Input the tags you want or leave blank if you don't need them.
8. Press again the `Next` button.
9. Input the role name in the form and press the `Create role` button.
10. Choose the role you just created from the list.
11. Press the `Add inline policy` text on the right of the screen.
12. Click the `Choose service` text and choose `DynamoDB`.
13. Select the following operations:
    - Read
        - BatchGetItem
        - ConditionCheckItem
        - GetItem
        - Query
        - Scan
    - Write
        - PutItem
        - UpdateItem
14. Click `Resources`, find the `Table` entry and click `Add ARN`.
15. Fill in the ARN with the ARNs of the tables you created earlier, you will find them at the bottom of the table page.
16. Do the same thing for the `Index` fields: `ARN/index/theIndexesWeChoseEarlier`.
17. Press `Verify policy`.
18. Name the policy.
19. Create the role.

### Lambda function creation

1. Go to Lambda's web page: [https://console.aws.amazon.com/lambda](https://console.aws.amazon.com/lambda)
2. Press `Create function` on the top right.
3. Choose `Author from scratch`.
4. Name your function and choose the `Go 1.x` runtime.
5. In the execution roles, choose `Use an existing role` and choose the one we created earlier.
6. Once the function has been created, fill in the following environment variables:
   - `AMAZON_DOMAIN`: the domain for which you want to use the bot (e.g. amazon.it).
   - `BITLY_KEY`: your Bitly API key.
   - `REF_ID`: your referral id from the Amazon affiliates program.
   - `REQUEST_TABLE_NAME`: the name you gave to the Requests table.
   - `TG_KEY`: a bot token from Telegram's [BotFather](https://t.me/BotFather).
   - `USER_TABLE_NAME`: the name you gave to the Users table.
7. Write `main` as the function handler.

### API Gateway configuration

1. Go to the API Gateway's web page: [https://console.aws.amazon.com/apigateway](https://console.aws.amazon.com/apigateway)
2. Go to API and choose `Create API`.
3. Choose `New API` and use a `Regional` endpoint.
4. Click on the newly created API and, from the dropdown `Actions` menu, choose `Create Method`.
5. Choose the `POST` method and confirm by pressing on the tick.
6. Make sure that `Lambda function` is selected as the `Integration type`.
7. Make sure that `Lambda Proxy Integration` is **disabled**.
8. Choose the appropriate region and write name of the function you've created in the `Lambda function` field.
9. Make sure that in the `Body mapping templates` of the function, `When there are no templates defined (recommended)"` is selected.
10. Deploy the API by choosing the option from the dropdown menu. This way you'll be given the URL we'll use to set up the bot's webhooks.

## Compiling

Now that we have (finally) set everything up, we can compile. To do so, we need to get Amazon's Go SDK with

```bash
go get -u github.com/aws/aws-sdk-go
```

### Linux

To compile on Linux we need to run:

```bash
GOOS=linux go build main.go
zip function.zip main
```

### Windows

On Windows, instead:

```cmd
set GOOS=linux
go build -o main main.go
%USERPROFILE%\Go\bin\build-lambda-zip.exe -o main.zip main
```

> I have included the `buildzip.bat` script that does exactly this.

You can now upload the function via the web interface and save the changes.

## Webhook setup

From the Lambda page, get the API Endpoint and from Telegram your bot token.
Perform the appropriate CURL request.

### Webhook creation

```bash
curl --request POST --url https://api.telegram.org/bot<BOT-TOKEN>/setWebhook --header 'content-type: application/json' --data '{"url": "<API-GATEWAY-URL>"}'
```

### Webhook deletion

```bash
curl --request POST --url https://api.telegram.org/bot<BOT-TOKEN>/setWebhook --header 'content-type: application/json' --data '{"url": ""}'
```

## Disclaimer

This project is meant for educational purposes.
I am not responsible for possible breaches of Amazon Affiliates's ToS, accidental charges from Amazon Web Services, etc.
By using this software, you take all the responsibilities.
