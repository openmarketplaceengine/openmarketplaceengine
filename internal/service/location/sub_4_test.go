package location

import (
	"encoding/json"

	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate/detector"

	"context"

	"github.com/openmarketplaceengine/openmarketplaceengine/redis/subscriber"
)

func subscribe(channel string) <-chan detector.Crossing {
	s := subscriber.NewSubscriber()

	messages := make(chan string)
	s.Subscribe(context.Background(), channel, messages)

	crossings := make(chan detector.Crossing, 1)

	go func() {
		for m := range messages {
			var crossing detector.Crossing
			_ = json.Unmarshal([]byte(m), &crossing)
			crossings <- crossing
		}
	}()

	return crossings
}
