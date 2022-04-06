// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import (
	"github.com/go-redis/redis/v8"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
)

type r struct {
	State        cfg.State64
	StoreClient  *redis.Client
	PubSubClient *redis.Client
}

var Reds = new(r)

func (r *r) Boot() (err error) {
	if !r.State.TryBoot() {
		return r.stateError()
	}

	defer r.State.BootOrFail(&err)

	r.PubSubClient = redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Pubsub.Addr,
		Password: string(cfg.Redis.Pubsub.Pass),
		PoolSize: cfg.Redis.Pubsub.Pool,
	})

	r.StoreClient = redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Store.Addr,
		Password: string(cfg.Redis.Store.Pass),
		PoolSize: cfg.Redis.Store.Pool,
	})

	ctx := cfg.Context()

	err = r.PubSubClient.Ping(ctx).Err()

	if err != nil {
		r.abort()
		return
	}

	err = r.StoreClient.Ping(ctx).Err()

	if err != nil {
		r.abort()
		return
	}

	r.State.SetRunning()

	return nil
}

func (r *r) Stop() error {
	if r.State.TryStop() {
		err := r.State.StopOrFail(r.PubSubClient.Close)
		if err != nil {
			return err
		}
		err = r.State.StopOrFail(r.StoreClient.Close)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *r) stateError() error {
	return r.State.StateError("redis")
}

func (r *r) abort() {
	if r.PubSubClient != nil {
		_ = r.PubSubClient.Close()
	}
	if r.StoreClient != nil {
		_ = r.StoreClient.Close()
	}
}
