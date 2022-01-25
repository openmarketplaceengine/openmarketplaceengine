package publisher

import (
	"context"

	"github.com/openmarketplaceengine/openmarketplaceengine/src/redis/client"

	"github.com/go-redis/redis/v8"
)

type Publisher interface {
	Publish(ctx context.Context, channel string, message interface{}) error
}

type publisher struct {
	pubsubClient *redis.Client
}

func NewPublisher() (pub Publisher) {
	return &publisher{
		pubsubClient: client.NewPubSubClient(),
	}
}

func (p *publisher) Publish(ctx context.Context, channel string, message interface{}) error {
	_, err := p.pubsubClient.Publish(ctx, channel, message).Result()
	return err
}
