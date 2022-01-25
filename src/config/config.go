package config

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var (
	v             *viper.Viper
	propertyNames map[string]struct{}
)

func init() {
	v = viper.New()
}

const EnvPrefix = "OME"

const (
	ServicePort         = "service.port"
	RedisStorePool      = "redis.store.pool"
	RedisStoreAddr      = "redis.store.addr"
	RedisStorePassword  = "redis.store.password"
	RedisPubSubPool     = "redis.pubsub.pool"
	RedisPubSubAddr     = "redis.pubsub.addr"
	RedisPubSubPassword = "redis.pubsub.password"
)

func Read() error {
	v.SetEnvPrefix(EnvPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	v.SetConfigName("config")
	v.SetConfigType("yml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.AddConfigPath("../config")
	v.AddConfigPath("../../config")
	v.AddConfigPath("../../../config")
	v.AddConfigPath("./src/config")
	err := v.ReadInConfig()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error reading config file %s", v.ConfigFileUsed()))
	}

	propertyNames = make(map[string]struct{})
	for _, k := range v.AllKeys() {
		propertyNames[k] = struct{}{}
	}
	return nil
}

func requireNotMissing(key string) {
	_, ok := propertyNames[key]
	if !ok {
		panic(fmt.Sprintf("config %s is missing property %s", v.ConfigFileUsed(), key))
	}
}

func GetString(key string) string {
	requireNotMissing(key)
	value := v.GetString(key)
	return value
}

func GetInt(key string) int {
	requireNotMissing(key)
	value := v.GetInt(key)
	return value
}

func GetBool(key string) bool {
	requireNotMissing(key)
	return v.GetBool(key)
}
