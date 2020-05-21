package storage

import (
	"errors"
	"github.com/srgyrn/mdns-todo-api/model"
	"reflect"
	"testing"
)

func tearDown(db Gateway) {
	db = nil
}

func TestDBHandler_AddNewItem(t *testing.T) {
	type fields struct {
		db Gateway
	}
	type args struct {
		content string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.Item
		wantErr bool
	}{
		{
			name: "fails when content is empty",
			fields: fields{
				db: initDbMock(nil),
			},
			args:    args{content: ""},
			want:    model.Item{},
			wantErr: true,
		},
		{
			name: "adds new item successfully",
			fields: fields{
				db: initDbMock(nil),
			},
			args: args{content: "buy some milk"},
			want: model.Item{
				ID:          1,
				Content:     "buy some milk",
				IsCompleted: false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := DBHandler{
				db: tt.fields.db,
			}
			got, err := h.AddNewItem(tt.args.content)
			if (err != nil) != tt.wantErr && !errors.Is(err, ErrNotFound) {
				t.Errorf("AddNewItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddNewItem() got = %v, want %v", got, tt.want)
			}
		})

		tearDown(tt.fields.db)
	}
}

func TestDBHandler_GetItems(t *testing.T) {
	type fields struct {
		db Gateway
	}

	tests := []struct {
		name    string
		fields  fields
		want    []model.Item
		wantErr bool
	}{
		{
			name: "returns not found error when no items are found",
			fields: fields{
				db: initDbMock(nil),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "returns slice of items when there are records",
			fields: fields{
				db: initDbMock([]string{"buy some milk", "get mail"}),
			},
			want: []model.Item{
				{
					ID:          1,
					Content:     "buy some milk",
					IsCompleted: false,
				},
				{
					ID:          2,
					Content:     "get mail",
					IsCompleted: false,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := DBHandler{
				db: tt.fields.db,
			}
			got, err := h.GetItems()
			if (err != nil) != tt.wantErr && !errors.Is(err, ErrNotFound) {
				t.Errorf("GetItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetItems() got = %v, want %v", got, tt.want)
			}
		})

		tearDown(tt.fields.db)
	}
}

func TestDBHandler_DeleteItem(t *testing.T) {
	db := initDbMock([]string{"buy some milk", "get mail"})
	defer tearDown(db)
	handler := DBHandler{
		db: db,
	}

	got, err := handler.DeleteItem(1)
	if err != nil {
		t.Errorf("DeleteItem error = %v", err)
	}

	want := true
	if got != want {
		t.Errorf("DeleteItem() failed. got = %v, want = %v", got, want)
	}
}

// Helper functions below

type dbMock struct {
	items map[int]model.Item
}

func initDbMock(data []string) Gateway {
	var tsx = dbMock{
		make(map[int]model.Item),
	}

	if data != nil {
		for i, d := range data {
			tsx.items[i+1] = model.Item{
				ID:          i + 1,
				Content:     d,
				IsCompleted: false,
			}
		}
	}
	var dbm Gateway
	dbm = tsx

	return dbm
}

func (d dbMock) FindAll() ([]model.Item, error) {
	if len(d.items) == 0 {
		return nil, ErrNotFound
	}

	var itemSlice []model.Item
	for _, v := range d.items {
		itemSlice = append(itemSlice, v)
	}
	return itemSlice, nil
}

func (d dbMock) Find(id int) (model.Item, error) {
	if i, ok := d.items[id]; ok {
		return i, nil
	}

	return model.Item{}, ErrNotFound
}

func (d dbMock) Insert(content string, isCompleted bool) (model.Item, error) {
	id := len(d.items) + 1
	item := model.Item{
		ID:          id,
		Content:     content,
		IsCompleted: isCompleted,
	}

	d.items[id] = item

	return item, nil
}

func (d dbMock) Update(id int, content string, isCompleted bool) (model.Item, error) {
	if _, ok := d.items[id]; ok {
		d.items[id] = model.Item{id, content, isCompleted}
		return d.items[id], nil
	}

	return model.Item{}, errors.New("")
}

func (d dbMock) Delete(id int) error {
	if _, ok := d.items[id]; ok {
		delete(d.items, id)
		return nil
	}

	return errors.New("")
}