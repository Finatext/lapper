package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

func HandleEvent(ctx context.Context, payload json.RawMessage) (string, error) {
	lc, _ := lambdacontext.FromContext(ctx)
	os.Setenv("AWS_REQUEST_ID", lc.AwsRequestID)

	notifyCond := getEnv("LAPPER_NOTIFY_COND", "stderr")
	debug := getEnv("DEBUG", "false")

	slackWebhookURL := os.Getenv("LAPPER_SLACK_WEBHOOK_URL")
	handler := strings.Split(os.Getenv("_HANDLER"), " ")

	command := handler[0]
	args := handler[1:]

	if debug == "true" {
		fmt.Println("[Lapper] Command: ", command, strings.Join(args, " "))
		fmt.Println("[Lapper] Payload: ", string(payload))
	}

	fn := NewFunction(command, args, payload)
	stdout, stderr, err := fn.Run()

	if slackWebhookURL != "" &&
		(notifyCond == "all" ||
			(notifyCond == "stderr" && stderr != "") ||
			(notifyCond == "exitcode" && err != nil)) {

		af := buildSlackAttachmentFields(os.Getenv("AWS_LAMBDA_FUNCTION_NAME"), lc.AwsRequestID, stdout, stderr, err)
		postSlack(slackWebhookURL, af)
	}

	if err != nil {
		return fmt.Sprintf("Failed"), err
	}

	return fmt.Sprintf("Succeeded"), nil
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func main() {
	if flag.Lookup("test.v") != nil {
		return
	}

	lambda.Start(HandleEvent)
}
