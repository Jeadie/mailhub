package airtable

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CreateRecordJson struct {
	Records []struct {
		Fields map[string]string `json:"fields"`
	} `json:"records"`
}

// Airtable provides methods to interact with a specific Airtable Table.
type Airtable struct {
	apiKey, baseId, tableName string
}

func CreateAirtable(baseId, tableName, apiKey string) *Airtable {
	return &Airtable{
		baseId:    baseId,
		apiKey:    apiKey,
		tableName: tableName,
	}
}

func (a Airtable) CreateRecord(r Record) error {
	return a.create(map[string]string{
		"To":         r.To,
		"From":       r.From,
		"Message":    r.Message,
		"ReceivedOn": r.ReceivedOn,
	})
}

func (a Airtable) create(fields map[string]string) error {
	body, err := a.createRecordBody(fields)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", a.getBaseUrl(), body)
	a.addHeaders(req)
	if err != nil {
		return err
	}

	_, err = (&http.Client{}).Do(req)
	return err
}

func (a Airtable) addHeaders(r *http.Request) {
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.apiKey))
	r.Header.Add("Content-Type", "application/json")
}

func (a Airtable) getBaseUrl() string {
	return fmt.Sprintf("https://api.airtable.com/v0/%s/%s", a.baseId, a.tableName)
}

func (a Airtable) createRecordBody(fields map[string]string) (io.Reader, error) {
	createJson := CreateRecordJson{Records: []struct {
		Fields map[string]string `json:"fields"`
	}{{Fields: fields}}}

	createBytes, err := json.Marshal(createJson)
	return bytes.NewReader(createBytes), err
}
