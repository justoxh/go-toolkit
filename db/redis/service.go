package redis

import (
	"errors"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/s3dteam/go-toolkit/log"
)

type RedisOptions struct {
	Host        string
	Port        string
	Password    string
	IdleTimeout int
	MaxIdle     int
	MaxActive   int
}

// RedisCacheService redis service define
type RedisCacheService struct {
	options RedisOptions
	pool    *redis.Pool
	log     log.Logger
}

// Initialize redis service init
func (service *RedisCacheService) Initialize(redisCfg interface{}, l log.Logger) {
	options, ok := redisCfg.(RedisOptions)
	if !ok {
		panic("redis service config error!")
	}
	service.log = l
	service.options = options
	service.pool = &redis.Pool{
		IdleTimeout: time.Duration(options.IdleTimeout) * time.Second,
		MaxIdle:     options.MaxIdle,
		MaxActive:   options.MaxActive,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			address := fmt.Sprintf("%s:%s", options.Host, options.Port)
			var (
				c   redis.Conn
				err error
			)
			if c, err = redis.Dial("tcp", address); err != nil {
				service.log.Error("redis conn fail", "conn", address, "error", err)
				return nil, err
			}
			if len(options.Password) > 0 {
				if _, err = c.Do("AUTH", options.Password); err != nil {
					c.Close()
					service.log.Error("redis auth fail", "auth", options.Password, "error", err)
					return nil, err
				}
			}
			return c, nil
		},
	}
}

// Start xx
func (service *RedisCacheService) Start() {

}

// Stop xx
func (service *RedisCacheService) Stop() {

}

// -----------------string operation------------------
// when set exist key, old key ttl must reset it
func (service *RedisCacheService) Set(key string, value []byte, ttl int64) error {
	conn := service.pool.Get()
	defer conn.Close()

	if _, err := conn.Do("set", key, value); err != nil {
		service.log.Error("redis String set", "key", key, "error", err.Error())
		return err
	}

	if ttl > 0 {
		if _, err := conn.Do("expire", key, ttl); err != nil {
			service.log.Error("redis String expire", "key", key, "ttl", ttl, "error", err.Error())
			return err
		}
	}
	return nil
}

func (service *RedisCacheService) Get(key string) ([]byte, error) {
	conn := service.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("get", key)
	if nil != err {
		service.log.Error("redis String get", "key", key, "error", err.Error())
		return []byte{}, err
	} else if nil == reply {
		service.log.Debug("redis String get, no this key", "key", key)
		return []byte{}, err
	} else {
		return reply.([]byte), err
	}
}

func (service *RedisCacheService) Mset(keyvalues [][2][]byte, ttl int64) error {
	conn := service.pool.Get()
	defer conn.Close()

	if len(keyvalues) <= 0 {
		return nil
	}

	var list []interface{}
	for _, v := range keyvalues {
		list = append(list, v[0])
		list = append(list, v[1])
	}

	if _, err := conn.Do("mset", list...); err != nil {
		service.log.Error("redis String mset", "keyvalues", keyvalues, "error", err.Error())
		return err
	}

	if ttl > 0 {
		for _, kv := range keyvalues {
			if _, err := conn.Do("expire", kv[0], ttl); err != nil {
				service.log.Error("redis String expire", "key", kv[0], "error", err.Error())
				return err
			}
		}
	}
	return nil
}

func (service *RedisCacheService) Mget(keys []string) ([][]byte, error) {
	conn := service.pool.Get()
	defer conn.Close()

	if keys == nil || len(keys) <= 0 {
		return [][]byte{}, nil
	}
	var list []interface{}
	for _, v := range keys {
		list = append(list, v)
	}
	reply, err := conn.Do("mget", list...)
	res := [][]byte{}
	if nil != err {
		service.log.Error("redis String mget", "key", keys, "error", err.Error())
		return [][]byte{}, err
	} else if nil == reply {
		service.log.Debug("redis String get, no this key", "key", keys)
		return [][]byte{}, err
	} else {
		rs := reply.([]interface{})
		for _, r := range rs {
			if nil != r {
				res = append(res, r.([]byte))
			} else {
				res = append(res, nil)
			}
		}
		return res, err
	}
}

func (service *RedisCacheService) Exists(key string) (bool, error) {
	conn := service.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("exists", key)
	if err != nil {
		service.log.Error("redis String exists", "key", key, "error", err.Error())
		return false, err
	} else {
		exists, _ := reply.(int64)
		if exists == 1 {
			return true, nil
		} else {
			return false, nil
		}
	}
}

func (service *RedisCacheService) Del(key string) error {
	conn := service.pool.Get()
	defer conn.Close()

	_, err := conn.Do("del", key)
	if nil != err {
		service.log.Error("redis String del", "key", key, "error", err.Error())
	}
	return err
}

func (service *RedisCacheService) Dels(keys []string) error {
	conn := service.pool.Get()
	defer conn.Close()

	if keys == nil || len(keys) <= 0 {
		return nil
	}

	var list []interface{}
	for _, v := range keys {
		list = append(list, v)
	}

	num, err := conn.Do("del", list...)
	if err != nil {
		service.log.Error("redis String dels", "key", list, "error", err.Error())
	} else {
		service.log.Debug("redis String dels", "num", num.(int64))
	}

	return nil
}

func (service *RedisCacheService) Keys(keyFormat string) ([][]byte, error) {
	conn := service.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("keys", keyFormat)

	res := [][]byte{}
	if nil != err {
		service.log.Error("redis String keys", "error", err.Error())
	} else if nil != reply {
		rs := reply.([]interface{})
		for _, r := range rs {
			if nil != r {
				res = append(res, r.([]byte))
			}
		}
	}
	return res, err
}

// -----------------set operation---------------------
func (service *RedisCacheService) SAdd(key string, ttl int64, members ...[]byte) error {
	conn := service.pool.Get()
	defer conn.Close()

	vs := []interface{}{}
	vs = append(vs, key)
	for _, v := range members {
		vs = append(vs, v)
	}
	_, err := conn.Do("sadd", vs...)
	if nil != err {
		service.log.Error("redis set sadd", "key", key, "error", err.Error())
	}
	if ttl > 0 {
		if _, err := conn.Do("expire", key, ttl); err != nil {
			service.log.Error("redis set expire", "key", key, "error", err.Error())
			return err
		}
	}
	return err
}

func (service *RedisCacheService) SRem(key string, members ...[]byte) (int64, error) {
	conn := service.pool.Get()
	defer conn.Close()

	vs := []interface{}{}
	vs = append(vs, key)
	for _, v := range members {
		vs = append(vs, v)
	}
	reply, err := conn.Do("srem", vs...)

	if err != nil {
		service.log.Error("redis set srem", "key", key, "error", err.Error())
		return 0, err
	} else {
		res := reply.(int64)
		return res, err
	}
}

func (service *RedisCacheService) SCard(key string) (int64, error) {
	conn := service.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("scard", key)

	if err != nil {
		service.log.Error("redis set scard", "key", key, "error", err.Error())
		return 0, err
	} else {
		res := reply.(int64)
		return res, err
	}
}

func (service *RedisCacheService) SIsMember(key string, member []byte) (bool, error) {
	conn := service.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("sismember", key, member)
	if err != nil {
		service.log.Error("redis set sismember", "key", key, "member", member, "error", err.Error())
		return false, err
	} else {
		return reply.(int64) > 0, nil
	}
}

func (service *RedisCacheService) SMembers(key string) ([][]byte, error) {
	conn := service.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("smembers", key)

	res := [][]byte{}
	if nil != err {
		service.log.Error("redis set smembers", "key", key, "error", err.Error())
	} else if nil != reply {
		rs := reply.([]interface{})
		for _, r := range rs {
			res = append(res, r.([]byte))
		}
	}
	return res, err
}

// -----------------zset operation-------------------
func (service *RedisCacheService) ZAdd(key string, ttl int64, args ...[]byte) error {
	conn := service.pool.Get()
	defer conn.Close()

	if len(args)%2 != 0 {
		return errors.New("the length of `args` must be even")
	}
	vs := []interface{}{}
	vs = append(vs, key)
	for _, v := range args {
		vs = append(vs, v)
	}
	_, err := conn.Do("zadd", vs...)
	if nil != err {
		service.log.Error("redis zset zadd", "key", key, "error", err.Error())
	}
	if ttl > 0 {
		if _, err := conn.Do("expire", key, ttl); err != nil {
			service.log.Error("redis zset expire", "key", key, "error", err.Error())
			return err
		}
	}
	return err
}

func (service *RedisCacheService) ZRem(key string, args ...[]byte) (int64, error) {
	conn := service.pool.Get()
	defer conn.Close()

	vs := []interface{}{}
	vs = append(vs, key)
	for _, v := range args {
		vs = append(vs, v)
	}
	reply, err := conn.Do("zrem", vs...)
	if nil != err {
		service.log.Error("redis zset zrem", "key", key, "error", err.Error())
		return 0, err
	} else {
		res := reply.(int64)
		return res, err
	}
}

func (service *RedisCacheService) ZCard(key string) (int64, error) {
	conn := service.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("zcard", key)
	if nil != err {
		service.log.Error("redis zset zcard", "key", key, "error", err.Error())
		return 0, err
	} else {
		res := reply.(int64)
		return res, err
	}
}

func (service *RedisCacheService) ZRank(key string, member []byte) (int64, error) {
	conn := service.pool.Get()
	defer conn.Close()

	vs := []interface{}{}
	vs = append(vs, key)
	vs = append(vs, member)

	reply, err := conn.Do("zrank", vs...)
	if nil != err {
		service.log.Error("redis zset zrank", "key", key, "member", member, "error", err.Error())
		return -1, err
	} else if reply == nil {
		return -1, err
	} else {
		res := reply.(int64)
		return res, err
	}
}

func (service *RedisCacheService) ZRevRank(key string, member []byte) (int64, error) {
	conn := service.pool.Get()
	defer conn.Close()

	vs := []interface{}{}
	vs = append(vs, key)
	vs = append(vs, member)

	reply, err := conn.Do("zrevrank", vs...)
	if nil != err {
		service.log.Error("redis zset zrevrank", "key", key, "member", member, "error", err.Error())
		return -1, err
	} else if reply == nil {
		return -1, err
	} else {
		res := reply.(int64)
		return res, err
	}
}

func (service *RedisCacheService) ZRange(key string, start, stop int64, withScores bool) ([][]byte, error) {
	conn := service.pool.Get()
	defer conn.Close()

	vs := []interface{}{}
	vs = append(vs, key, start, stop)
	if withScores {
		vs = append(vs, []byte("WITHSCORES"))
	}
	reply, err := conn.Do("zrange", vs...)

	res := [][]byte{}
	if nil != err {
		service.log.Error("redis zset zrange", "key", key, "range", start, stop, "error", err.Error())
	} else if nil == err && nil != reply {
		rs := reply.([]interface{})
		for _, r := range rs {
			if nil == r {
				res = append(res, []byte{})
			} else {
				res = append(res, r.([]byte))
			}
		}
	}
	return res, err
}

func (service *RedisCacheService) ZRangeByScore(key string, min, max interface{}, withScores bool) ([][]byte, error) {
	conn := service.pool.Get()
	defer conn.Close()

	vs := []interface{}{}
	vs = append(vs, key, min, max)
	if withScores {
		vs = append(vs, []byte("WITHSCORES"))
	}
	reply, err := conn.Do("zrangebyscore", vs...)

	res := [][]byte{}
	if nil != err {
		service.log.Error("redis zset zrangebyscore", "key", key, "range", min, max, "error", err.Error())
	} else if nil == err && nil != reply {
		rs := reply.([]interface{})
		for _, r := range rs {
			if nil == r {
				res = append(res, []byte{})
			} else {
				res = append(res, r.([]byte))
			}
		}
	}
	return res, err
}

func (service *RedisCacheService) ZRevRange(key string, start, stop int64, withScores bool) ([][]byte, error) {
	conn := service.pool.Get()
	defer conn.Close()

	vs := []interface{}{}
	vs = append(vs, key, start, stop)
	if withScores {
		vs = append(vs, []byte("WITHSCORES"))
	}
	reply, err := conn.Do("zrevrange", vs...)

	res := [][]byte{}
	if nil != err {
		service.log.Error("redis zset zrevrange", "key", key, "range", start, stop, "error", err.Error())
	} else if nil == err && nil != reply {
		rs := reply.([]interface{})
		for _, r := range rs {
			if nil == r {
				res = append(res, []byte{})
			} else {
				res = append(res, r.([]byte))
			}
		}
	}
	return res, err
}

func (service *RedisCacheService) ZRemRangeByScore(key string, start, stop int64) (int64, error) {

	//log.Info("[REDIS-ZRemRangeByScore] key : " + key)

	conn := service.pool.Get()
	defer conn.Close()

	vs := []interface{}{}
	vs = append(vs, key, start, stop)

	reply, err := conn.Do("ZREMRANGEBYSCORE", vs...)

	if err != nil {
		service.log.Error(" key:%s, err:%s", key, err.Error())
		return 0, err
	} else {
		res := reply.(int64)
		return res, err
	}
}

// -----------------hash operation-------------------
func (service *RedisCacheService) HSet(key string, ttl int64, field string, value []byte) error {
	conn := service.pool.Get()
	defer conn.Close()

	vs := []interface{}{}
	vs = append(vs, key)
	vs = append(vs, field)
	vs = append(vs, value)

	_, err := conn.Do("hset", vs...)
	if nil != err {
		service.log.Error("redis hash hset", "key", key, "filed", field, "value", value, "error", err.Error())
	}
	if ttl > 0 {
		if _, err := conn.Do("expire", key, ttl); err != nil {
			service.log.Error("redis hash hset", "key", key, "filed", field, "error", err.Error())
			return err
		}
	}
	return err
}

func (service *RedisCacheService) HGet(key string, field []byte) ([]byte, error) {
	conn := service.pool.Get()
	defer conn.Close()

	var vs []interface{}
	vs = append(vs, key)
	vs = append(vs, field)

	reply, err := conn.Do("hget", vs...)

	if nil != err {
		service.log.Error("redis hash hget", "key", key, "filed", field, "error", err.Error())
		return []byte{}, err
	} else if nil == reply {
		service.log.Debug("redis hash hget", "key", key, "filed", field)
		return []byte{}, err
	} else {
		return reply.([]byte), err
	}
}

func (service *RedisCacheService) HMSet(key string, ttl int64, args ...[]byte) error {
	conn := service.pool.Get()
	defer conn.Close()

	if len(args)%2 != 0 {
		return errors.New("the length of `args` must be even")
	}
	vs := []interface{}{}
	vs = append(vs, key)
	for _, v := range args {
		vs = append(vs, v)
	}
	_, err := conn.Do("hmset", vs...)
	if nil != err {
		service.log.Error("redis hash hmset", "key", key, "error", err.Error())
	}
	if ttl > 0 {
		if _, err := conn.Do("expire", key, ttl); err != nil {
			service.log.Error("redis hash hmset", "key", key, "error", err.Error())
			return err
		}
	}
	return err
}

func (service *RedisCacheService) HMGet(key string, fields ...[]byte) ([][]byte, error) {
	conn := service.pool.Get()
	defer conn.Close()

	vs := []interface{}{}
	vs = append(vs, key)
	for _, v := range fields {
		vs = append(vs, v)
	}
	reply, err := conn.Do("hmget", vs...)

	res := [][]byte{}
	if nil != err {
		service.log.Error("redis hash hmget", "key", key, "error", err.Error())
	} else if nil == err && nil != reply {
		rs := reply.([]interface{})
		for _, r := range rs {
			if nil == r {
				res = append(res, []byte{})
			} else {
				res = append(res, r.([]byte))
			}
		}
	}
	return res, err
}

func (service *RedisCacheService) HDel(key string, fields ...[]byte) (int64, error) {
	conn := service.pool.Get()
	defer conn.Close()

	vs := []interface{}{}
	vs = append(vs, key)
	for _, v := range fields {
		vs = append(vs, v)
	}
	reply, err := conn.Do("hdel", vs...)

	if err != nil {
		service.log.Error("redis hash hdel", "key", key, "fields", fields, "error", err.Error())
		return 0, err
	} else {
		res := reply.(int64)
		return res, err
	}
}

func (service *RedisCacheService) HExists(key string, field []byte) (bool, error) {
	conn := service.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("hexists", key, field)
	if nil != err {
		service.log.Error("redis hash hexists", "key", key, "field", field, "error", err.Error())
	} else if nil == err && nil != reply {
		exists := reply.(int64)
		return exists > 0, nil
	}

	return false, err
}

func (service *RedisCacheService) HKeys(key string) ([][]byte, error) {
	conn := service.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("hkeys", key)

	res := [][]byte{}
	if nil != err {
		service.log.Error("redis hash hkeys", "key", key, "error", err.Error())
	} else if nil != reply {
		rs := reply.([]interface{})
		for _, r := range rs {
			res = append(res, r.([]byte))
		}
	}
	return res, err
}

func (service *RedisCacheService) HVals(key string) ([][]byte, error) {
	conn := service.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("hvals", key)

	res := [][]byte{}
	if nil != err {
		service.log.Error("redis hash hvals", "key", key, "error", err.Error())
	} else if nil != reply {
		rs := reply.([]interface{})
		for _, r := range rs {
			res = append(res, r.([]byte))
		}
	}
	return res, err
}

func (service *RedisCacheService) HGetAll(key string) ([][]byte, error) {
	conn := service.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("hgetall", key)

	res := [][]byte{}
	if nil != err {
		service.log.Error("redis hash hgetall", "key", key, "error", err.Error())
	} else if nil != reply {
		rs := reply.([]interface{})
		for _, r := range rs {
			res = append(res, r.([]byte))
		}
	}
	return res, err
}

func (service *RedisCacheService) HLen(key string) (int64, error) {
	conn := service.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("hlen", key)

	if nil != err {
		service.log.Error("redis hash hlen", "key", key, "error", err.Error())
		return 0, err
	} else {
		res := reply.(int64)
		return res, err
	}
}

// -----------------list operation--------------------
func (service *RedisCacheService) LRpush(key string, args ...[]byte) (int64, error) {
	conn := service.pool.Get()
	defer conn.Close()

	vs := []interface{}{}
	vs = append(vs, key)
	for _, v := range args {
		vs = append(vs, v)
	}
	reply, err := conn.Do("lpush", vs...)

	if err != nil {
		service.log.Error("redis list lpush", "key", key, "value", args, "error", err.Error())
		return 0, err
	} else {
		res := reply.(int64)
		return res, err
	}
}

func (service *RedisCacheService) LLpush(key string, args ...[]byte) (int64, error) {
	conn := service.pool.Get()
	defer conn.Close()

	vs := []interface{}{}
	vs = append(vs, key)
	for _, v := range args {
		vs = append(vs, v)
	}
	reply, err := conn.Do("rpush", vs...)

	if err != nil {
		service.log.Error("redis list rpush", "key", key, "value", args, "error", err.Error())
		return 0, err
	} else {
		res := reply.(int64)
		return res, err
	}
}

func (service *RedisCacheService) LRpop(key string) ([]byte, error) {
	conn := service.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("rpop", key)

	if nil != err {
		service.log.Error("redis list rpop", "key", key, "error", err.Error())
		return []byte{}, err
	} else if nil == reply {
		service.log.Debug("redis list rpop", "key", key)
		return []byte{}, err
	} else {
		return reply.([]byte), err
	}
}

func (service *RedisCacheService) LLpop(key string) ([]byte, error) {
	conn := service.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("lpop", key)

	if nil != err {
		service.log.Error("redis list lpop", "key", key, "error", err.Error())
		return []byte{}, err
	} else if nil == reply {
		service.log.Debug("redis list lpop", "key", key)
		return []byte{}, err
	} else {
		return reply.([]byte), err
	}
}

func (service *RedisCacheService) LIndex(key string, index int64) ([]byte, error) {
	conn := service.pool.Get()
	defer conn.Close()

	vs := []interface{}{}
	vs = append(vs, key)
	vs = append(vs, index)

	reply, err := conn.Do("lindex", vs)

	if nil != err {
		service.log.Error("redis list lindex", "key", key, "index", index, "error", err.Error())
		return []byte{}, err
	} else if nil == reply {
		service.log.Debug("redis list lindex", "key", key, "index", index)
		return []byte{}, err
	} else {
		return reply.([]byte), err
	}
}
