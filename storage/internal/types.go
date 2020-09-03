package internal

import "github.com/srgyrn/mdns-todo-api/model"

type Gateway interface {
	FindAll() ([]model.Item, error)
	Find(id int) (model.Item, error)
	Insert(content string, isCompleted bool) (model.Item, error)
	Update(id int, content string, isCompleted bool) (model.Item, error)
	Delete(id int) error
}
