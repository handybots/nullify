package handler

import (
	"strings"

	"github.com/handybots/nullify/storage"
	tele "gopkg.in/tucnak/telebot.v3"
)

func (h handler) OnText(c tele.Context) error {
	count, err := h.db.Links.CountByUser(c.Sender())
	if err != nil {
		return err
	}

	limits := h.lt.Get("limits")
	if count >= limits.Int("subscribed") {
		return c.Send(h.lt.Text(c, "limit"))
	}

	if count >= limits.Int("default") &&
		!h.checkSubscription(c.Sender()) {
		return c.Send(
			h.lt.Text(c, "limit"),
			h.lt.Markup(c, "more"),
		)
	}

	var ent tele.MessageEntity
	for _, ent = range c.Message().Entities {
		if ent.Type == tele.EntityURL {
			break
		}
	}
	if ent.Type != tele.EntityURL {
		return c.Reply(h.lt.Text(c, "bad"))
	}

	url := c.Text()[ent.Offset : ent.Offset+ent.Length]
	if !strings.HasPrefix(url, "http://") &&
		!strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	link := storage.Link{
		UserID: c.Message().Sender.ID,
		URL:    url,
	}

	link, err = h.db.Links.Create(link)
	if err != nil {
		return err
	}

	return c.Reply(
		h.lt.Text(c, "created", link.String()),
		tele.NoPreview,
	)
}
