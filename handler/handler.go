package handler

import (
	"github.com/handybots/inzerobot/storage"
	tele "gopkg.in/tucnak/telebot.v3"
	"gopkg.in/tucnak/telebot.v3/layout"
)

func New(c Handler) handler {
	return handler{
		lt:     c.Layout,
		b:      c.Bot,
		db:     c.DB,
		domain: c.Domain,
	}
}

type (
	Handler struct {
		Layout *layout.Layout
		Bot    *tele.Bot
		DB     *storage.DB
		Domain string
	}
	handler struct {
		lt     *layout.Layout
		b      *tele.Bot
		db     *storage.DB
		domain string
	}
)
