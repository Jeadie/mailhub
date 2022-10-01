package main

type Dao interface {
	Save(to string, s SmsMessage) error
	GetSmssTo(to string) ([]SmsMessage, error)
	GetAllSmss() ([]SmsMessage, error)
}

func CreateDao(inMemory bool) Dao {
	return CreatePersistentDao(inMemory)
}
