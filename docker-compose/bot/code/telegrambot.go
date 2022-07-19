package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"os"
	"reflect"
	"time"

	"github.com/rawello/wikitggo-bot/docker-compose/bot/code/wiki"
)

func telegramBot() {

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {

			switch update.Message.Text {
			case "/start":

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "привет, напиши что ты хочешь найти на википедии")
				bot.Send(msg)

			case "/number_of_users":

				if os.Getenv("DB_SWITCH") == "on" {

					num, err := getNumberOfUsers()
					if err != nil {

						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "дб упала.")
						bot.Send(msg)
					}

					ans := fmt.Sprintf("%d peoples used me for search information in Wikipedia", num)

					msg := tgbotapi.NewMessage(update.Message.Chat.ID, ans)
					bot.Send(msg)
				} else {

					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "хз скок человек, дб неконнект")
					bot.Send(msg)
				}
			default:

				language := os.Getenv("LANGUAGE")

				ms, _ := wiki.URLEncoded(update.Message.Text)

				url := ms
				request := "https://" + language + ".wikipedia.org/w/api.php?action=opensearch&search=" + url + "&limit=3&origin=*&format=json"

				message := wiki.WikipediaAPI(request)

				if os.Getenv("DB_SWITCH") == "on" {

					if err := collectData(update.Message.Chat.UserName, update.Message.Chat.ID, update.Message.Text, message); err != nil {

						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "ошибка дб, напиши @rawello, бот работает")
						bot.Send(msg)
					}
				}

				for _, val := range message {

					msg := tgbotapi.NewMessage(update.Message.Chat.ID, val)
					bot.Send(msg)
				}
			}
		} else {

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "используй слова для поиска")
			bot.Send(msg)
		}
	}
}

func main() {

	time.Sleep(1 * time.Minute)

	if os.Getenv("CREATE_TABLE") == "yes" {

		if os.Getenv("DB_SWITCH") == "on" {

			if err := createTable(); err != nil {

				panic(err)
			}
		}
	}

	time.Sleep(1 * time.Minute)

	telegramBot()
}
