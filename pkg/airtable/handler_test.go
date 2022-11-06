package airtable

import (
	"encoding/xml"
	"github.com/Jeadie/mailhub/pkg/db"
	"testing"
	"time"
)

func TestRecordFromDbEvent(t *testing.T) {
	r := recordFromDbEvent(db.Event{
		To: "Someone",
		S: db.SmsMessage{
			XMLName:  xml.Name{},
			Smstat:   0,
			Index:    0,
			Phone:    "61412345678",
			Content:  "Hello World",
			Date:     "1970-01-01 12:00:01",
			Sca:      nil,
			SaveType: 0,
			Priority: 0,
			SmsType:  0,
		},
	})
	if r.To != "Someone" {
		t.Errorf("Incorrect `To` field")
	}
	if r.SmsTimestamp != "1970-01-01 12:00:01" {
		t.Errorf("Incorrect `SmsTimestamp` field %s", r.SmsTimestamp)
	}

	receivedOn, err := time.Parse("2006-01-02T15:04:05.000Z", r.ReceivedOn)
	if err != nil {
		t.Errorf("Invalid `ReceivedOn` date format")
	}

	if time.Now().Sub(receivedOn).Seconds() > 10 {
		t.Errorf("Incorrect `RecievedOn` field %s. Should be closer to present time", r.ReceivedOn)
	}
}
