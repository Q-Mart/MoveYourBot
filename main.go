package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

var activeUsers = make(map[int]*time.Ticker)

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

func motivate(chatID int, ch chan telegramMessage) {
	motivationalMessages := [3]string{"Do 10 pull ups!", "Do 10 press ups!", "Do 10 chin ups!"}
	i := 0

	for range activeUsers[chatID].C {
		ch <- telegramMessage{chatID, motivationalMessages[i%3]}
		i++
	}
}

func main() {
	bot, err := tgbotapi.NewBotAPI(getAccessToken())
	messages := make(chan telegramMessage)
	keyboard := tgbotapi.ReplyKeyboardMarkup{
		Keyboard:       [][]string{{"start"}, {"stop"}},
		ResizeKeyboard: true}

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
					messages <- telegramMessage{chatID, "Get ready to rumble!"}

					go motivate(chatID, messages)
				}

			case "STOP":
				if userActive {
					activeUsers[chatID].Stop()
					delete(activeUsers, chatID)
				} else {
					messages <- telegramMessage{chatID, "I haven't started"}
				}

			case "STATUS":
				if userActive {
					messages <- telegramMessage{chatID, "You are currently in a session"}
				} else {
					messages <- telegramMessage{chatID, "You are not currently in session"}
				}

			}

		}
	}()

	for message := range messages {
		msg := tgbotapi.NewMessage(message.chatID, message.text)
		msg.ReplyMarkup = keyboard
		bot.Send(msg)
	}

}
