package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"

	"golang.org/x/sync/errgroup"
)

func main() {
	group, ctx := errgroup.WithContext(context.Background())
	group.Go(func() error {
		listener, err := net.Listen("tcp", ":12345")
		if err != nil {
			return err
		}

		go func() {
			<-ctx.Done()
			listener.Close()
		}()

		for {
			conn, err := listener.Accept()
			if err != nil {
				select {
				case <-ctx.Done():
					fmt.Println("Listener Close!")
					return nil
				default:
				}

				fmt.Printf("Server accept err=%+v\n", err)
				continue
			}

			Handle(conn)
		}
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

// Handle conn
func Handle(conn net.Conn) {
	c := make(chan []byte)

	go ServerWriter(conn, c)
	go ServerReader(conn, c)
}

// ServerReader conn
func ServerReader(r io.ReadCloser, c chan<- []byte) {
	defer r.Close()

	scan := bufio.NewScanner(r)

	for {
		if !scan.Scan() {
			if scan.Err() != nil {
				//not EOF
				fmt.Printf("ServerReader error: %+v\n", scan.Err())
			}
			close(c)
			return
		}

		if b := scan.Bytes(); len(b) > 0 {
			c <- append(append([]byte("echo:"), scan.Bytes()...), '\n')
		}
	}

}

// ServerWriter conn
func ServerWriter(w io.WriteCloser, c <-chan []byte) {
	defer w.Close()

	for b := range c {
		_, err := w.Write(b)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("ServerReader error: %+v\n", err)
			}
			return
		}
	}

	return
}
