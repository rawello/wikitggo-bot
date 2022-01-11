package main

import (
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/trigun117/TelegramBot-Go/docker-compose/bot/code/wiki"
)

func telegramBot() {

	//Create bot
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		panic(err)
	}

	//Set update timeout
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	//Get updates from bot
	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		//Check if message from user is text
		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {

			switch update.Message.Text {
			case "/start":

				//Send message
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hi, i'm a wikipedia bot, i can search information in a wikipedia, send me something what you want find in Wikipedia.")
				bot.Send(msg)

			case "/number_of_users":

				if os.Getenv("DB_SWITCH") == "on" {

					//Assigning number of users to num variable
					num, err := getNumberOfUsers()
					if err != nil {

						//Send message
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Database error.")
						bot.Send(msg)
					}

					//Creating string which contains number of users
					ans := fmt.Sprintf("%d peoples used me for search information in Wikipedia", num)

					//Send message
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, ans)
					bot.Send(msg)
				} else {

					//Send message
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Database not connected, so i can't say you how many peoples used me.")
					bot.Send(msg)
				}
			default:

				//Set search language
				language := os.Getenv("LANGUAGE")

				//Create search url
				ms, _ := wiki.URLEncoded(update.Message.Text)

				url := ms
				request := "https://" + language + ".wikipedia.org/w/api.php?action=opensearch&search=" + url + "&limit=3&origin=*&format=json"

				//assigning value of answer slice to variable message
				message := wiki.WikipediaAPI(request)

				if os.Getenv("DB_SWITCH") == "on" {

					//Putting username, chat_id, message, answer to database
					if err := collectData(update.Message.Chat.UserName, update.Message.Chat.ID, update.Message.Text, message); err != nil {

						//Send message
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Database error, but bot still working.")
						bot.Send(msg)
					}
				}

				//Loop throug message slice
				for _, val := range message {

					//Send message
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, val)
					bot.Send(msg)
				}
			}
		} else {

			//Send message
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Use the words for search.")
			bot.Send(msg)
		}
	}
}

func main() {

	time.Sleep(1 * time.Minute)

	//Creating Table
	if os.Getenv("CREATE_TABLE") == "yes" {

		if os.Getenv("DB_SWITCH") == "on" {

			if err := createTable(); err != nil {

				panic(err)
			}
		}
	}

	time.Sleep(1 * time.Minute)

	//Call Bot
	telegramBot()
}
