package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/itoi10/go-webapp/entity"
	"github.com/itoi10/go-webapp/store"
	"github.com/jmoiron/sqlx"
)

// タスク登録

type AddTask struct {
	DB        *sqlx.DB
	Repo      *store.Repository
	Validator *validator.Validate
}

// http.HandlerFunc型を満たすメソッドを実装
// リクエスト処理が正常完了する場合、RespondJSONを使いJSONレスポンスを返す
// エラー時はErrResponse型に情報を含めてJ、RespondJSONを使いJSONレスポンスを返す
func (at *AddTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Request.BodyをUnmarshalした型
	var b struct {
		//                         バリデーション設定 タイトル必須
		Title string `json:"title" validate:"required"`
	}
	// Unmarshal
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	// 検証
	err := at.Validator.Struct(b)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	// タスク設定
	t := &entity.Task{
		Title:  b.Title,
		Status: entity.TaskStatusTodo,
	}
	// DBに保存
	err = at.Repo.AddTask(ctx, at.DB, t)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	// 正常レスポンス
	rsp := struct {
		ID entity.TaskID `json:"id"`
	}{ID: t.ID}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
