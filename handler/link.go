package handler

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/handybots/nullifybot/storage"
	tele "gopkg.in/tucnak/telebot.v3"
)

const limitLinks = 5

func (h handler) limitLinks(c tele.Context) error {
	count, err := h.db.Links.CountByUserID(c.Chat())
	if err != nil {
		return err
	}
	if count >= limitLinks {
		c.Send(h.lt.Text(c, "link_limit"))
		return fmt.Errorf("create links: limit")
	}

	return nil
}

func (h handler) OnCreateLink(c tele.Context) error {
	if err := h.limitLinks(c); err != nil {
		return err
	}
	return c.Send(h.lt.Text(c, "link"), tele.ForceReply)
}

func (h handler) onCreateLinkReply(c tele.Context) error {
	if err := h.limitLinks(c); err != nil {
		return err
	}

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

	if len(links) == 0 {
		return c.Send(h.lt.Text(c, "my_empty"))
	}

	menu := h.b.NewMarkup()
	menu.ResizeKeyboard = true

	type btnInfo struct {
		ID     string
		Number string
	}

	buttons := make([]tele.Btn, 0)
	for idx, link := range links {
		buttons = append(buttons, *h.lt.Button(c, "link", btnInfo{
			ID:     strconv.FormatInt(link.ID, 10),
			Number: strconv.Itoa(idx + 1),
		}))
	}
	menu.Inline(menu.Row(buttons...))

	return c.Send(
		h.lt.Text(c, "my", links),
		menu,
		tele.NoPreview)
}

func (h handler) OnDeleteLink(c tele.Context) error {
	id, err := strconv.ParseInt(c.Callback().Data, 10, 64)
	if err != nil {
		return err
	}

	link, err := h.db.Links.ByID(id)
	if err != nil {
		return err
	}

	confirm := h.lt.Markup(c, "link_delete")
	confirm.InlineKeyboard[0][0].Data = c.Callback().Data

	return c.Send(
		h.lt.Text(c, "delete_link", link.URL),
		confirm,
		tele.NoPreview,
	)
}

func (h handler) OnDeleteLinkConfirm(c tele.Context) error {
	defer c.Delete()
	if c.Callback().Data == "-1" {
		return nil
	}

	id, err := strconv.ParseInt(c.Callback().Data, 10, 64)
	if err != nil {
		return err
	}

	link, err := h.db.Links.ByID(id)
	if err != nil {
		return err
	}

	if link.UserID != c.Callback().Sender.ID {
		return c.Send(h.lt.Text(c, "where_link"))
	}

	link.Deleted = true
	if err := h.db.Links.SetDeleted(link); err != nil {
		return err
	}

	return h.OnLinkList(c)
}
