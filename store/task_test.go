package store

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/itoi10/go-webapp/clock"
	"github.com/itoi10/go-webapp/entity"
	"github.com/itoi10/go-webapp/testutil"
	"github.com/itoi10/go-webapp/testutil/fixture"

	"github.com/google/go-cmp/cmp"
	"github.com/jmoiron/sqlx"
)

// RDBMSが必要。ローカルの場合はdocker-composeで起動しておく
func TestRepository_ListTasks(t *testing.T) {
	t.Parallel()

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
	wantUserID, wants := prepareTasks(ctx, t, tx)

	// テスト対象 タスク取得処理
	sut := &Repository{}
	gots, err := sut.ListTasks(ctx, tx, wantUserID)
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
		UserID:   33,
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
		`INSERT INTO task \( user_id, title, status, created, modified \) VALUES \( \?, \?, \?, \?, \? \)`,
	).WithArgs(
		// 期待するクエリパラメータ
		okTask.UserID,
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

// ユーザデータ準備
func prepareUser(ctx context.Context, t *testing.T, db Execer) entity.UserID {
	t.Helper()
	u := fixture.User(nil)
	result, err := db.ExecContext(ctx, `
		INSERT INTO user (
			name,
			password,
			role,
			created,
			modified
		) VALUES
			(?, ?, ?, ?, ?);`,
		u.Name, u.Password, u.Role, u.Created, u.Modified)
	if err != nil {
		t.Fatalf("insert user: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("got user_id: %v", err)
	}
	return entity.UserID(id)
}

// タスクデータ準備
func prepareTasks(ctx context.Context, t *testing.T, con Execer) (entity.UserID, entity.Tasks) {
	t.Helper()
	// ユーザデータ作成
	userID := prepareUser(ctx, t, con)
	otherUserID := prepareUser(ctx, t, con)
	// テスト用の固定時刻取得
	c := clock.FixedClocker{}
	wants := entity.Tasks{
		{
			UserID: userID,
			Title:  "want task 1", Status: "todo",
			Created: c.Now(), Modified: c.Now(),
		},
		{
			UserID: userID,
			Title:  "want task 2", Status: "done",
			Created: c.Now(), Modified: c.Now(),
		},
	}
	tasks := entity.Tasks{
		wants[0],
		{
			UserID: otherUserID,
			Title:  "not want task", Status: "todo",
			Created: c.Now(), Modified: c.Now(),
		},
		wants[1],
	}
	// インサート
	result, err := con.ExecContext(ctx,
		`INSERT INTO task (
			user_id,
			title,
			status,
			created,
			modified
		) VALUES
			(?, ?, ?, ?, ?),
			(?, ?, ?, ?, ?),
			(?, ?, ?, ?, ?);`,
		tasks[0].UserID, tasks[0].Title, tasks[0].Status, tasks[0].Created, tasks[0].Modified,
		tasks[1].UserID, tasks[1].Title, tasks[1].Status, tasks[1].Created, tasks[1].Modified,
		tasks[2].UserID, tasks[2].Title, tasks[2].Status, tasks[2].Created, tasks[2].Modified,
	)
	if err != nil {
		t.Fatal(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	// ID取得
	tasks[0].ID = entity.TaskID(id)
	tasks[1].ID = entity.TaskID(id + 1)
	tasks[2].ID = entity.TaskID(id + 2)
	return userID, wants
}
