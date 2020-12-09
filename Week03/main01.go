package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"golang.org/x/sync/errgroup"
)

//1. 基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够 一个退出，全部注销退出

func main() {
	signalCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group, ctx := errgroup.WithContext(signalCtx)
	group.Go(func() error {
		return Serve(ctx, ":10001")
	})
	group.Go(func() error {
		return Serve(ctx, ":10002")
	})

	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, os.Interrupt)
	go func() {
		<-sigC
		fmt.Println("SIGINT!")
		cancel()
	}()

	err := group.Wait()
	fmt.Printf("Exit: %+v\n", err)
}

// Serve ...
func Serve(ctx context.Context, addr string) error {
	svr := &http.Server{Addr: addr, Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Listen: %s", addr)
	})}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		svr.Shutdown(shutdownCtx)
	}()

	return svr.ListenAndServe()
}
