package suitetest

import (
	"errors"
	"github.com/srgyrn/mdns-todo-api/model"
	"github.com/srgyrn/mdns-todo-api/storage"
	"reflect"
	"testing"
)

type GatewayTest struct {
	GW storage.Gateway

	Before func()
	After func()
}

func Gateway(t *testing.T, gwt GatewayTest) {
	gwt.Before()

	t.Run("valid insert", validInsert(gwt.GW))

	t.Run("valid find", find(gwt.GW))
	t.Run("invalid find", invalidFind(gwt.GW))
	t.Run("valid findAll", findAll(gwt.GW))

	t.Run("valid update", validUpdate(gwt.GW))
	t.Run("invalid update", invalidUpdate(gwt.GW))

	t.Run("valid delete", validDelete(gwt.GW))

	gwt.After()
}

func validInsert(gw storage.Gateway) func(*testing.T) {
	return func(t *testing.T) {
		want := model.Item{
			ID:          1,
			Content:     "test content",
			IsCompleted: true,
		}

		got, err := gw.Insert(want.Content, true)
		if err != nil {
			t.Errorf("Insert operation error: %v", err)
		}
		if !reflect.DeepEqual(want, got) {
			t.Errorf("Insert operation failed. Want: %v, got: %v", want, got)
		}
	}
}

func find(gw storage.Gateway) func(t *testing.T) {
	return func(t *testing.T) {
		want := model.Item{
			ID:          1,
			Content:     "test content",
			IsCompleted: true,
		}

		got, err := gw.Find(1)
		if err != nil {
			t.Errorf("find error: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("find operation failed. got: %v, want: %v", got, want)
		}
	}
}

func invalidFind(gw storage.Gateway) func(t *testing.T) {
	return func(t *testing.T) {
		got, err := gw.Find(10000)
		if !errors.Is(err,  nil) {
			t.Errorf("find error: %v", err)
		}

		var want model.Item
		if !reflect.DeepEqual(got, want) {
			t.Errorf("find got: %v, want: %v", got, want)
		}
	}
}

func findAll(gw storage.Gateway) func(t *testing.T) {
	return func(t *testing.T) {
		want := []model.Item{
			{
				ID:          1,
				Content:     "test content",
				IsCompleted: true,
			},
		}
		got, err := gw.FindAll()

		if !errors.Is(err, nil) {
			t.Errorf("find all error: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("find all failed. got: %v, want: %v", got, want)
		}
	}
}

func validUpdate(gw storage.Gateway) func(t *testing.T) {
	content := "make a sandwich"
	return func(t *testing.T) {
		want := model.Item{
			ID:          1,
			Content:     content,
			IsCompleted: true,
		}

		got, err := gw.Update(1, content, true)
		if !errors.Is(err, nil) {
			t.Errorf("update error: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("update operation failed. got: %v, want: %v", got, want)
		}
	}
}

func invalidUpdate(gw storage.Gateway) func(t *testing.T) {
	content := "make a sandwich"
	return func(t *testing.T) {
		if _, err := gw.Update(10000, content, true); errors.Is(err, nil) {
			t.Errorf("update operation failed. expected error, got %v", err)
		}
	}
}

func validDelete(gw storage.Gateway) func(t *testing.T) {
	return func(t *testing.T) {
		if err := gw.Delete(1); !errors.Is(err, nil) {
			t.Errorf("delete operation error: %v", err)
		}
	}
}
