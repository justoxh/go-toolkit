package redis

// Service implement interface for redis servece
type Service interface {
	// string
	Set(key string, value []byte, ttl int64) error
	Get(key string) ([]byte, error)
	Mset(keyvalues [][2][]byte, ttl int64) error
	Mget(keys []string) ([][]byte, error)
	Exists(key string) (bool, error)
	Del(key string) error
	Dels(keys []string) error
	Keys(keyFormat string) ([][]byte, error)

	// set
	SAdd(key string, ttl int64, members ...[]byte) error
	SRem(key string, members ...[]byte) (int64, error)
	SCard(key string) (int64, error)
	SIsMember(key string, member []byte) (bool, error)
	SMembers(key string) ([][]byte, error)

	// zset
	ZAdd(key string, ttl int64, args ...[]byte) error
	ZRem(key string, args ...[]byte) (int64, error)
	ZCard(key string) (int64, error)
	ZRank(key string, member []byte) (int64, error)
	ZRevRank(key string, member []byte) (int64, error)
	ZRange(key string, start, stop int64, withScores bool) ([][]byte, error)
	ZRangeByScore(key string, min, max interface{}, withScores bool) ([][]byte, error)
	ZRevRange(key string, start, stop int64, withScores bool) ([][]byte, error)

	// hash
	HSet(key string, ttl int64, field string, value []byte) error
	HGet(key string, field []byte) ([]byte, error)
	HMSet(key string, ttl int64, args ...[]byte) error
	HMGet(key string, fields ...[]byte) ([][]byte, error)
	HDel(key string, fields ...[]byte) (int64, error)
	HExists(key string, field []byte) (bool, error)
	HKeys(key string) ([][]byte, error)
	HVals(key string) ([][]byte, error)
	HGetAll(key string) ([][]byte, error)
	HLen(key string) (int64, error)

	// list
	LRpush(key string, args ...[]byte) (int64, error)
	LLpush(key string, args ...[]byte) (int64, error)
	LRpop(key string) ([]byte, error)
	LLpop(key string) ([]byte, error)
	LIndex(key string, index int64) ([]byte, error)
	LLlen(key string) (int64, error)
}
