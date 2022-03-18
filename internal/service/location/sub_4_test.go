package location

import (
	"encoding/json"

	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate"
	"github.com/openmarketplaceengine/openmarketplaceengine/redis/subscriber"
	"golang.org/x/net/context"
)

func subscribe(channel string) <-chan tollgate.Crossing {
	s := subscriber.NewSubscriber()

	messages := make(chan string)
	s.Subscribe(context.Background(), channel, messages)

	crossings := make(chan tollgate.Crossing, 1)

	go func() {
		for m := range messages {
			var crossing tollgate.Crossing
			_ = json.Unmarshal([]byte(m), &crossing)
			crossings <- crossing
		}
	}()

	return crossings
}
