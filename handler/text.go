package handler

import tele "gopkg.in/tucnak/telebot.v3"

func (h handler) OnText(c tele.Context) error {
	message := c.Message()
	if message.ReplyTo != nil {
		return h.onCreateLinkReply(c)
	}
	return nil
}
