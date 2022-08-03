package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

// http.Server型をラップした型を作りHTTPサーバに関わる部分を分割する
// ルーティングは引数で受け取りServer型の責務から除外する

type Server struct {
	srv *http.Server
	l   net.Listener
}

func NewServer(l net.Listener, mux http.Handler) *Server {
	return &Server{
		srv: &http.Server{Handler: mux},
		l:   l,
	}
}

func (s *Server) Run(ctx context.Context) error {
	// シグナルハンドリング
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	eg, ctx := errgroup.WithContext(ctx)
	// 別のゴルーチンでHTTPサーバを起動
	eg.Go(func() error {
		if err := s.srv.Serve(s.l); err != nil &&
			// http.ErrServerClosedはhttp.Serve.Shutdown()が正常終了したことを示すので異常ではない。
			err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})

	// チャネルからの終了通知待機
	<-ctx.Done()
	if err := s.srv.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdonw: %+v", err)
	}

	// グレースフルシャットダウンの終了を待つ
	return eg.Wait()
}
