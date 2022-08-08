package service

import (
	"context"
	"fmt"

	"github.com/itoi10/go-webapp/entity"
	"github.com/itoi10/go-webapp/store"
)

// handler.ListTasksServiceインターフェースの実装

type ListTask struct {
	DB   store.Queryer
	Repo TaskLister
}

func (l *ListTask) ListTasks(ctx context.Context) (entity.Tasks, error) {
	ts, err := l.Repo.ListTasks(ctx, l.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to list: %w", err)
	}
	return ts, nil
}
