package main

import (
	"log"
	"os"

	"github.com/demget/clickrus"
	"github.com/sirupsen/logrus"

	tele "gopkg.in/tucnak/telebot.v3"
	"gopkg.in/tucnak/telebot.v3/layout"
	"gopkg.in/tucnak/telebot.v3/middleware"

	"github.com/handybots/nullifybot/handler"
	"github.com/handybots/nullifybot/server"
	"github.com/handybots/nullifybot/storage"
)

func main() {
	logrus.SetLevel(log.Lshortfile)

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

	ch, err := clickrus.NewHook(clickHouseConfig)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		srv := server.New(b.Me.Username, "127.0.0.1:8050", db)
		if err := srv.Listen(); err != nil {
			logrus.Fatalln(err)
		}
	}()

	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.AddHook(ch)

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
	b.Handle(tele.OnText, h.OnText)
	b.Handle(lt.Callback("lang"), h.OnLang)
	b.Handle(lt.Callback("link"), h.OnLink)
	b.Handle(lt.Callback("link_delete"), h.OnLinkDelete)

	// Locale-dependent handlers.
	for _, loc := range []string{"ru", "en"} {
		b.Handle(lt.ButtonLocale(loc, "my"), h.OnMy)
	}

	b.Start()
}

var clickHouseConfig = clickrus.Config{
	Addr:    os.Getenv("CLICKHOUSE_URL"),
	Columns: []string{"event", "user_id"},
	Table:   "nullifybot.logs",
}
