package v1

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"

	"google.golang.org/grpc"
)

type server struct {
	UnimplementedHelloWorldServiceServer
}

func runServer(port int32) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	RegisterHelloWorldServiceServer(grpcServer, &server{})
	return grpcServer.Serve(lis)
}

func (s server) SayHelloSimple(ctx context.Context, request *SayHelloSimpleRequest) (*SayHelloSimpleResponse, error) {
	return &SayHelloSimpleResponse{
		Text: fmt.Sprintf("reply to %s", request.GetText()),
	}, nil
}

func (s server) SayHelloServerStreaming(messageTo *SayHelloServerStreamingRequest, stream HelloWorldService_SayHelloServerStreamingServer) error {
	for i := 0; i < 3; i++ {
		if err := stream.Send(&SayHelloServerStreamingResponse{Text: fmt.Sprintf("%d reply to %s", i, messageTo.GetText())}); err != nil {
			return err
		}
	}
	return nil
}

func (s server) SayHelloClientStreaming(stream HelloWorldService_SayHelloClientStreamingServer) error {
	var messages []*SayHelloClientStreamingRequest
	startTime := time.Now()
	for {
		message, err := stream.Recv()
		if err == io.EOF {
			endTime := time.Now()
			return stream.SendAndClose(&SayHelloClientStreamingResponse{
				StartTime: startTime.Unix(),
				EndTime:   endTime.Unix(),
				Messages:  messages,
			})
		}
		if err != nil {
			return err
		}
		messages = append(messages, message)
	}
}

func (s server) SayHelloBothStreaming(stream HelloWorldService_SayHelloBothStreamingServer) error {
	var messages []*SayHelloBothStreamingRequest
	for {
		message, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		for i := 0; i < 3; i++ {
			if err := stream.Send(&SayHelloBothStreamingResponse{
				Messages: messages,
			}); err != nil {
				return err
			}
		}
		messages = append(messages, message)
	}
}
