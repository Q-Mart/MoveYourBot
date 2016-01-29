package main

import (
	"log"
	"time"
)

type countdownTimer struct {
	secondsElapsed, chatID int
	ticker                 *time.Ticker
}

func (ct *countdownTimer) Stop() {
	ct.ticker.Stop()
}

func (ct *countdownTimer) StartNew(ch chan telegramMessage) {
	ct.secondsElapsed = 0
	ct.ticker = time.NewTicker(time.Second)
	log.Print("Timer started!")

	go func() {
		for range ct.ticker.C {
			if ct.secondsElapsed == 5 {
				ch <- telegramMessage{ct.chatID, "Do 10 pull ups!"}
				ct.secondsElapsed = 0
			}
			ct.secondsElapsed++
		}
	}()

}
