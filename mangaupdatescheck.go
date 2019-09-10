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
	fmt.Println(id)
	query := fmt.Sprintf("UPDATE database SET VALUE = %d WHERE ID = %d", count, id)
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	log.Println(id, count)
	return nil
}

func dbConnect() (info string) {
	var host = "ec2-54-75-245-196.eu-west-1.compute.amazonaws.com"
	var port = "5432"
	var user = "paxfbkluuhzzxd"
	var password = "e3dd8f47e1e2bc1b0c266248f5a449ebcc8f19dcca546091fe25087b386321d7"
	var dbname = "d6knbc375gophn"
	var sslmode = ""
	var dbInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)
	return dbInfo
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

	//client := &http.Client{Transport: &http.Transport{DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
	//	return dialer.Dial(network, addr)
	//}}}
	//bot, err := tgbotapi.NewBotAPIWithClient("944404078:AAG9Rk5JFkolvU4EwdSTXFqF2hnF3gLqBZQ", client)
	bot, err := tgbotapi.NewBotAPI("944404078:AAG9Rk5JFkolvU4EwdSTXFqF2hnF3gLqBZQ")
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Logged on %s", bot.Self.UserName)
	//baseURL := "https://mangaupdatescheck.herokuapp.com/"
	//url := baseURL + bot.Token
	//_, err = bot.SetWebhook(tgbotapi.NewWebhook(url))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//updates := bot.ListenForWebhook("/" + bot.Token)
	//
	//for update := range updates {
	//	log.Printf("%+v\n", update)
	//}
	msg := tgbotapi.NewMessage(37434600, "")
	if dif > 1 {
		msg.Text = fmt.Sprintln(link, "\n", dif, "new chapters")
	} else {
		msg.Text = fmt.Sprintln(link, "\n", dif, "new chapter")
	}
	bot.Send(msg)
}
//func hello(w http.ResponseWriter, r *http.Request) {
//	fmt.Fprintln(w, "Hello World")
//}

func main() {
	var (
		id    int
		title string //Full name of the manga
		link  string //Link to manga
		value int    //Number of chapters currently saved in database
		count = 0    //Number of found chapters
	)

	//Open database
	db, err := sql.Open("postgres", dbConnect())
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
		if count == value {
			// Send notification that new chapters are available
			botCore(link, count-value)
			//Change chapter count inside th db
			changeCount(id, count, db)
		}
		count = 0
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	//bot, err := tgbotapi.NewBotAPI("944404078:AAG9Rk5JFkolvU4EwdSTXFqF2hnF3gLqBZQ")
	//if err != nil {
	//	log.Panic(err)
	//}
	//bot.Debug = true
	//log.Printf("Logged on %s", bot.Self.UserName)
	//baseURL := "https://mangaupdatescheck.herokuapp.com/"
	//url := baseURL + bot.Token
	//_, err = bot.SetWebhook(tgbotapi.NewWebhook(url))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//updates := bot.ListenForWebhook("/" + bot.Token)
	//http.ListenAndServe(":"+port, nil)
	//for update := range updates {
	//	log.Printf("%+v\n", update)
	//}

}
