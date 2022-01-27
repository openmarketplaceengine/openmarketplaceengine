package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/config"
	"github.com/openmarketplaceengine/openmarketplaceengine/redis/publisher"
	"github.com/openmarketplaceengine/openmarketplaceengine/redis/subscriber"
	"github.com/openmarketplaceengine/openmarketplaceengine/uids/timestampuid"
	"github.com/stretchr/testify/require"
)

func TestPublisher(t *testing.T) {
	err := config.Read()
	require.NoError(t, err)
	pub := publisher.NewPublisher()
	sub := subscriber.NewSubscriber()

	t.Run("testPublisherAndSubscriber", func(t *testing.T) {
		testPublisherAndSubscriber(t, pub, sub)
	})
}

type message struct {
	Payload string `json:"payload"`
}

func (m *message) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func (m *message) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &m)
}

func testPublisherAndSubscriber(t *testing.T, publisher publisher.Publisher, subscriber subscriber.Subscriber) {

	channel := fmt.Sprintf("chan-%s", timestampuid.NewTimestampUid())
	rcv := make(chan string)
	defer close(rcv)
	ctx := context.Background()
	subscriber.Subscribe(ctx, channel, rcv)

	syncChan := make(chan string)

	go func() {
		defer close(syncChan)
		for {
			select {
			case m := <-rcv:
				syncChan <- m
			}
		}
	}()

	var sent []message
	for i := 0; i < 10; i++ {
		m := message{
			Payload: fmt.Sprintf("payload-%v", i),
		}
		sent = append(sent, m)
		err := publisher.Publish(ctx, channel, &m)
		require.NoError(t, err)
	}

	var received []message

outer:
	for {
		select {
		case m := <-syncChan:
			var mm message
			err := json.Unmarshal([]byte(m), &mm)
			require.NoError(t, err)
			received = append(received, mm)
			if len(received) == len(sent) {
				break outer
			}
		case <-time.After(5 * time.Second):
			require.Fail(t, "expected to receive published message")
		}
	}
	require.Equal(t, sent, received)

}
