package db

import (
	"encoding/json"
	"fmt"
	"github.com/srgyrn/mdns-todo-api/model"
	"go.etcd.io/bbolt"
	"log"
	"strconv"
)

type BucketCreationError struct {
	bucketName string
}

type BoltDB struct {
	dbName         string
	rootBucketName string
	bucketName     string
	db             *bbolt.DB
}

func (e *BucketCreationError) Error() string {
	return fmt.Sprintf("failed to create bucket %s", e.bucketName)
}

func NewBoltDB() (*BoltDB, error) {
	bolt := &BoltDB{
		dbName:         "todolist.db",
		rootBucketName: "DB",
		bucketName:     "ITEMS",
		db:             nil,
	}

	err := initBoltDB(bolt)
	return bolt, err
}

func initBoltDB(bolt *BoltDB) error {
	conn, err := bbolt.Open(bolt.dbName, 0600, nil)

	if err != nil {
		log.Panic(err)
	}

	err = conn.Update(func(tx *bbolt.Tx) error {
		root, err := tx.CreateBucketIfNotExists([]byte(bolt.rootBucketName))

		if err != nil {
			return &BucketCreationError{bucketName: bolt.rootBucketName}
		}

		_, err = root.CreateBucketIfNotExists([]byte(bolt.bucketName))

		if err != nil {
			return &BucketCreationError{bucketName: bolt.bucketName}
		}

		return nil
	})

	if err != nil {
		log.Fatal("failed to create DB")
	}

	bolt.db = conn
	return nil
}

func (b *BoltDB) FindAll() ([]model.Item, error) {
	var item model.Item
	var items []model.Item
	err := b.db.View(func(tx *bbolt.Tx) error {
		tx.Bucket([]byte(b.rootBucketName)).Bucket([]byte(b.bucketName)).ForEach(func(_, v []byte) error {
			json.Unmarshal(v, &item)
			items = append(items, item)

			return nil
		})
		return nil
	})

	return items, err
}

func (b *BoltDB) Find(id int) (model.Item, error) {
	var item model.Item
	err := b.db.View(func(tx *bbolt.Tx) error {
		result := tx.Bucket([]byte(b.rootBucketName)).Bucket([]byte(b.bucketName)).Get([]byte(strconv.FormatUint(uint64(id), 10)))
		json.Unmarshal(result, &item)
		return nil
	})

	return item, err
}

func (b *BoltDB) Insert(content string, isCompleted bool) (model.Item, error) {
	item := model.Item{
		ID:          0,
		Content:     content,
		IsCompleted: isCompleted,
	}

	err := b.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(b.rootBucketName)).Bucket([]byte(b.bucketName))
		ns, err := bucket.NextSequence()
		if err != nil {
			return fmt.Errorf("failed to get next sequence: %v", err)
		}

		item.ID = int(ns)
		serializedItem, err := json.Marshal(item)
		if err != nil {
			return fmt.Errorf("failed to serialize item: %v", item)
		}

		err = bucket.Put([]byte(strconv.FormatUint(ns, 10)), serializedItem)

		if err != nil {
			return fmt.Errorf("failed to insert new item: %v", item)
		}

		return nil
	})

	return item, err
}

func (b *BoltDB) Update(id int, content string, isCompleted bool) (model.Item, error) {
	var item model.Item

	err := b.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(b.rootBucketName)).Bucket([]byte(b.bucketName))
		result := bucket.Get([]byte(strconv.Itoa(id)))

		if result == nil {
			return fmt.Errorf("item not found")
		}

		err := json.Unmarshal(result, &item)
		if err != nil {
			return fmt.Errorf("failed to serialize item: %v", item)
		}

		item.Content = content
		item.IsCompleted = isCompleted

		serializedItem, err := json.Marshal(item)

		if err != nil {
			return fmt.Errorf("failed to serialize item: %v", err)
		}

		err = bucket.Put([]byte(strconv.Itoa(id)), serializedItem)

		if err != nil {
			return fmt.Errorf("failed to insert new item: %v", item)
		}

		return nil
	})

	return item, err
}

func (b *BoltDB) Delete(id int) error {
	err := b.db.Update(func(tx *bbolt.Tx) error {
		err := tx.Bucket([]byte(b.rootBucketName)).Bucket([]byte(b.bucketName)).Delete([]byte(strconv.FormatUint(uint64(id), 10)))

		if err != nil {
			return fmt.Errorf("failed to delete item: %v", err)
		}

		return nil
	})

	return err
}
