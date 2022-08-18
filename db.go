package main

type Dao struct {
	db map[string][]SmsMessage
}

func (d *Dao) Save(to string, s SmsMessage) error {
	d.db[to] = append(d.db[to], s)
	return nil
}

func (d *Dao) GetSmssTo(to string) ([]SmsMessage, error) {
	return d.db[to], nil
}

func (d *Dao) GetAllSmss() ([]SmsMessage, error) {
	result := make([]SmsMessage, 1)
	for _, val := range d.db {
		result = append(result, val...)
	}
	return result, nil
}
