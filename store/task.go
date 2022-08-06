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
) (entity.Tasks, error) {
	tasks := entity.Tasks{}
	sql := `SELECT id, title, status, created, modified FROM task;`
	if err := db.SelectContext(ctx, &tasks, sql); err != nil {
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
	t.Created = r.Clocker.Now()
	t.Modified = r.Clocker.Now()
	sql := `INSERT INTO task (title, status, created, modified) VALUES (?, ?, ?, ?)`
	result, err := db.ExecContext(
		ctx, sql,
		t.Title,
		t.Status,
		t.Created,
		t.Modified,
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
