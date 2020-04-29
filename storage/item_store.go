package storage

import (
	"errors"
	"github.com/srgyrn/mdns-todo-api/model"
	"github.com/srgyrn/mdns-todo-api/storage/db"
	"log"
	"strings"
)

type (
	Gateway interface {
		FindAll() ([]model.Item, error)
		Find(id int) (model.Item, error)
		Insert(content string, isCompleted bool) (model.Item, error)
		Update(id int, content string, isCompleted bool) (model.Item, error)
		Delete(id int) error
	}

	DBHandler struct {
		db Gateway
	}
)

var ErrNotFound = errors.New("item not found")

func NewDBHandler() DBHandler {
	gw, err := db.NewBoltDB()

	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	return DBHandler{gw}
}

func (h *DBHandler) GetItems() ([]model.Item, error) {
	items, err := h.db.FindAll()

	if !errors.Is(err, nil) {
		return nil, ErrNotFound
	}

	return items, nil
}

func (h DBHandler) AddNewItem(content string) (model.Item, error) {
	var item model.Item

	if len(strings.TrimSpace(content)) == 0 {
		return item, errors.New("content cannot be empty")
	}

	item, err := h.db.Insert(content, false)

	if !errors.Is(err, nil) {
		return model.Item{}, err
	}

	return item, nil
}
