package service

import (
	"context"
	"fmt"

	"github.com/itoi10/go-webapp/store"
)

// # ログイン情報の検証とアクセストークンの生成を行う

// テストを考慮し *store.Repository, *auth.JWTer を直接参照せず、インターフェースを使用する
type Login struct {
	DB             store.Queryer
	Repo           UserGetter
	TokenGenerator TokenGenerator
}

func (l *Login) Login(ctx context.Context, name, pw string) (string, error) {
	u, err := l.Repo.GetUser(ctx, l.DB, name)
	if err != nil {
		return "", fmt.Errorf("failed to list: %w", err)
	}
	if err = u.ComparePassword(pw); err != nil {
		return "", fmt.Errorf("wrong password: %w", err)
	}
	jwt, err := l.TokenGenerator.GenerateToken(ctx, *u)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	return string(jwt), nil
}
