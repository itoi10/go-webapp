package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// # ログイン処理
// - POST /register で登録済みのユーザが対象
// - POST /login でユーザ名とパスワードを含んだリクエストを受け取る
// - 認証に成功したユーザにアクセストークンを発行する
//    - アクセストークン有効期限は30分
//    - アクセストークンは改竄防止用に署名がなされ、ログイン情報が含まれる
//    - アクセストークンには次の情報が含まれる
//      - ユーザ名
//      - 権限ロール
// - アプリケーションはアクセストークンからユーザIDを検索できる

type Login struct {
	Service   LoginService
	Validator *validator.Validate
}

func (l *Login) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body struct {
		UserName string `json:"user_name" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	// デコード
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	// 検証
	err := l.Validator.Struct(body)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}
	// アクセストークン発行
	jwt, err := l.Service.Login(ctx, body.UserName, body.Password)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	// 成功レスポンス
	rsp := struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: jwt,
	}
	RespondJSON(r.Context(), w, rsp, http.StatusOK)

}
