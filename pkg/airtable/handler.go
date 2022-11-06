package airtable

import (
	"fmt"
	"github.com/Jeadie/mailhub/pkg/db"
	"time"
)

type Record struct {
	From, To, Message, ReceivedOn, SmsTimestamp string
}

func recordFromDbEvent(e db.Event) Record {
	return Record{
		From:         e.S.Phone,
		To:           e.To,
		Message:      e.S.Content,
		ReceivedOn:   time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
		SmsTimestamp: e.S.Date,
	}
}

func HandleDbEvent(air *Airtable, e db.Event) error {
	fmt.Printf("Sending to airtable URL: %s. with data %+v\n", air.GetBaseUrl(), e)
	err := air.CreateRecord(recordFromDbEvent(e))
	if err != nil {
		fmt.Printf("An error occurred sending %+v to airtable url %s. Error: %s", e, air.GetBaseUrl(), err)
	}
	return err
}
