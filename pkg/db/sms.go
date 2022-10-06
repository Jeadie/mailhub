package db

import "encoding/xml"

type SmsMessage struct {
    XMLName  xml.Name `json:"Message,omitempty"`
    Smstat   uint     `json:"Smstat,omitempty"`
    Index    uint     `json:"Index,omitempty"`
    Phone    string   `json:"Phone,omitempty"`
    Content  string   `json:"Content,omitempty"`
    Date     string   `json:"Date,omitempty"`
    Sca      any      `json:"Sca,omitempty"`
    SaveType uint     `json:"SaveType,omitempty"`
    Priority uint     `json:"Priority,omitempty,omitempty"`
    SmsType  uint     `json:"SmsType,omitempty"`
}

func (s SmsMessage) isEqualTo(t SmsMessage) bool {
    return s.Phone == t.Phone && s.Content == t.Content && s.Date == t.Date
}
