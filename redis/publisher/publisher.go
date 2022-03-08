package publisher

import (
	"context"

	"github.com/openmarketplaceengine/openmarketplaceengine/redis/client"

	"github.com/go-redis/redis/v8"
)

type Publisher interface {
	Publish(ctx context.Context, channel string, message interface{}) error
}

type publisher struct {
	pubSubClient *redis.Client
}

func NewPublisher() (pub Publisher) {
	return &publisher{
		pubSubClient: client.NewPubSubClient(),
	}
}

func New(client *redis.Client) (pub Publisher) {
	return &publisher{
		pubSubClient: client,
	}
}

func (p *publisher) Publish(ctx context.Context, channel string, message interface{}) error {
	_, err := p.pubSubClient.Publish(ctx, channel, message).Result()
	return err
}
