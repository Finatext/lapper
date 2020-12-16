package main

import (
	"github.com/nlopes/slack"
)

func postSlack(url string, afs []slack.AttachmentField) error {
	msg := slack.WebhookMessage{}
	msg.Username = "Lapper"
	msg.IconEmoji = ":mega:"
	msg.Text = "Notification from Lapper"
	msg.Attachments = []slack.Attachment{
		{
			Fields: afs,
		},
	}

	return slack.PostWebhook(url, &msg)
}

func buildSlackAttachmentFields(fName string, rid string, stdout string, stderr string, err error) []slack.AttachmentField {
	postStdout := getEnv("LAPPER_POST_STDOUT", "false")
	postStderr := getEnv("LAPPER_POST_STDERR", "true")
	postError := getEnv("LAPPER_POST_ERROR", "true")

	af := []slack.AttachmentField{}

	af = append(af, slack.AttachmentField{Title: "Function Name", Value: fName, Short: true})
	af = append(af, slack.AttachmentField{Title: "Request Id", Value: rid, Short: true})

	if postStdout == "true" {
		af = append(af, slack.AttachmentField{Title: "STDOUT", Value: stdout, Short: false})
	}

	if postStderr == "true" {
		af = append(af, slack.AttachmentField{Title: "STDERR", Value: stderr, Short: false})
	}

	if postError == "true" && err != nil {
		af = append(af, slack.AttachmentField{Title: "ERROR", Value: err.Error(), Short: false})
	}

	return af
}
