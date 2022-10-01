package airtable

import "github.com/Jeadie/mailhub/pkg/db"

type Record struct {
	From, To, Message, ReceivedOn string
}

func HandleDbEvent(air *Airtable, e db.Event) error {
	return air.CreateRecord(Record{
		From:       e.S.Phone,
		To:         e.To,
		Message:    e.S.Content,
		ReceivedOn: e.S.Date,
	})
}
