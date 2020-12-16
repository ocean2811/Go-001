package service

import (
	"context"
	pb "go001/api"
	"go001/internal/biz"
	"log"
	"net"

	"google.golang.org/grpc"
)

type RPCServer struct {
	bizHandler *biz.BizHandler
	grpcServer *grpc.Server
}

func NewRPCServer(bizHandler *biz.BizHandler) Server {
	r := &RPCServer{bizHandler: bizHandler}

	r.grpcServer = grpc.NewServer()

	pb.RegisterHelloServiceServer(r.grpcServer, r)

	return r
}

func (rpc *RPCServer) Hello(ctx context.Context, rq *pb.HelloRequest) (*pb.HelloResponse, error) {
	if rpc.bizHandler == nil {
		return nil, ErrServerUnavailable
	}

	msg, err := rpc.bizHandler.GenHelloMsg(rq.GetId())
	if err != nil {
		log.Printf("RPCServer Hello has error=%+v,rq=%s", err, rq)
		return nil, ErrUserInvalid
	}

	return &pb.HelloResponse{Msg: msg}, nil
}

func (rpc *RPCServer) Serve(l net.Listener) error {
	if rpc.grpcServer == nil {
		return ErrServerUnavailable
	}

	return rpc.grpcServer.Serve(l)
}

func (rpc *RPCServer) Stop() {
	if rpc.grpcServer == nil {
		return
	}

	rpc.grpcServer.Stop()
	return
}
