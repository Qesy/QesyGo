package QesyGo

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

type CacheRedis struct {
	Pool     *redis.Pool // redis connection pool
	Conninfo string
	Auth     string
}

func (cr *CacheRedis) newPool() {
	cr.Pool = &redis.Pool{
		MaxIdle:     30,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", cr.Conninfo)
			if err != nil {
				return nil, err
			}
			if cr.Auth != "" {
				if _, err := c.Do("AUTH", cr.Auth); err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, err
		},
	}
}

//var pool = newPool()
func (cr *CacheRedis) Connect() error {
	cr.newPool()
	c := cr.Pool.Get()
	defer c.Close()
	return c.Err()
}

func (cr *CacheRedis) do(commandName string, args ...interface{}) (reply interface{}, err error) {
	c := cr.Pool.Get()
	defer c.Close()
	return c.Do(commandName, args...)
}

func (cr *CacheRedis) FlushAll() error {
	_, err := cr.do("FLUSHALL")
	return err
}

func (cr *CacheRedis) Get(key string) (string, error) {
	str, err := redis.String(cr.do("GET", key))
	return str, err
}

func (cr *CacheRedis) Set(key string, value string) error {
	_, err := cr.do("SET", key, value)
	return err
}

func (cr *CacheRedis) Del(key string) error {
	_, err := cr.do("DEL", key)
	return err
}

func (cr *CacheRedis) Exists(key string) (bool, error) {
	return redis.Bool(cr.do("EXISTS", key))
}

func (cr *CacheRedis) Expire(key string, second int) (bool, error) {
	return redis.Bool(cr.do("EXPIRE", key, second))
}

func (cr *CacheRedis) Keys(key string) ([]string, error) {
	return redis.Strings(cr.do("KEYS", key))
}

func (cr *CacheRedis) Ttl(key string) (interface{}, error) {
	return redis.Int(cr.do("TTL", key))
}

func (cr *CacheRedis) HMset(key string, arr interface{}) (interface{}, error) {
	return cr.do("HMSET", redis.Args{}.Add(key).AddFlat(arr)...)
}

func (cr *CacheRedis) HGet(key string, subKey string) (interface{}, error) {
	return cr.do("HGET", key, subKey)
}

func (cr *CacheRedis) HGetAll(key string) (map[string]string, error) {
	rsByte, err := cr.do("HGETALL", key)
	if err != nil {
		return nil, err
	}
	rs, err := redis.StringMap(rsByte, err)
	if err != nil {
		return nil, err
	}
	return rs, err
}

func (cr *CacheRedis) ZAdd(Key string, Score int64, Name string) (interface{}, error) {
	return cr.do("ZADD", Key, Score, Name)
}

func (cr *CacheRedis) ZCard(Key string) (int, error) {
	return redis.Int(cr.do("ZCARD", Key))
}

func (cr *CacheRedis) ZRange(Key string, Start int32, End int32) ([]map[string]string, error) {
	rank := []map[string]string{}
	if rsByte, err := cr.do("ZRANGE", Key, Start, End, "WITHSCORES"); err == nil {
		rsByteArr := rsByte.([]interface{})
		for i := 0; i < len(rsByteArr)/2; i++ {
			index := i * 2
			key := string(rsByteArr[index].([]byte))
			val := string(rsByteArr[index+1].([]byte))
			rank = append(rank, map[string]string{key: val})
		}
		return rank, err
	} else {
		return rank, err
	}
}

func (cr *CacheRedis) ZRevRange(Key string, Start int32, End int32) (interface{}, error) {
	rank := []map[string]string{}
	if rsByte, err := cr.do("ZREVRANGE", Key, Start, End, "WITHSCORES"); err == nil {
		rsByteArr := rsByte.([]interface{})
		for i := 0; i < len(rsByteArr)/2; i++ {
			index := i * 2
			key := string(rsByteArr[index].([]byte))
			val := string(rsByteArr[index+1].([]byte))
			rank = append(rank, map[string]string{key: val})
		}
		return rank, err
	} else {
		return rank, err
	}
}
