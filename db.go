package main

type Dao interface {
	Save(to string, s SmsMessage) error
	GetSmssTo(to string) ([]SmsMessage, error)
	GetAllSmss() ([]SmsMessage, error)
}

func CreateDao() Dao {
	return InMemoryDao{make(map[string][]SmsMessage)}
}

type InMemoryDao struct {
	db map[string][]SmsMessage
}

func (d InMemoryDao) Save(to string, s SmsMessage) error {
	d.db[to] = append(d.db[to], s)
	return nil
}

func (d InMemoryDao) GetSmssTo(to string) ([]SmsMessage, error) {
	return d.db[to], nil
}

func (d InMemoryDao) GetAllSmss() ([]SmsMessage, error) {
	result := make([]SmsMessage, 1)
	for _, val := range d.db {
		result = append(result, val...)
	}
	return result, nil
}
