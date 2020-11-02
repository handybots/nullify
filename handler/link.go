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
	l := generateLink()

	if err := h.db.Links.Create(storage.Link{
		UserID:     c.Message().Sender.ID,
		SourceLink: c.Message().Text,
		Link:       l,
	}); err != nil {
		return err
	}

	return c.Reply(h.domain + l)
}

func (h handler) OnLinkList(c tele.Context) error {
	links, err := h.db.Links.ByUserID(c.Message().Sender.ID)
	if err != nil {
		return err
	}
	return c.Send(h.lt.Text(c, "my", links), tele.NoPreview)
}
