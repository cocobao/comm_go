package rediscache

import (
	"encoding/json"
	"time"

	redis "gopkg.in/redis.v5"
)

var redisClient *redis.ClusterClient

func SetRedisClient(c *redis.ClusterClient) {
	redisClient = c
}

func IncrBy(key string, value int64) error {
	return redisClient.IncrBy(key, value).Err()
}

func TTL(key string) (time.Duration, error) {
	return redisClient.TTL(key).Result()
}

func CacheIsExist(key string) bool {
	return redisClient.Exists(key).Val()
}

func CacheSet(key string, val interface{}, timeOut time.Duration) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return redisClient.Set(key, string(data), timeOut).Err()
}

func CacheGet(key string, val interface{}) error {
	data, err := redisClient.Get(key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), &val)
}

func Set(key string, val interface{}, timeOut time.Duration) error {
	return redisClient.Set(key, val, timeOut).Err()
}

func Get(key string) (string, error) {
	return redisClient.Get(key).Result()
}

func CacheDel(key string) error {
	return redisClient.Del(key).Err()
}

func GetRedisClient() *redis.ClusterClient {
	return redisClient
}

func SessionSet(key string, field string, val interface{}) error {
	return redisClient.HSet(key, field, val).Err()
}

func SessionGetAll(key string) (map[string]string, error) {
	return redisClient.HGetAll(key).Result()
}

func SessionGet(key string, field string) *redis.StringCmd {
	return redisClient.HGet(key, field)
}

func SessionIsFieldExist(key string, field string) bool {
	return redisClient.HExists(key, field).Val()
}

func SessonFieldDelete(key string, field string) error {
	return redisClient.HDel(key, field).Err()
}
