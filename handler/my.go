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

	const maxInRow = 5
	rows := make([]tele.Row, (len(links)-1)/maxInRow+1)
	for i, link := range xlinks {
		i /= maxInRow
		rows[i] = append(rows[i], *h.lt.Button(c, "link", link))
	}

	menu := h.b.NewMarkup()
	menu.ResizeKeyboard = true
	menu.Inline(rows...)

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
