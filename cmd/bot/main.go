package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/demget/clickrus"
	"github.com/sirupsen/logrus"

	tele "gopkg.in/tucnak/telebot.v3"
	"gopkg.in/tucnak/telebot.v3/layout"
	"gopkg.in/tucnak/telebot.v3/middleware"

	"github.com/handybots/inzerobot/storage"
	handler "github.com/handybots/inzerobothandler"
)

const domain = "http://127.0.0.1/"

func main() {
	logrus.SetLevel(log.Lshortfile)

	layout.AddFunc("inc", func(i int) int {
		return i + 1
	})
	layout.AddFunc("format", func(n int) string {
		if n >= 1000000 {
			if n2 := n / 100000 % 10; n2 != 0 {
				return fmt.Sprint(n/1000000, ".", n2, " M")
			}
			return fmt.Sprint(n/1000000, " M")
		}
		if n >= 1000 {
			if n2 := n / 100 % 10; n2 != 0 {
				return fmt.Sprint(n/1000, ".", n2, " K")
			}
			return fmt.Sprint(n/1000, " K")
		}
		return fmt.Sprint(n)
	})
	layout.AddFunc("link", func(l string) string {
		u, err := url.ParseRequestURI(l)
		if err != nil {
			return "-"
		}
		return u.Host
	})

	lt, err := layout.New("bot.yml")
	if err != nil {
		log.Fatal(err)
	}

	b, err := tele.NewBot(lt.Settings())
	if err != nil {
		log.Fatal(err)
	}

	// db, err := storage.Open(os.Getenv("DB_URL"))
	db, err := storage.Open("tester:test123@/magiclinkbot?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	// ch, err := clickrus.NewHook(clickHouseConfig)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	// logger.AddHook(ch)

	h := handler.New(handler.Handler{
		Layout: lt,
		Bot:    b,
		DB:     db,
		Domain: domain,
	})

	// Middleware
	b.OnError = h.OnError
	b.Use(middleware.Logger(logger, h.LoggerFields))
	b.Use(lt.Middleware("ru", h.LocaleFunc))

	// Handlers
	b.Handle("/start", h.OnStart)
	b.Handle("/create", h.OnCreateLink)
	b.Handle("/my", h.OnLinkList)

	b.Handle(tele.OnText, h.OnText)
	b.Handle(tele.OnCallback, h.OnLang)

	b.Start()
}

var clickHouseConfig = clickrus.Config{
	Addr:    os.Getenv("CLICKHOUSE_URL"),
	Columns: []string{"date", "time", "level", "message", "event", "user_id"},
	Table:   "logs",
}
