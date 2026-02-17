package redis

import rdb "github.com/redis/go-redis/v9"

func CreateRedisClient(addr string) *rdb.Client {
	return rdb.NewClient(&rdb.Options{
		Addr: addr,
	})
}
