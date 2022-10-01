package db

import (
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"log"
	"time"
)

type PersistentDao struct {
	db *badger.DB
	stream chan Event
}

func CreatePersistentDao(inMemory bool, events chan Event) PersistentDao {
	var opt badger.Options
	if inMemory {
		opt = badger.DefaultOptions("").WithInMemory(true)
	} else {
		opt = badger.DefaultOptions(".badger")
	}

	db, err := badger.Open(opt)
	if err != nil {
		log.Fatal(err)
	}
	return PersistentDao{
		db: db,
		stream: events,
	}
}

func (d PersistentDao) Save(to string, s SmsMessage) error {
	if !d.keyExists(to) {
		// New keys won't get sent to stream in the Merge Operator. Must be done manually
		d.stream <- Event{
			to: to,
			s:  s,
		}
	}

	m := d.db.GetMergeOperator([]byte(to),  func(existingVal, newVal []byte) []byte {
		return d.MergeSms(to, existingVal, newVal)
	}, 200*time.Millisecond)

	defer m.Stop()

	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return m.Add(b)

}

func (d PersistentDao) keyExists(to string) bool {
	exists := false

	d.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(to))
		if err == badger.ErrKeyNotFound {
			exists = false
		} else {
			exists = true
		}
		return err
	})
	return exists
}

func (d PersistentDao) GetSmssTo(to string) ([]SmsMessage, error) {
	var result []SmsMessage
	err := d.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(to))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			var sms []SmsMessage
			err := json.Unmarshal(val, &sms)

			result = append(result, sms...)
			return err
		})
	})
	return result, err
}

func (d PersistentDao) GetAllSmss() ([]SmsMessage, error) {
	var result []SmsMessage
	err := d.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			err := it.Item().Value(func(v []byte) error {
				var itSms []SmsMessage
				err := json.Unmarshal(v, &itSms)

				result = append(result, itSms...)
				return err
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return result, err
}

// MergeSms for exising keys (phone number/user) with a new, single SmsMessage.
func (d PersistentDao) MergeSms(to string, originalValue, newValue []byte) []byte {

	// Decode value to SmsMessage/s
	var newSms SmsMessage
	err := json.Unmarshal(newValue, &newSms)
	if err != nil {
		fmt.Println(fmt.Errorf("[ERROR]: Could not unmarshall new SMS object, %w\n", err))
		return originalValue
	}

	// Add newValue, if not duplicate
	existing := UnmarshalExistingSms(originalValue)
	if !SmsInList(newSms, existing) {
		existing = append(existing, newSms)
		d.stream <- Event{
			to: to,
			s:  newSms,
		}
	}


	// Encode result back to byte[]
	bytes, err := json.Marshal(existing)
	if err != nil {
		fmt.Println(fmt.Errorf("[ERROR]: Could not Marshall combined []SmsMessage, %w\n", err))
		return originalValue
	}
	return bytes
}

// SmsInList returns true iff s is in the list of SMSs (as defined by SmsMessage's isEqualTo).
func SmsInList(s SmsMessage, sList []SmsMessage) bool {
	for _, ss := range sList {
		if s.isEqualTo(ss) {
			return true
		}
	}
	return false
}

// UnmarshalExistingSms bytes, which can be either []SmsMessage or SmsMessage.
func UnmarshalExistingSms(v []byte) []SmsMessage {
	var currentSmss []SmsMessage
	err := json.Unmarshal(v, &currentSmss)
	if err == nil {
		return currentSmss
	}

	fmt.Println("Could not unmarshal []SmsMessage. Trying a single SmsMessage.")

	var currentSms SmsMessage
	err = json.Unmarshal(v, &currentSms)
	if err != nil {
		fmt.Println(fmt.Errorf("[ERROR]: Could not unmarshall existing SMS object/s, %w\n", err))
		return []SmsMessage{}
	}
	return []SmsMessage{currentSms}
}
