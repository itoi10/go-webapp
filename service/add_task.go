package service

import (
	"context"
	"fmt"

	"github.com/itoi10/go-webapp/auth"
	"github.com/itoi10/go-webapp/entity"
	"github.com/itoi10/go-webapp/store"
)

type AddTask struct {
	DB   store.Execer
	Repo TaskAdder
}

// タスク追加
func (a *AddTask) AddTask(ctx context.Context, title string) (*entity.Task, error) {
	// アクセストークンからユーザID取得
	id, ok := auth.GetUserID(ctx)
	if !ok {
		// エラーメッセージ設定 （JSONでエラーレスポンスを返すのはHandlerの役割）
		return nil, fmt.Errorf("useer_id not found")
	}
	// ユーザID, タイトル, ステータス 設定
	t := &entity.Task{
		UserID: id,
		Title:  title,
		Status: entity.TaskStatusTodo,
	}
	// 登録
	err := a.Repo.AddTask(ctx, a.DB, t)
	if err != nil {
		return nil, fmt.Errorf("failed to register: %w", err)
	}
	return t, nil
}
