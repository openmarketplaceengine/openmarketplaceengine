// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package cfg

const (
	redisMinPool = 0
	redisMaxPool = 1024
	redisMinPort = 80
)

// RedisConfig wraps Redis store and pub-sub connection params.
type RedisConfig struct {
	Store  RedisConnect
	Pubsub RedisConnect
}

// RedisConnect contains necessary params to connect to a Redis server.
type RedisConnect struct {
	Pool      int      `default:"10" usage:"Redis server pool #size#"`
	Addr      string   `default:"127.0.0.1:6379" usage:"Redis server #host:port# address"`
	Pass      Password `usage:"Redis server #password#"`
	TLSServer string   `usage:"Redis tls server name, optional"`
}

// Check validates RedisConfig.
func (c *RedisConfig) Check(name ...string) (err error) {
	if err = c.Store.Check(append(name, "store")); err == nil {
		err = c.Pubsub.Check(append(name, "pubsub"))
	}
	return
}

// Check validates RedisConnect.
func (c *RedisConnect) Check(name []string) (err error) {
	err = checkRange(c.Pool, redisMinPool, redisMaxPool, append(name, "pool"))
	if err == nil {
		err = checkAddr(c.Addr, false, redisMinPort, name)
	}
	return
}
