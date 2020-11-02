package handler

import (
	"log"

	tele "gopkg.in/tucnak/telebot.v3"
)

func (h handler) OnStart(c tele.Context) error {
	chat := c.Chat()

	has, err := h.db.Users.Exists(chat)
	if err != nil {
		return err
	}

	if !has {
		log.Println("Start from", chat.Recipient())
		if err := h.db.Users.Create(chat); err != nil {
			return err
		}
	}

	link, err := h.db.Links.ByID(1)
	if err != nil {
		return err
	}

	defer c.Send(
		h.lt.Text(c, "send"),
		h.lt.Markup(c, "menu"))

	return c.Send(
		h.lt.Text(c, "start", link.String()),
		h.lt.Markup(c, "lang"),
		tele.NoPreview)
}

func (h handler) OnLang(c tele.Context) error {
	lang := c.Data()
	if locale, _ := h.lt.Locale(c); locale == lang {
		return c.Respond()
	}

	if err := h.db.Users.SetLang(c.Sender(), lang); err != nil {
		return err
	} else {
		h.lt.SetLocale(c, lang)
	}

	return c.Edit(
		h.lt.Text(c, "start"),
		h.lt.Markup(c, "lang"))
}
