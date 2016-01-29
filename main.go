package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
)

type telegramMessage struct {
	chatID int
	text   string
}

func getAccessToken() string {
	buff, err := ioutil.ReadFile("access.token")
	if err != nil {
		log.Panic(err)
	}

	return string(buff[:len(buff)-1])
}

func main() {
	bot, err := tgbotapi.NewBotAPI(getAccessToken())
	messages := make(chan telegramMessage)

	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	activeUsers := make(map[int]countdownTimer)

	go func() {
		for update := range updates {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			command := update.Message.Text
			chatID := update.Message.Chat.ID

			if _, ok := activeUsers[chatID]; ok {
				userTimer := activeUsers[chatID]

				switch command {
				case "start":
					go userTimer.StartNew(messages)

				case "stop":
					log.Println("Stopping timer")
					userTimer.Stop()
					delete(activeUsers, chatID)

				}
			} else {
				activeUsers[chatID] = countdownTimer{chatID: chatID}
				userTimer := activeUsers[chatID]

				if command == "start" {
					go userTimer.StartNew(messages)
				}
			}

		}
	}()

	for message := range messages {
		msg := tgbotapi.NewMessage(message.chatID, message.text)
		bot.Send(msg)
	}
}
