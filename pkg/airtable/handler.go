package airtable

import (
	"fmt"
	"github.com/Jeadie/mailhub/pkg/db"
)

type Record struct {
	From, To, Message, ReceivedOn string
}

func HandleDbEvent(air *Airtable, e db.Event) error {
	fmt.Printf("Sending to airtable URL: %s. with data %+v\n", air.GetBaseUrl(), e)
	err := air.CreateRecord(Record{
		From:       e.S.Phone,
		To:         e.To,
		Message:    e.S.Content,
		ReceivedOn: e.S.Date,
	})
	if err != nil {
		fmt.Printf("An error occurred sending %+v to airtable url %s. Error: %s", e, air.GetBaseUrl(), err)
	}
	return err
}
