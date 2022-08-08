package service

import (
	"context"
	"fmt"

	"github.com/itoi10/go-webapp/entity"
	"github.com/itoi10/go-webapp/store"
)

type AddTask struct {
	DB   store.Execer
	Repo TaskAdder
}

// 実装
func (a *AddTask) AddTask(ctx context.Context, title string) (*entity.Task, error) {
	t := &entity.Task{
		Title:  title,
		Status: entity.TaskStatusTodo,
	}
	err := a.Repo.AddTask(ctx, a.DB, t)
	if err != nil {
		return nil, fmt.Errorf("failed to register: %w", err)
	}
	return t, nil
}
