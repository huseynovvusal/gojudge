package main

import (
	"context"
	"fmt"
	"huseynovvusal/gojudge/internal/executor"
	pb "huseynovvusal/gojudge/internal/proto"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedExecutorServiceServer
}

func (*server) Execute(ctx context.Context, req *pb.ExecutorRequest) (*pb.ExecutorResponse, error) {
	result, err := executor.RunCode(req.Language, req.Code, req.Input, int16(req.TimeLimit), int16(req.MemoryLimit), int16(req.CpuLimit))

	if err != nil {
		return &pb.ExecutorResponse{}, err
	}

	return &pb.ExecutorResponse{
		Output:        result.Output,
		ExecutionTime: result.ExecutionMs,
	}, nil

}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Println(err)
		return
	}

	s := grpc.NewServer()

	pb.RegisterExecutorServiceServer(s, &server{})

	fmt.Println("Server is running on port :50051")

	if err := s.Serve(lis); err != nil {
		fmt.Println(err)
	}

}
