package handler

import tele "gopkg.in/tucnak/telebot.v3"

func (h handler) OnText(c tele.Context) error {
	if c.Message().ReplyTo != nil {
		return h.onCreateLinkReply(c)
	}
	return nil
}
