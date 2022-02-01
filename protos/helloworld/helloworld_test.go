package helloworld

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const port = 10123

func TestHelloWorld(t *testing.T) {

	go func() {
		err := runServer(port)
		if err != nil {
			require.NoError(t, err)
		}
	}()

	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))

	require.NoError(t, err)

	defer func(conn *grpc.ClientConn) {
		innerErr := conn.Close()
		if innerErr != nil {
			require.NoError(t, innerErr)
		}
	}(conn)

	client := NewHelloWorldClient(conn)

	t.Run("testSayHelloSimple", func(t *testing.T) {
		testSayHelloSimple(t, client)
	})
	t.Run("testSayHelloClientStreaming", func(t *testing.T) {
		testSayHelloClientStreaming(t, client)
	})
}

func testSayHelloSimple(t *testing.T, client HelloWorldClient) {
	to := &Message{Text: "test 1"}
	messageFrom, err := client.SayHelloSimple(context.Background(), to)
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("reply to %s", to.Text), messageFrom.GetText())
}

func testSayHelloClientStreaming(t *testing.T, client HelloWorldClient) {
	streaming, err := client.SayHelloClientStreaming(context.Background())
	require.NoError(t, err)

	for i := 0; i < 3; i++ {
		err = streaming.Send(&Message{Text: fmt.Sprintf("test-%v", i)})
		require.NoError(t, err)
	}

	messages, err := streaming.CloseAndRecv()
	require.NoError(t, err)
	require.GreaterOrEqual(t, messages.GetEndTime(), messages.GetStartTime())
	require.Len(t, messages.GetMessages(), 3)
}
