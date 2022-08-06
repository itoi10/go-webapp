package store

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/itoi10/go-webapp/clock"
	"github.com/itoi10/go-webapp/entity"
	"github.com/itoi10/go-webapp/testutil"

	"github.com/google/go-cmp/cmp"
	"github.com/jmoiron/sqlx"
)

// RDBMSが必要。ローカルの場合はdocker-composeで起動しておく
func TestRepository_ListTasks(t *testing.T) {
	ctx := context.Background()
	// entity.Taskを作成する他のテストケースと混ざるとテストがフェイルする。
	// そのため、トランザクションをはることでこのテストケースの中だけのテーブル状態にする。
	tx, err := testutil.OpenDBForTest(t).BeginTxx(ctx, nil)
	// このテストケースが完了したら元に戻す
	t.Cleanup(func() { _ = tx.Rollback() })
	if err != nil {
		t.Fatal(err)
	}
	// テーブルを一旦空にしてテストデータを登録する関数
	wants := prepareTasks(ctx, t, tx)

	// テスト対象 タスク取得処理
	sut := &Repository{}
	gots, err := sut.ListTasks(ctx, tx)
	if err != nil {
		t.Fatalf("unexected error: %v", err)
	}
	if d := cmp.Diff(gots, wants); len(d) != 0 {
		t.Errorf("differs: (-got +want)\n%s", d)
	}
}

// 単体テスト等でRDBMSに依存したテストを書きたくない場合は
// DATA-DOG/go-sqlmockパッケージが使える
//
// cf. Goでデータベースを簡単にモック化する【sqlmock】
// https://qiita.com/gold-kou/items/cb174690397f651e2d7f
func TestRepository_AddTask(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c := clock.FixedClocker{}
	var wantID int64 = 20
	okTask := &entity.Task{
		Title:    "ok task",
		Status:   "todo",
		Created:  c.Now(),
		Modified: c.Now(),
	}

	// DATA-DOG/go-sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })

	// 実際にDBにインサートされるわけではなく、SQLドライバの戻り値に注目して正常性を確認する
	mock.ExpectExec(
		// 期待するSQLクエリ。 エスケープが必要
		`INSERT INTO task \(title, status, created, modified\) VALUES \(\?, \?, \?, \?\)`,
	).WithArgs(
		// 期待するクエリパラメータ
		okTask.Title,
		okTask.Status,
		c.Now(),
		c.Now(),
	).WillReturnResult(
		// SQLドライバの戻り値
		sqlmock.NewResult(wantID, 1),
	)

	xdb := sqlx.NewDb(db, "mysql")
	r := &Repository{Clocker: c}
	if err := r.AddTask(ctx, xdb, okTask); err != nil {
		t.Errorf("want no error, but got %v", err)
	}
}

func prepareTasks(ctx context.Context, t *testing.T, con Execer) entity.Tasks {
	t.Helper()
	// テーブルを空にする
	if _, err := con.ExecContext(ctx, "DELETE FROM task;"); err != nil {
		t.Logf("failed to initialize task: %v", err)
	}
	// テスト用の固定時刻取得
	c := clock.FixedClocker{}
	wants := entity.Tasks{
		{
			Title: "want task 1", Status: "todo",
			Created: c.Now(), Modified: c.Now(),
		},
		{
			Title: "want task 2", Status: "todo",
			Created: c.Now(), Modified: c.Now(),
		},
		{
			Title: "want task 3", Status: "done",
			Created: c.Now(), Modified: c.Now(),
		},
	}
	// インサート
	result, err := con.ExecContext(ctx,
		`INSERT INTO task (title, status, created, modified)
			VALUES
			    (?, ?, ?, ?),
			    (?, ?, ?, ?),
			    (?, ?, ?, ?);`,
		wants[0].Title, wants[0].Status, wants[0].Created, wants[0].Modified,
		wants[1].Title, wants[1].Status, wants[1].Created, wants[1].Modified,
		wants[2].Title, wants[2].Status, wants[2].Created, wants[2].Modified,
	)
	if err != nil {
		t.Fatal(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	// ID取得
	wants[0].ID = entity.TaskID(id)
	wants[1].ID = entity.TaskID(id + 1)
	wants[2].ID = entity.TaskID(id + 2)
	return wants
}
