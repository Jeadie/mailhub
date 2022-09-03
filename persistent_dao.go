package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"log"
	"time"
)

type PersistentDao struct {
	db *badger.DB
}

func CreatePersistentDao(inMemory bool) PersistentDao {
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
	}
}

func (d PersistentDao) Save(to string, s SmsMessage) error {
	m := d.db.GetMergeOperator([]byte(to), MergeSms, 200*time.Millisecond)
	defer m.Stop()

	b, err := json.Marshal(s)
	err = m.Add(b)

	return err
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
func MergeSms(originalValue, newValue []byte) []byte {
	var newSms SmsMessage
	err := json.Unmarshal(newValue, &newSms)
	if err != nil {
		fmt.Println(fmt.Errorf("[ERROR]: Could not unmarshall new SMS object, %w\n", err))
		return originalValue
	}

	merge := append(UnmarshalExistingSms(originalValue), newSms)
	bytes, err := json.Marshal(merge)
	if err != nil {
		fmt.Println(fmt.Errorf("[ERROR]: Could not Marshall combined []SmsMessage, %w\n", err))
		return originalValue
	}
	return bytes
}

// UnmarshalExistingSms bytes, which can be either []SmsMessage or SmsMessage.
func UnmarshalExistingSms(v []byte) []SmsMessage {
	var currentSmss []SmsMessage
	err := json.Unmarshal(v, &currentSmss)
	if err == nil {
		return currentSmss
	}

	fmt.Println("Could not unmarshal []SmsMessage. Trying SmsMessage.")

	var currentSms SmsMessage
	err = json.Unmarshal(v, &currentSms)
	if err != nil {
		fmt.Println(fmt.Errorf("[ERROR]: Could not unmarshall existing SMS object/s, %w\n", err))
		return []SmsMessage{}
	}
	return []SmsMessage{currentSms}
}
