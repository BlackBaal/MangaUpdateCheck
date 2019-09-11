package main

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gocolly/colly"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func changeCount(id int, count int, db *sql.DB) error {
	query := fmt.Sprintf("UPDATE database SET VALUE = %d WHERE ID = %d", count, id)
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func botCore(link string, dif int) {

	//dialer, proxyErr := proxy.SOCKS5(
	//	"tcp",
	//	"62.112.11.204:80",
	//	nil,
	//	proxy.Direct,
	//)
	//if proxyErr != nil {
	//	log.Panicf("Error in proxy %s", proxyErr)
	//}
	//
	//client := &http.Client{Transport: &http.Transport{DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
	//	return dialer.Dial(network, addr)
	//}}}
	//bot, err := tgbotapi.NewBotAPIWithClient("944404078:AAG9Rk5JFkolvU4EwdSTXFqF2hnF3gLqBZQ", client)
	bot, err := tgbotapi.NewBotAPI("944404078:AAG9Rk5JFkolvU4EwdSTXFqF2hnF3gLqBZQ")
	if err != nil {
		log.Panic(err)
	}
	//bot.Debug = true
	log.Printf("Logged on %s", bot.Self.UserName)
	msg := tgbotapi.NewMessage(37434600, "")
	if dif > 1 {
		msg.Text = fmt.Sprintln(link, "\n", dif, "new chapters")
	} else {
		msg.Text = fmt.Sprintln(link, "\n", dif, "new chapter")
	}
	bot.Send(msg)
}

func main() {
	var (
		id    int
		title string //Full name
		link  string //Link to page for scraping
		value int    //Number of chapters currently saved in database
		count = 0    //Number of found chapters
	)

	//Open database
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Instantiate default collector
	c := colly.NewCollector(
		colly.AllowedDomains("mangarock.com", "www.magarock.com"),
	)

	// On every a element which has class attribute call callback
	c.OnHTML("a[class]", func(e *colly.HTMLElement) {
		if e.Attr("class") == "_1A2Dc rZ05K" {
			count += 1
		}
	})

	res, _ := db.Query("select * from database")
	for res.Next() {
		err := res.Scan(&id, &title, &link, &value)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(id, title, link, value)
		c.Visit(link)
		if count > value {
			// Send notification when new chapters are available
			botCore(link, count-value)
			//Change chapter count inside the db
			changeCount(id, count, db)
		}
		count = 0
	}

}
