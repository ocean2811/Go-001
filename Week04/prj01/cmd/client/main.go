package main

import (
	"context"
	"log"

	pb "go001/api"

	"google.golang.org/grpc"
)

func main() {
	//TODO: config
	conn, err := grpc.Dial(":10001", grpc.WithInsecure())
	if err != nil {
		log.Printf("connect error: %v\n", err)
		return
	}
	defer conn.Close()

	c := pb.NewHelloServiceClient(conn)

	rq := &pb.HelloRequest{Id: "9008000000048942"}
	rs, err := c.Hello(context.Background(), rq)
	if err != nil {
		log.Printf("Hello has error: %+v\n", err)
		return
	}

	log.Printf("%s\n", rs.GetMsg())
}
