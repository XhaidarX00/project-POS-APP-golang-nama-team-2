package database

import (
	"context"
	"fmt"
	"project_pos_app/config"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cache struct {
	rdb     *redis.Client
	expired time.Duration
	prefix  string
}

func newRedisClient(url, password string, dbIndex int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     url,
		Password: password,
		DB:       dbIndex,
	})
}

func NewCache(cfg config.Config, expired int) Cache {
	return Cache{
		rdb:     newRedisClient(cfg.Redis.Url, cfg.Redis.Password, 0),
		expired: time.Duration(expired) * time.Second,
		prefix:  cfg.Redis.Prefix,
	}
}

func (c *Cache) Push(name string, value []byte) error {
	return c.rdb.RPush(context.Background(), c.prefix+"_"+name, value).Err()
}

func (c *Cache) Pop(name string) (string, error) {
	return c.rdb.LPop(context.Background(), c.prefix+"_"+name).Result()
}

func (c *Cache) GetLength(name string) int64 {
	return c.rdb.LLen(context.Background(), c.prefix+"_"+name).Val()
}

func (c *Cache) Set(name string, value string) error {
	return c.rdb.Set(context.Background(), c.prefix+"_"+name, value, c.expired).Err()
}

func (c *Cache) SaveToken(name string, value string) error {
	return c.rdb.Set(context.Background(), c.prefix+"_"+name, value, 24*time.Hour).Err()
}

func (c *Cache) Get(name string) (string, error) {
	return c.rdb.Get(context.Background(), c.prefix+"_"+name).Result()
}

func (c *Cache) Delete(name string) error {
	return c.rdb.Del(context.Background(), c.prefix+"_"+name).Err()
}

func (c *Cache) DeleteByKey(key string) error {
	return c.rdb.Del(context.Background(), key).Err()
}

func (c *Cache) PrintKeys() {
	var cursor uint64
	for {
		var keys []string
		var err error
		keys, cursor, err = c.rdb.Scan(context.Background(), cursor, "", 0).Result()
		if err != nil {
			panic(err)
		}

		for _, key := range keys {
			fmt.Println("key", key)
		}

		if cursor == 0 { // no more keys
			break
		}
	}
}

func (c *Cache) GetKeys() []string {
	var cursor uint64
	var result []string
	for {
		var keys []string
		var err error
		keys, cursor, err = c.rdb.Scan(context.Background(), cursor, "", 0).Result()
		if err != nil {
			panic(err)
		}

		result = append(result, keys...)

		if cursor == 0 { // no more keys
			break
		}
	}

	return result
}

func (c *Cache) GetKeysByPattern(pattern string) []string {
	var cursor uint64
	var result []string
	for {
		var keys []string
		var err error
		keys, cursor, err = c.rdb.Scan(context.Background(), cursor, pattern, 0).Result()
		if err != nil {
			panic(err)
		}

		result = append(result, keys...)

		if cursor == 0 { // no more keys
			break
		}
	}

	return result
}

// Pub and Sub
func (c *Cache) Publish(channelName string, message string) error {
	return c.rdb.Publish(context.Background(), channelName, message).Err()
}

func (c *Cache) Subcribe(channelName string) (*redis.Message, error) {
	subscriber := c.rdb.Subscribe(context.Background(), channelName)
	message, err := subscriber.ReceiveMessage(context.Background())
	return message, err
}

func (c *Cache) GetClient() *redis.Client {
	return c.rdb
}
