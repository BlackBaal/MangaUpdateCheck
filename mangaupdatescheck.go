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

type Data struct {
	id    int
	title string //Full name
	link  string //Link to page for scraping
	value int    //Number of chapters currently saved in database
}
func changeCount(id int, count int, db *sql.DB) error {
	query := fmt.Sprintf("UPDATE database SET VALUE = %d WHERE ID = %d", count, id)
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func botCore(link string, dif int) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
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

func dbOpen() (db *sql.DB) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	return db
}
func main() {
	var (
		count = 0    //Number of found chapters
		db *sql.DB
		data = Data{}
	)

	//Open database
	db = dbOpen()

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
		err := res.Scan(&data.id, &data.title, &data.link, &data.value)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(data.id, data.title, data.link, data.value)
		c.Visit(data.link)
		if count > data.value {
			// Send notification when new chapters are available
			botCore(data.link, count - data.value)
			//Change chapter count inside the db
			changeCount(data.id, count, db)
		}
		count = 0
	}
	db.Close()
}
