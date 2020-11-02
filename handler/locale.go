package handler

import tele "gopkg.in/tucnak/telebot.v3"

func (h handler) LocaleFunc(r tele.Recipient) string {
	locale, _ := h.db.Users.Lang(r)
	return locale
}
