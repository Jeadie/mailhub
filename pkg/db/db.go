package db

type Dao interface {
	Save(to string, s SmsMessage) error
	GetSmssTo(to string) ([]SmsMessage, error)
	GetAllSmss() ([]SmsMessage, error)
}

// Event represents an SMS sent to the database for storage.
type Event struct {
	to string
	s SmsMessage
}

// CreateDao (optionally just in memory) If newEvents is non-nil, it will stream new, non-duplicate
// events to the channel. Will not ever close channel (even on DB close/crash).
func CreateDao(inMemory bool, newEvents chan Event) Dao {
	return CreatePersistentDao(inMemory, newEvents)
}

