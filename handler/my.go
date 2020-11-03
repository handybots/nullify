package handler

import (
	"database/sql"

	tele "gopkg.in/tucnak/telebot.v3"
)

func (h handler) OnMy(c tele.Context) error {
	links, err := h.db.Links.ByUser(c.Sender())
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if len(links) == 0 {
		return c.Send(h.lt.Text(c, "send"))
	}

	var xlinks []Link
	for i, link := range links {
		stats, err := h.db.Views.Stats(link.ID)
		if err != nil {
			return err
		}

		xlinks = append(xlinks, Link{
			Link:      link,
			ViewStats: stats,
			Number:    i + 1,
		})
	}

	var btns []tele.Btn
	for _, link := range xlinks {
		btns = append(btns, *h.lt.Button(c, "link", link))
	}

	menu := h.b.NewMarkup()
	menu.ResizeKeyboard = true
	menu.Inline(menu.Row(btns...))

	if c.Callback() != nil {
		return c.Edit(
			h.lt.Text(c, "my", xlinks),
			menu, tele.NoPreview)
	} else {
		return c.Send(
			h.lt.Text(c, "my", xlinks),
			menu, tele.NoPreview)
	}
}
