package location

import (
	"encoding/json"

	"github.com/openmarketplaceengine/openmarketplaceengine/dao"

	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate/detector"

	"context"
)

func subscribe(channel string) <-chan detector.Crossing {
	s := dao.Reds.PubSubClient

	pubSub := s.PSubscribe(context.Background(), channel)

	crossings := make(chan detector.Crossing, 1)

	go func() {
		for m := range pubSub.Channel() {
			var crossing detector.Crossing
			_ = json.Unmarshal([]byte(m.Payload), &crossing)
			crossings <- crossing
		}
	}()

	return crossings
}
