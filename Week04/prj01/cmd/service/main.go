package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"golang.org/x/sync/errgroup"
)

func main() {
	//TODO: config
	s := InitRPCServer("redis1")

	listen, err := net.Listen("tcp", ":10001")
	if err != nil {
		log.Printf("failed to listen: %+v\n", err)
		return
	}

	group, ctx := errgroup.WithContext(context.Background())
	group.Go(func() error {
		errC := make(chan error, 1)
		go func() { errC <- s.Serve(listen) }()

		var err error
		select {
		case err = <-errC:
			log.Printf("Server down! err=%+v\n", err)

		case <-ctx.Done():
			s.Stop()
			log.Printf("Server Shutdown!\n")
		}

		return err
	})

	group.Go(func() error {
		sigC := make(chan os.Signal, 1)
		signal.Notify(sigC, os.Interrupt)
		select {
		case <-sigC:
			log.Println("SIGINT!")
			return errors.New("Stop by SIGINT")

		case <-ctx.Done():
			return nil
		}
	})

	err = group.Wait()
	fmt.Printf("Exit: %+v\n", err)
}
