package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/itoi10/go-webapp/auth"
	"github.com/itoi10/go-webapp/clock"
	"github.com/itoi10/go-webapp/config"
	"github.com/itoi10/go-webapp/handler"
	"github.com/itoi10/go-webapp/service"
	"github.com/itoi10/go-webapp/store"
)

// ルーティング

func NewMux(ctx context.Context, cfg *config.Config) (http.Handler, func(), error) {
	// go-chi/chi/v5
	mux := chi.NewRouter()
	// ヘルスチェック
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		// 静的解析エラー回避のため戻り値を捨てている
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})

	v := validator.New()
	db, cleanup, err := store.New(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}
	clocker := clock.RealClocker{}
	r := store.Repository{Clocker: clocker}
	rcli, err := store.NewKVS(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}
	jwter, err := auth.NewJWTer(rcli, clocker)
	if err != nil {
		return nil, cleanup, err
	}
	l := &handler.Login{
		Service: &service.Login{
			DB:             db,
			Repo:           &r,
			TokenGenerator: jwter,
		},
		Validator: v,
	}
	mux.Post("/login", l.ServeHTTP)

	at := &handler.AddTask{
		Service:   &service.AddTask{DB: db, Repo: &r},
		Validator: v,
	}
	mux.Post("/tasks", at.ServeHTTP)

	lt := &handler.ListTask{
		Service: &service.ListTask{DB: db, Repo: &r},
	}
	mux.Get("/tasks", lt.ServeHTTP)

	ru := &handler.RegisterUser{
		Service:   &service.RegisterUser{DB: db, Repo: &r},
		Validator: v,
	}
	mux.Post("/register", ru.ServeHTTP)

	return mux, cleanup, nil
}
