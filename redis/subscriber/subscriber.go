package subscriber

import (
	"context"
	"fmt"

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

func (s *subscriber) Subscribe(ctx context.Context, channel string, messages chan<- string) {
	pubSub := s.pubSubClient.Subscribe(ctx, channel)
	log.GetLogger().Info(fmt.Sprintf("[Subscriber] subscribed to %s", channel))
	go func() {
		for {
			select {
			case m := <-pubSub.Channel():
				log.GetLogger().Info(fmt.Sprintf("[Subscriber] received from channel=%s message=%v", channel, m.Payload))
				messages <- m.Payload
			case <-ctx.Done():
				log.GetLogger().Info(fmt.Sprintf("[Subscriber] stopped by context.Done channel=%s", channel))
				return
			}
		}
	}()
}
