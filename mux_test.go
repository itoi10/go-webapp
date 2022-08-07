package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/itoi10/go-webapp/config"
)

func TestNewMux(t *testing.T) {
	// httptest.NewRecorderとhttptest.NewRequestを使用すると
	// HTTPサーバを起動せずにHTTPハンドラのテストコードを作成できる
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/health", nil)

	// 環境変数読み込み
	cfg, err := config.New()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	sut, cleanup, err := NewMux(context.Background(), cfg)
	if err != nil {
		t.Fatalf("failed NewMux: %v", err)
	}
	defer cleanup()

	sut.ServeHTTP(w, r)
	resp := w.Result()
	t.Cleanup(func() { _ = resp.Body.Close() })

	if resp.StatusCode != http.StatusOK {
		t.Error("want status code 200, but", resp.StatusCode)
	}
	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	want := `{"status": "ok"}`
	if string(got) != want {
		t.Errorf("want %q, but got %q", want, got)
	}
}
