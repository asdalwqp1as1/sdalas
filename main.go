package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	port     = "8088"
	timeZone = "Asia/Almaty"
	// chatID   = -4102987152
	chatID = -4221126429
	token  = "7248085977:AAH26NYJUreuju8d16HX3_FoIGxdRH8yjl0"
)

var (
	count    = 0
	daysLeft = 10
)

func getDaysLeft() int {
	var (
		today = time.Now().Weekday()
	)

	if today == time.Monday {
		switch count {
		case 0:
			count++
			return daysLeft
		case 1:
			count++
			daysLeft--
			return daysLeft
		case 2:
			count = 1
			daysLeft = 10
			return daysLeft
		}
	}

	daysLeft--
	return daysLeft
}

func sendReminder(bot *tgbotapi.BotAPI, chatID int64) {
	var (
		daysLeft = getDaysLeft()
		message  string
	)

	if daysLeft == -1 {
		return
	} else if daysLeft == 1 {
		message = "❗❗❗ ПОСЛЕДНИЙ ДЕНЬ СПРИНТА ❗❗❗\n\n\nСсылка на встречу:\nhttps://meet.google.com/jqu-kwiv-nmt"
	} else if daysLeft >= 2 && daysLeft <= 4 {
		message = fmt.Sprintf("❗❗❗ ДО КОНЦА СПРИНТА ОСТАЛОСЬ %d ДНЯ ❗❗❗\n\n\nСсылка на встречу:\nhttps://meet.google.com/jqu-kwiv-nmt", daysLeft)
	} else if daysLeft >= 5 && daysLeft <= 10 {
		message = fmt.Sprintf("❗❗❗ ДО КОНЦА СПРИНТА ОСТАЛОСЬ %d ДНЕЙ ❗❗❗\n\n\nСсылка на встречу:\nhttps://meet.google.com/jqu-kwiv-nmt", daysLeft)
	}

	msg := tgbotapi.NewMessage(chatID, message)
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("ERROR: failed send message: %v\n", err)
	}
}

func commentCards(bot *tgbotapi.BotAPI, chatID int64) {
	message := "❗ НЕ ЗАБУДЬТЕ ОСТАВИТЬ КОММЕНТАРИИ К ЗАДАЧАМ ❗\n\n\nСсылка на ClickUp:\nhttps://app.clickup.com/24579426/v/b/li/901802397345"

	msg := tgbotapi.NewMessage(chatID, message)
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("ERROR: failed send message: %v\n", err)
	}
}

func main() {
	log.Println("Application running ...")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	bot.Debug = true
	// sendReminder(bot, int64(chatID))
	// commentCards(bot, int64(chatID))

	location, err := time.LoadLocation(timeZone)
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	// fmt.Println(location.String())

	cron := gocron.NewScheduler(location)
	cron.Every(1).Monday().Tuesday().Wednesday().Thursday().Friday().At("9:50").Do(sendReminder, bot, int64(chatID))
	cron.Every(1).Monday().Tuesday().Wednesday().Thursday().Friday().At("17:50").Do(commentCards, bot, int64(chatID))
	cron.StartAsync()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Bot is running")
		})

		log.Printf("Starting server on port %s...\n", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatalf("Failed to start server: %v\n", err)
		}
	}()

	<-c

	cron.Stop()
	log.Println("Application closing ...")
}
