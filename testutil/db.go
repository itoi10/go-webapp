package testutil

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// 実行中の環境によりDB接続先を切り替える
func OpenDBForTest(t *testing.T) *sqlx.DB {
	// ローカル(docker-compose.ymlでマッピング)のポート番号
	port := 33306
	// GitHub Actions環境用のポート番号 (環境変数"CI"はGitHub Actions上でのみ定義する)
	if _, defined := os.LookupEnv("CI"); defined {
		port = 3306
	}

	db, err := sql.Open(
		"mysql",
		fmt.Sprintf("todo:todo@tcp(127.0.0.1:%d)/todo?parseTime=true", port),
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(
		func() { _ = db.Close() },
	)
	return sqlx.NewDb(db, "mysql")
}
