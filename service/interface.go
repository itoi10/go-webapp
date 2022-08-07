package service

// storeパッケージの直接参照を避けるためのインターフェース

import (
	"context"

	"github.com/itoi10/go-webapp/entity"
	"github.com/itoi10/go-webapp/store"
)

//go:generate go run github.com/matryer/moq -out moq_test.go . TaskAdder TaskLister
type TaskAdder interface {
	AddTask(ctx context.Context, db store.Execer, t *entity.Task) error
}
type TaskLister interface {
	ListTasks(ctx context.Context, db store.Queryer) (entity.Tasks, error)
}
