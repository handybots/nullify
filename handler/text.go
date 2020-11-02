package handler

import (
	"strings"

	"github.com/handybots/nullifybot/storage"
	tele "gopkg.in/tucnak/telebot.v3"
)

func (h handler) OnText(c tele.Context) error {
	count, err := h.db.Links.CountByUser(c.Sender())
	if err != nil {
		return err
	}

	if count >= h.lt.Int("limit") {
		return c.Send(h.lt.Text(c, "limit"))
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

	return c.Reply(h.lt.Text(c, "created", link.String()))
}
