package main

import (
	"backend-engineering/cmd/3.28.grpc/backend.engineering/gen"
	"context"
	"sync"
)

var storage = new(sync.Map)

var ai = new(autoInc)

type autoInc struct {
	sync.Mutex // ensures autoInc is goroutine-safe
	id         int32
}

func (a *autoInc) ID() int32 {
	a.Lock()
	defer a.Unlock()

	a.id++
	return a.id
}

type server struct {
	gen.UnimplementedTodoServer
}

func (s *server) CreateTodo(ctx context.Context, todo *gen.TodoRequest) (*gen.TodoItem, error) {
	item := &gen.TodoItem{
		Id:   ai.ID(),
		Text: todo.Text,
	}

	storage.Store(item.Id, item)

	return item, nil
}
func (s *server) ReadTodos(context.Context, *gen.VoidRequest) (*gen.TodoItems, error) {
	var items = []*gen.TodoItem{}
	storage.Range(func(_, value any) bool {
		items = append(items, value.(*gen.TodoItem))

		return true
	})

	return &gen.TodoItems{Items: items}, nil
}
