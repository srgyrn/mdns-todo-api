package storage

import (
	"errors"
	"log"
	"sort"
	"strings"

	"github.com/srgyrn/mdns-todo-api/model"
	"github.com/srgyrn/mdns-todo-api/storage/db"
	"github.com/srgyrn/mdns-todo-api/storage/internal"
)
type (
	DBHandler struct {
		db internal.Gateway
	}

	ById []model.Item
)

var ErrNotFound = errors.New("item not found")

func (a ById) Len() int           { return len(a) }
func (a ById) Less(i, j int) bool { return a[i].ID < a[j].ID }
func (a ById) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

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

	sort.Sort(ById(items))
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

func (h DBHandler) DeleteItem(id int) (bool, error) {
	err := h.db.Delete(id)
	return err == nil, err
}
