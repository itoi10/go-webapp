package store

import (
	"context"

	"github.com/itoi10/go-webapp/entity"
)

// 全てのタスクを取得する
// 引数としてQueryerインターフェースを満たす型の値を受け取る
func (r *Repository) ListTasks(
	ctx context.Context,
	db Queryer,
	id entity.UserID,
) (entity.Tasks, error) {
	tasks := entity.Tasks{}
	sql := `
	SELECT
		id,
		user_id,
		title,
		status,
		created,
		modified
	FROM
		task
	WHERE user_id = ?;
	`
	if err := db.SelectContext(ctx, &tasks, sql, id); err != nil {
		return nil, err
	}
	return tasks, nil
}

// タスクを保存する
// 引数としてExecerインターフェースを満たす型の値を受け取る
func (r *Repository) AddTask(
	ctx context.Context,
	db Execer,
	t *entity.Task,
) error {
	sql := `
	INSERT INTO task (
		user_id,
		title,
		status,
		created,
		modified
	) VALUES (
		?, ?, ?, ?, ?
	)`
	result, err := db.ExecContext(
		ctx, sql,
		t.UserID,
		t.Title,
		t.Status,
		r.Clocker.Now(),
		r.Clocker.Now(),
	)
	if err != nil {
		return err
	}
	// 発行されたID取得
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	// *entity.TaskのIDを更新することで呼び出し元に発行された値を伝える
	t.ID = entity.TaskID(id)
	return nil
}
