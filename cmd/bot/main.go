package main

import (
	"log"
	"os"

	"github.com/demget/clickrus"
	"github.com/sirupsen/logrus"

	tele "gopkg.in/tucnak/telebot.v3"
	"gopkg.in/tucnak/telebot.v3/layout"
	"gopkg.in/tucnak/telebot.v3/middleware"

	"github.com/handybots/inzerobot/handler"
	"github.com/handybots/inzerobot/server"
	"github.com/handybots/inzerobot/storage"
)

func main() {
	logrus.SetLevel(log.Lshortfile)

	layout.AddFunc("inc", func(i int) int { return i + 1 })
	layout.AddFunc("format", formatN)
	layout.AddFunc("host", host)

	lt, err := layout.New("bot.yml")
	if err != nil {
		log.Fatal(err)
	}

	b, err := tele.NewBot(lt.Settings())
	if err != nil {
		log.Fatal(err)
	}

	db, err := storage.Open(os.Getenv("MYSQL_URL"))
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

	go func() {
		srv := server.New(b.Me.Username, ":8050", db)
		if err := srv.Listen(); err != nil {
			logrus.Fatalln(err)
		}
	}()

	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	// logger.AddHook(ch)

	h := handler.New(handler.Handler{
		Layout: lt,
		Bot:    b,
		DB:     db,
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
	b.Handle(lt.Button("lang"), h.OnLang)

	b.Start()
}

var clickHouseConfig = clickrus.Config{
	Addr:    os.Getenv("CLICKHOUSE_URL"),
	Columns: []string{"event", "user_id"},
	Table:   "inzerobot.logs",
}
