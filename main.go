package main

import (
	"context"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gocolly/colly"
	_ "github.com/lib/pq"
	"golang.org/x/net/proxy"
	"log"
	"net"
	"net/http"
)

// Course stores information about a coursera course
//var (
//	titles = [...]string()
//)

func botCore(count int) {
	dialer, proxyErr := proxy.SOCKS5(
		"tcp",
		"185.162.235.83:1080",
		nil,
		proxy.Direct,
	)
	if proxyErr != nil {
		log.Panicf("Error in proxy %s", proxyErr)
	}
	client := &http.Client{Transport: &http.Transport{DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
		return dialer.Dial(network, addr)
	}}}
	bot, err := tgbotapi.NewBotAPIWithClient("944404078:AAG9Rk5JFkolvU4EwdSTXFqF2hnF3gLqBZQ", client)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	log.Printf("Logged on %s", bot.Self.UserName)
	msg := tgbotapi.NewMessage(37434600, "")
	msg.Text = fmt.Sprintln("Found ", count)
	bot.Send(msg)
}

func main() {
	// Instantiate default collector
	c := colly.NewCollector(
		colly.AllowedDomains("mangarock.com", "www.magarock.com"),
	)

	count := 0
	// On every a element which has class attribute call callback
	c.OnHTML("a[class]", func(e *colly.HTMLElement) {
		if e.Attr("class") == "_1A2Dc rZ05K" {
		count += 1
		}
	})

	c.Visit("https://mangarock.com/manga/mrs-serie-100159890")
	botCore(count)
	fmt.Println("found", count)
}