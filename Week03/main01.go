package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"golang.org/x/sync/errgroup"
)

//1. 基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够 一个退出，全部注销退出

func main() {
	group, ctx := errgroup.WithContext(context.Background())
	group.Go(func() error {
		return Serve(ctx, ":10001")
	})
	group.Go(func() error {
		return Serve(ctx, ":10002")
	})

	group.Go(func() error {
		sigC := make(chan os.Signal, 1)
		signal.Notify(sigC, os.Interrupt)
		select {
		case <-sigC:
			fmt.Println("SIGINT!")
			return errors.New("Stop by SIGINT")

		case <-ctx.Done():
			return nil
		}
	})

	err := group.Wait()
	fmt.Printf("Exit: %+v\n", err)
}

// Serve ...
func Serve(ctx context.Context, addr string) error {
	svr := &http.Server{Addr: addr, Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Listen: %s", addr)
	})}

	errC := make(chan error, 1)
	go func() { errC <- svr.ListenAndServe() }()

	var err error
	select {
	case err = <-errC:
		fmt.Printf("Server %s down! err=%+v\n", addr, err)

	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		svr.Shutdown(shutdownCtx)
		fmt.Printf("Server %s Shutdown!\n", addr)
	}

	return err
}
