package subscriber

import (
	"context"

	"github.com/openmarketplaceengine/openmarketplaceengine/log"
	"github.com/openmarketplaceengine/openmarketplaceengine/redis/client"

	"github.com/go-redis/redis/v8"
)

type Subscriber interface {
	Subscribe(ctx context.Context, channel string, messages chan<- string)
}

type subscriber struct {
	pubSubClient *redis.Client
}

func NewSubscriber() (sub Subscriber) {
	return &subscriber{
		pubSubClient: client.NewPubSubClient(),
	}
}

func New(client *redis.Client) (sub Subscriber) {
	return &subscriber{
		pubSubClient: client,
	}
}

func (s *subscriber) Subscribe(ctx context.Context, channel string, messages chan<- string) {
	pubSub := s.pubSubClient.Subscribe(ctx, channel)
	log.Infof("[Subscriber] subscribed to %s", channel)
	go func() {
		for {
			select {
			case m := <-pubSub.Channel():
				log.Infof("[Subscriber] received from channel=%s message=%v", channel, m.Payload)
				messages <- m.Payload
			case <-ctx.Done():
				log.Infof("[Subscriber] stopped by context.Done channel=%s", channel)
				return
			}
		}
	}()
}
