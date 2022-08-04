package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/itoi10/go-webapp/entity"
	"github.com/itoi10/go-webapp/store"
	"github.com/itoi10/go-webapp/testutil"
)

func TestAddTask(t *testing.T) {
	type want struct {
		status  int
		rspFile string
	}
	// - ゴールデンテスト
	//   テストの入力値や期待値を別ファイルとして保存したデータを使用する
	tests := map[string]struct {
		reqFile string
		want    want
	}{
		// 正常系
		"ok": {
			reqFile: "testdata/add_task/ok_req.json.golden",
			want: want{
				status:  http.StatusOK,
				rspFile: "testdata/add_task/ok_rsp.json.golden",
			},
		},
		// 異常系
		"badRequest": {
			reqFile: "testdata/add_task/bad_req.json.golden",
			want: want{
				status:  http.StatusBadRequest,
				rspFile: "testdata/add_task/bad_req_rsp.json.golden",
			},
		},
	}
	// - テーブルドリブンテスト
	//   複数の入力や期待値の組み合わせを共通化した実行手順で実行させるテストの実装パターン
	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			// テスト並列化
			t.Parallel()

			// 擬似リクエスト
			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodPost,
				"/tasks",
				bytes.NewReader(testutil.LoadFile(t, tt.reqFile)),
			)

			// ハンドラ
			sut := AddTask{
				Store: &store.TaskStore{
					Tasks: map[entity.TaskID]*entity.Task{},
				},
				Validator: validator.New(),
			}
			sut.ServeHTTP(w, r)

			// レスポンス検証
			resp := w.Result()
			testutil.AssertResponse(t,
				resp, tt.want.status, testutil.LoadFile(t, tt.want.rspFile),
			)
		})
	}
}
