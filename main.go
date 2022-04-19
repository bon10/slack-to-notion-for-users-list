package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	api := slack.New(os.Getenv("SLACK_TOKEN"))

	users, err := api.GetUsers()
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	totalCount := 0
	resCount := 0

	// TODO: goroutine化したい
	for i := 2; i < len(users); i++ {
		fmt.Printf("ユーザー情報%#v\n", users[i])
		if users[i].IsBot != true && users[i].Profile.DisplayName != "" && users[i].Deleted != true {
			user := generateUser(users[i])
			resCount += NotionCreatePage(user)
			totalCount += 1
		} else {
			if users[i].Deleted == true {
				fmt.Printf("%vさんは削除されたユーザーです\n", users[i].Profile.DisplayName)
			} else {
				if users[i].IsBot != true && users[i].RealName != "" {
					users[i].Profile.DisplayName = users[i].RealName
					user := generateUser(users[i])
					resCount += NotionCreatePage(user)
					totalCount += 1
				} else {
					fmt.Printf("%vさんはdisplayNameがないかBotユーザーです\n", users[i].Profile.DisplayName)
				}
			}
		}
	}
	fmt.Printf("%d件中%d件処理しました", totalCount, resCount)
	os.Exit(0)
}

func generateUser(user slack.User) User {
	fmt.Println(user.Profile.DisplayName)

	// いろんな記号を置換
	displayName := strings.ReplaceAll(user.Profile.DisplayName, ",", "/")
	displayName = strings.ReplaceAll(displayName, "、", "/")
	displayName = strings.ReplaceAll(displayName, "／", "/")
	userAttributes := strings.Split(displayName, "/")

	fmt.Println(userAttributes)
	name := userAttributes[0]
	organizations := []string{}
	if len(userAttributes) > 2 {
		for i := 1; i < len(userAttributes); i++ {
			organizations = append(organizations, userAttributes[i])
		}
	} else if len(userAttributes) > 1 {
		organizations = append(organizations, userAttributes[1])
	} else {
		fmt.Printf("%sさんは組織が存在しないようですよ\n", user.Profile.DisplayName)
	}

	notionUser := &User{
		Name:          name,
		Organizations: organizations,
		ImageURL:      user.Profile.Image512,
		Email:         user.Profile.Email,
	}

	return *notionUser
}
