package handler

import (
	"strconv"

	tele "gopkg.in/tucnak/telebot.v3"
)

func (h handler) OnLink(c tele.Context) error {
	args := c.Args()
	if len(args) != 2 {
		return nil
	}

	num, _ := strconv.Atoi(args[0])
	id, _ := strconv.ParseInt(args[1], 10, 64)

	exists, err := h.db.Links.ExistsForUser(id, c.Sender())
	if err != nil {
		return err
	}

	if !exists {
		return c.Edit(h.lt.Text(c, "whence"))
	}

	// TODO: Get the link from cache.

	link, err := h.db.Links.ByID(id)
	if err != nil {
		return err
	}

	stats, err := h.db.Views.Stats(link.ID)
	if err != nil {
		return err
	}

	xlink := Link{
		Link:      link,
		ViewStats: stats,
		Number:    num,
	}

	return c.Edit(
		h.lt.Text(c, "my_link", xlink),
		h.lt.Markup(c, "delete", link.ID),
		tele.NoPreview)
}

func (h handler) OnLinkDelete(c tele.Context) error {
	data := c.Data()
	if data == "-1" {
		return h.OnMy(c)
	}

	id, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return err
	}

	link, err := h.db.Links.ByID(id)
	if err != nil {
		return err
	}

	if link.UserID != c.Callback().Sender.ID {
		return c.Send(h.lt.Text(c, "whence"))
	}

	if err := h.db.Links.SetDeleted(link.ID, true); err != nil {
		return err
	}

	return h.OnMy(c)
}
