package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
)

func slackClient() *slack.Client {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	api := slack.New(os.Getenv("SLACK_TOKEN"))
	return api
}
