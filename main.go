package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"strings"
	"time"
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
	activeUsers := make(map[int]*time.Ticker)

	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	go func() {
		for update := range updates {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			command := strings.ToUpper(update.Message.Text)
			chatID := update.Message.Chat.ID
			_, userActive := activeUsers[chatID]

			switch command {
			case "START":
				if userActive {
					messages <- telegramMessage{chatID, "I've already started"}
				} else {
					activeUsers[chatID] = time.NewTicker(time.Second * 5)

					go func() {
						for range activeUsers[chatID].C {
							messages <- telegramMessage{chatID, "Do 10 pull ups"}
						}
					}()
				}

			case "STOP":
				if userActive {
					activeUsers[chatID].Stop()
				} else {
					messages <- telegramMessage{chatID, "I haven't started"}
				}

			}

		}
	}()

	for message := range messages {
		msg := tgbotapi.NewMessage(message.chatID, message.text)
		bot.Send(msg)
	}

}
