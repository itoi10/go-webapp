package service

import (
	"context"
	"fmt"

	"github.com/itoi10/go-webapp/auth"
	"github.com/itoi10/go-webapp/entity"
	"github.com/itoi10/go-webapp/store"
)

// handler.ListTasksServiceインターフェースの実装

type ListTask struct {
	DB   store.Queryer
	Repo TaskLister
}

func (l *ListTask) ListTasks(ctx context.Context) (entity.Tasks, error) {
	// アクセストークンからユーザID取得
	id, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, fmt.Errorf("user_id not found")
	}
	ts, err := l.Repo.ListTasks(ctx, l.DB, id)
	if err != nil {
		return nil, fmt.Errorf("failed to list: %w", err)
	}
	return ts, nil
}
