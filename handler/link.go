package handler

import (
	"net/url"

	"github.com/handybots/inzerobot/storage"
	tele "gopkg.in/tucnak/telebot.v3"
)

func (h handler) OnCreateLink(c tele.Context) error {
	return c.Send(h.lt.Text(c, "link"), tele.ForceReply)
}

func (h handler) onCreateLinkReply(c tele.Context) error {
	u, err := url.ParseRequestURI(c.Text())
	if err != nil || u.Scheme == "" || (u.Scheme != "http" && u.Scheme != "https") {
		return c.Send(h.lt.Text(c, "link"), tele.ForceReply)
	}

	link, err := h.db.Links.Create(storage.Link{
		UserID: c.Message().Sender.ID,
		URL:    c.Message().Text,
	})
	if err != nil {
		return err
	}

	return c.Reply(h.lt.String("domain") + "/" + link.String())
}

func (h handler) OnLinkList(c tele.Context) error {
	links, err := h.db.Links.ByUserID(c.Sender())
	if err != nil {
		return err
	}

	return c.Send(
		h.lt.Text(c, "my", links),
		tele.NoPreview)
}
