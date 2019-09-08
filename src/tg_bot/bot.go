package main

import (
	"fmt"
	"giphy_api"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/net/proxy"
)

const (
	TOKEN = "YOUR_TOKEN"
	PROXY = "socks5://999.999.999.999:9999"
)


func getHelp() string {
	reply :=
		`This bot can search GIFs by your query. 
	Command must have format:
	/get_gif [result number limit] [query]
	You can result number limit set it to value from 0 to 25
	`
	return reply
}

func sendMsgs(chatId int64, msgs []string, bot *tgbotapi.BotAPI) {
	for _, msg := range msgs {
		if msg != "" {
			fmt.Println("", msg)
			newMsg := tgbotapi.NewMessage(chatId, msg)
			bot.Send(newMsg)
		}
	}
}

func main() {
	client := &http.Client{}
	if len(PROXY) > 0 {
		tgProxyURL, err := url.Parse(PROXY)
		if err != nil {
			log.Printf("Failed to parse proxy URL:%s\n", err)
		}
		tgDialer, err := proxy.FromURL(tgProxyURL, proxy.Direct)
		if err != nil {
			log.Printf("Failed to obtain proxy dialer: %s\n", err)
		}
		tgTransport := &http.Transport{
			Dial: tgDialer.Dial,
		}
		client.Transport = tgTransport
	}

	bot, err := tgbotapi.NewBotAPIWithClient(TOKEN, client)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		userMsg := update.Message.Text
		replySplit := strings.Split(userMsg, " ")
		var reply string
		switch replySplit[0] {
		case "/start", "/help":
			reply = getHelp()
		case "/get_gif":
			err_reply, res := api.GetGif(userMsg[1:])
			fmt.Println(res)
			if err_reply == "" {
				sendMsgs(update.Message.Chat.ID, res, bot)
			} else {
				reply = err_reply
			}
		default:
			reply = "Unknown command. Send /help to get available commands"
		}
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		bot.Send(msg)
	}
}
