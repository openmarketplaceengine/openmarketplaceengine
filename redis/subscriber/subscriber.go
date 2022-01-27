package subscriber

import (
	"context"
	"fmt"

	"github.com/openmarketplaceengine/openmarketplaceengine/redis/client"

	"github.com/go-redis/redis/v8"
)

type Subscriber interface {
	Subscribe(ctx context.Context, channel string, messages chan<- string)
}

type subscriber struct {
	pubsubClient *redis.Client
}

func NewSubscriber() (sub Subscriber) {
	return &subscriber{
		pubsubClient: client.NewPubSubClient(),
	}
}

func (s *subscriber) Subscribe(ctx context.Context, channel string, messages chan<- string) {
	pubSub := s.pubsubClient.Subscribe(ctx, channel)
	logger.Info(fmt.Sprintf("[Subscriber] subscribed to %s", channel))
	go func() {
		for {
			select {
			case m := <-pubSub.Channel():
				logger.Info(fmt.Sprintf("[Subscriber] received from channel=%s message=%v", channel, m.Payload))
				messages <- m.Payload
			case <-ctx.Done():
				logger.Info(fmt.Sprintf("[Subscriber] stopped by context.Done channel=%s", channel))
				return
			}
		}
	}()
}
