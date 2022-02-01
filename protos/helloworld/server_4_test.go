package helloworld

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"

	"google.golang.org/grpc"
)

type server struct {
	UnimplementedHelloWorldServer
}

func runServer(port int32) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	RegisterHelloWorldServer(grpcServer, &server{})
	return grpcServer.Serve(lis)
}

func (s server) SayHelloSimple(ctx context.Context, messageTo *Message) (*ReplyMessage, error) {
	return &ReplyMessage{
		Text: fmt.Sprintf("reply to %s", messageTo.GetText()),
	}, nil
}

func (s server) SayHelloServerStreaming(messageTo *Message, stream HelloWorld_SayHelloServerStreamingServer) error {
	for i := 0; i < 3; i++ {
		if err := stream.Send(&ReplyMessage{Text: fmt.Sprintf("%d reply to %s", i, messageTo.GetText())}); err != nil {
			return err
		}
	}
	return nil
}

func (s server) SayHelloClientStreaming(stream HelloWorld_SayHelloClientStreamingServer) error {
	var messages []*Message
	startTime := time.Now()
	for {
		message, err := stream.Recv()
		if err == io.EOF {
			endTime := time.Now()
			return stream.SendAndClose(&ReplyMessages{
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

func (s server) SayHelloBothStreaming(stream HelloWorld_SayHelloBothStreamingServer) error {
	for {
		message, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		for i := 0; i < 3; i++ {
			if err := stream.Send(&ReplyMessage{
				Text: fmt.Sprintf("%d reply to %s", i, message.GetText()),
			}); err != nil {
				return err
			}
		}
	}
}
