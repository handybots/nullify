package handler

import tele "gopkg.in/tucnak/telebot.v3"

func (h handler) OnMore(c tele.Context) error {
	return c.Edit(
		h.lt.Text(c, "subscribe"),
		h.lt.Markup(c, "subscribe"),
		tele.NoPreview)
}

func (h handler) OnSubscribe(c tele.Context) error {
	if h.checkSubscription(c.Sender()) {
		return c.Edit(h.lt.Text(c, "thanks"))
	}
	return c.Respond()
}

func (h handler) checkSubscription(r tele.Recipient) bool {
	chat := &tele.Chat{
		Type:     tele.ChatChannel,
		Username: "@" + h.lt.String("subscribe"),
	}

	member, err := h.b.ChatMemberOf(chat, r.(*tele.User))
	if err != nil {
		h.OnError(err, nil)
		return false
	}

	switch member.Role {
	case tele.Creator, tele.Administrator, tele.Member:
		return true
	default:
		return false
	}
}
