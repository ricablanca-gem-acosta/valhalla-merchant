package api

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/brianvoe/gofakeit/v6"
	badger "github.com/dgraph-io/badger/v3"
	"github.com/google/uuid"
)

var db *badger.DB

func InitDb(init bool) error {
	var err error
	needCreateData := false
	if init {
		_, err := os.Stat("db")
		if err != nil {
			needCreateData = true
		}
	}
	opt := badger.DefaultOptions("db")
	db, err = badger.Open(opt)
	if err != nil {
		return err
	}
	if needCreateData {
		return createInitData()
	}
	return nil
}

func CloseDb(dropAll bool) error {
	if dropAll {
		err := db.DropAll()
		if err != nil {
			return err
		}
	}
	return db.Close()
}

func createInitData() error {
	for i := 0; i < 5; i++ {
		m, err := createMerchant(nil)
		if err != nil {
			return err
		}
		for i := 0; i < 200; i++ {
			email := gofakeit.Email()
			err = createMember(m.Code, email)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func createMerchant(code *string) (*Merchant, error) {
	if code == nil {
		generatedCode := uuid.New().String()
		code = &generatedCode
	}
	var value []byte
	merchant := &Merchant{*code, []Member{}}
	value, err := json.Marshal(merchant)
	if err != nil {
		return nil, err
	}
	txn := db.NewTransaction(true)
	err = txn.SetEntry(badger.NewEntry([]byte(*code), value))
	if err != nil {
		return nil, err
	}
	err = txn.Commit()
	if err != nil {
		return nil, err
	}
	return merchant, nil
}

func getMerchant(code string) (*Merchant, error) {
	var err error
	var merchant Merchant
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(code))
		if err != nil {
			return err
		}
		val, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		err = json.Unmarshal(val, &merchant)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &merchant, nil
}

func listMerchants() ([]Merchant, error) {
	var ret []Merchant
	err := db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			var merchant Merchant
			val, err := it.Item().ValueCopy(nil)
			if err != nil {
				return err
			}
			err = json.Unmarshal(val, &merchant)
			if err != nil {
				return err
			}
			ret = append(ret, merchant)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func deleteMerchant(code string) error {
	txn := db.NewTransaction(true)
	err := txn.Delete([]byte(code))
	if err != nil {
		return err
	}
	return txn.Commit()
}

func createMember(merchantCode string, email string) error {
	merchant, err := getMerchant(merchantCode)
	if err != nil {
		return err
	}
	merchant.Members = append(merchant.Members, Member{email})
	value, err := json.Marshal(merchant)
	if err != nil {
		return err
	}
	txn := db.NewTransaction(true)
	err = txn.Set([]byte(merchantCode), value)
	if err != nil {
		return err
	}
	return txn.Commit()
}

func deleteMember(merchantCode string, email string) error {
	merchant, err := getMerchant(merchantCode)
	if err != nil {
		return err
	}
	var remIdx = -1
	for i, m := range merchant.Members {
		if m.Email == email {
			remIdx = i
			break
		}
	}
	if remIdx == -1 {
		return errors.New("Member not found")
	}
	var l = len(merchant.Members)
	merchant.Members[remIdx] = merchant.Members[l-1]
	merchant.Members = merchant.Members[:l-1]
	value, err := json.Marshal(merchant)
	if err != nil {
		return err
	}
	txn := db.NewTransaction(true)
	err = txn.Set([]byte(merchantCode), value)
	if err != nil {
		return err
	}
	return txn.Commit()
}
