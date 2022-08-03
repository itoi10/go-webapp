package main

import "net/http"

// ルーティング

func NewMux() http.Handler {
	mux := http.NewServeMux()
	// ヘルスチェック
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf8")
		// 静的解析エラー回避のため戻り値を捨てている
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})

	return mux
}
