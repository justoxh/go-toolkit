package redis

import (
	"bufio"
	"errors"
	"fmt"
	_ "github.com/justoxh/go-toolkit/log"
	"github.com/justoxh/go-toolkit/log/logruslogger"
	"os"
	"strconv"
	"testing"

	"github.com/naoina/toml"
)

type RelayConfig struct {
	Name string
	Role string

	Redis RedisOptions
}

func Float64ToByte(f float64) []byte {
	// 8 precision
	s := strconv.FormatFloat(f, 'f', 8, 64)
	return []byte(s)
}

func loadConfig(file string, cfg *RelayConfig) {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = toml.NewDecoder(bufio.NewReader(f)).Decode(cfg)
	// Add file name to errors that have a line number.
	if _, ok := err.(*toml.LineError); ok {
		err = errors.New(file + ", " + err.Error())
		panic(err)
	}
}

func InitService() *RedisCacheService {
	var config_relay1 *RelayConfig = &RelayConfig{}
	var config_file1 string = "./pool_test_nopass.toml"
	loadConfig(config_file1, config_relay1)
	var redisService *RedisCacheService = &RedisCacheService{}
	var conf = logruslogger.Options{
	}
	log := logruslogger.GetLoggerWithOptions("test",&conf)
	redisService.Initialize(config_relay1.Redis,log)
	return redisService
}

func Test_Initialize_nopass(t *testing.T) {
	var config_relay1 *RelayConfig = &RelayConfig{}
	var config_file1 string = "./pool_test_nopass.toml"
	loadConfig(config_file1, config_relay1)
	if config_relay1.Name != "miner" {
		t.Errorf("config item Name[%s] not equal miner", config_relay1.Name)
	}
	var redisService *RedisCacheService = &RedisCacheService{}
	var conf = logruslogger.Options{
	}
	log := logruslogger.GetLoggerWithOptions("test",&conf)
	redisService.Initialize(config_relay1.Redis,log)
	// test connection no password
	conn := redisService.pool.Get()
	defer conn.Close()
	_, err1 := conn.Do("PING")
	if err1 != nil {
		t.Errorf("redis service conn fail, no need pass but config pass")
	}
	// test connection close
	conn.Close()
	_, err2 := conn.Do("PING")
	if err2 == nil {
		t.Errorf("redis service close fail!")
	}
}

func Test_Initialize_pass(t *testing.T) {
	var config_relay1 *RelayConfig = &RelayConfig{}
	var config_file1 string = "./pool_test_pass.toml"
	loadConfig(config_file1, config_relay1)
	if config_relay1.Name != "miner" {
		t.Errorf("config item Name[%s] not equal miner", config_relay1.Name)
	}
	var redisService *RedisCacheService = &RedisCacheService{}
	var conf = logruslogger.Options{
	}
	log := logruslogger.GetLoggerWithOptions("test",&conf)
	redisService.Initialize(config_relay1.Redis,log)
	// test connection with password
	conn := redisService.pool.Get()
	defer conn.Close()
	_, err := conn.Do("PING")
	if err == nil {
		t.Errorf("redis service conn fail, need pass but not config pass")
	}
}

func Test_String_Set(t *testing.T) {
	var redisService *RedisCacheService = InitService()
	// test key nil, value not nil
	err1 := redisService.Set("", []byte("123"), 1)
	if err1 != nil {
		t.Errorf("redis string set[key nil, value not nil] fail")
	}
	// test key not nil, value nil
	err2 := redisService.Set("test2", nil, 1)
	if err2 != nil {
		t.Errorf("redis string set[key not nil, value nil] fail")
	}
	// test key not nil, value not nil & tll > 0
	err3 := redisService.Set("test3", []byte("123"), 1)
	if err3 != nil {
		t.Errorf("redis string set[key not nil, value nil] fail")
	}
	// test tll == 0
	err5 := redisService.Set("test5", []byte("12345"), 0)
	if err5 != nil {
		t.Errorf("redis string set[key not nil, value nil] fail")
	}
	redisService.Del("test5")
	// test tll < 0
	err6 := redisService.Set("test6", []byte("123456"), -1)
	if err6 != nil {
		t.Errorf("redis string set[key not nil, value nil] fail")
	}
	redisService.Del("test6")
}

func Test_String_Get(t *testing.T) {
	var redisService *RedisCacheService = InitService()
	// get nil key
	reply1, err1 := redisService.Get("")
	if err1 != nil {
		t.Errorf("redis string set[key is nil] fail")
	}
	if len(reply1) != 0 {
		t.Errorf("redis string set[key is nil, value is nil] fail")
	}
	// get not exist key
	reply2, err2 := redisService.Get("test2")
	if err2 != nil {
		t.Errorf("redis string get[key not exist] fail")
	}
	if len(reply2) != 0 {
		t.Errorf("redis string set[key not exist, value is nil] fail")
	}
	// get exist key
	err3 := redisService.Set("test3", []byte("123456"), 5)
	if err3 != nil {
		t.Errorf("redis string set fail")
	}
	reply3, err4 := redisService.Get("test3")
	if err4 != nil {
		t.Errorf("redis string get[key exist] fail")
	}
	if string(reply3) != "123456" {
		t.Errorf("redis string set[key exist, value is not nil] fail")
	}
}

func Test_String_Mset(t *testing.T) {
	redisService := InitService()
	// keyvalues len 0
	err1 := redisService.Mset([][2][]byte{}, 1)
	if err1 != nil {
		t.Errorf("redis string mset[keyvalue is nil] fail")
	}
	// keyvalues len%2 != 0
	var keyvalues2 [][2][]byte
	var arr1 [2][]byte
	arr1[0] = []byte("test1")
	arr1[1] = []byte("value1")
	var arr2 [2][]byte
	arr2[0] = []byte("test2")
	keyvalues2 = append(keyvalues2, arr1)
	keyvalues2 = append(keyvalues2, arr2)

	if err2 := redisService.Mset(keyvalues2, 10); err2 != nil {
		t.Errorf("redis string mset[keyvalue len/2 != 0] fail")
	}
	if relay21, err21 := redisService.Get("test1"); err21 != nil || string(relay21) != "value1" {
		t.Errorf("redis string mset[keyvalue len/2 != 0] value check fail")
	}
	conn := redisService.pool.Get()
	defer conn.Close()
	reply, _ := conn.Do("EXISTS ", "test2")
	if exists, _ := reply.(int64); exists != 0 {
		t.Errorf("redis string mset[keyvalue len/2 != 0] value check fail")
	}

	// keyvalues len%2 == 0
	var keyvalues3 [][2][]byte
	var arr3 [2][]byte
	arr3[0] = []byte("test3")
	arr3[1] = []byte("value3")
	var arr4 [2][]byte
	arr4[0] = []byte("test4")
	arr4[1] = []byte("value4")
	keyvalues3 = append(keyvalues3, arr3)
	keyvalues3 = append(keyvalues3, arr4)
	if err4 := redisService.Mset(keyvalues3, 10); err4 != nil {
		t.Errorf("redis string mset[keyvalue len/2 == 0] fail")
	}
	if relay31, err31 := redisService.Get("test3"); err31 != nil || string(relay31) != "value3" {
		t.Errorf("redis string mset[keyvalue len/2 != 0] value check fail")
	}
	if relay32, err32 := redisService.Get("test4"); err32 != nil || string(relay32) != "value4" {
		t.Errorf("redis string mset[keyvalue len/2 != 0] value check fail")
	}
}

func Test_String_Mget(t *testing.T) {
	redisService := InitService()
	// keys is nil
	reply1, err1 := redisService.Mget(nil)
	if len(reply1) != 0 || err1 != nil {
		t.Errorf("redis string mget[keys is nil] fail")
	}
	// keys len 0
	reply2, err2 := redisService.Mget([]string{})
	if len(reply2) != 0 || err2 != nil {
		t.Errorf("redis string mget[keys len = 0] fail")
	}
	// keys len > 0, all exist
	var keyvalues3 [][2][]byte
	var arr31 [2][]byte
	arr31[0] = []byte("test3")
	arr31[1] = []byte("value3")
	var arr32 [2][]byte
	arr32[0] = []byte("test4")
	arr32[1] = []byte("value4")
	keyvalues3 = append(keyvalues3, arr31)
	keyvalues3 = append(keyvalues3, arr32)
	redisService.Mset(keyvalues3, 10)
	reply3, err3 := redisService.Mget([]string{"test3", "test4"})
	if len(reply3) != 2 || err3 != nil {
		t.Errorf("redis string mget[keys len > 0 & all exist] fail")
	}
	if reply31, _ := redisService.Get("test3"); string(reply31) != "value3" {
		t.Errorf("redis string mget[keys len > 0 & all exist] check fail")
	}

	if reply32, _ := redisService.Get("test4"); string(reply32) != "value4" {
		t.Errorf("redis string mget[keys len > 0 & all exist] check fail")
	}
	// keys len > 0, keys part exist
	var keyvalues4 [][2][]byte
	var arr41 [2][]byte
	arr41[0] = []byte("test3")
	arr41[1] = []byte("value3")
	var arr42 [2][]byte
	arr42[0] = []byte("test4")
	arr42[1] = []byte("value4")
	keyvalues4 = append(keyvalues4, arr41)
	keyvalues4 = append(keyvalues4, arr42)
	redisService.Mset(keyvalues4, 10)
	reply4, err4 := redisService.Mget([]string{"test1", "test2", "test3", "test4"})
	if len(reply4) != 4 || err4 != nil {
		t.Errorf("redis string mget[keys len > 0 & part exist] fail")
	}
	if reply41, _ := redisService.Get("test2"); string(reply41) != "" {
		t.Errorf("redis string mget[keys len > 0 & part exist] check fail")
	}

	if reply42, _ := redisService.Get("test4"); string(reply42) != "value4" {
		t.Errorf("redis string mget[keys len > 0 & part exist] check fail")
	}
}

func Test_String_Exists(t *testing.T) {
	redisService := InitService()
	// key is "" not exist
	if b1, err1 := redisService.Exists(""); b1 != false || err1 != nil {
		t.Errorf("redis string exists[key is '' not exist] check fail")
	}
	// key is "" exist
	redisService.Set("", []byte("value2"), 10)
	if b2, err2 := redisService.Exists(""); b2 != true || err2 != nil {
		t.Errorf("redis string exists[key is '' exist] check fail")
	}
	// key not exist
	if b3, err3 := redisService.Exists("test3"); b3 != false || err3 != nil {
		t.Errorf("redis string exists[key not exist ] check fail")
	}
	// key exist
	redisService.Set("test4", []byte("value4"), 10)
	if b4, err4 := redisService.Exists("test4"); b4 != true || err4 != nil {
		t.Errorf("redis string exists[key exist] check fail")
	}
}

func Test_String_Del(t *testing.T) {
	redisService := InitService()
	// key not exist
	if err1 := redisService.Del("test1"); err1 != nil {
		t.Errorf("redis string exists[key not exist] check fail")
	}
	// key exist
	redisService.Set("test2", []byte("value2"), 10)
	if err2 := redisService.Del("test2"); err2 != nil {
		t.Errorf("redis string exists[key exist] check fail")
	}
	if b21, err21 := redisService.Exists("test2"); b21 == true || err21 != nil {
		t.Errorf("redis string exists[key exist] check fail")
	}
}

func Test_String_Dels(t *testing.T) {
	redisService := InitService()
	// keys is nil
	if err1 := redisService.Dels(nil); err1 != nil {
		t.Errorf("redis string Dels[keys is nil] check fail")
	}
	// keys len 0
	if err2 := redisService.Dels([]string{}); err2 != nil {
		t.Errorf("redis string Dels[keys len 0] check fail")
	}
	// keys len > 0, all exist
	redisService.Set("test1", []byte("value1"), 10)
	redisService.Set("test2", []byte("value2"), 10)
	if err3 := redisService.Dels([]string{"test1", "test2"}); err3 != nil {
		t.Errorf("redis string Dels[ keys len > 0, all exist] check fail")
	}
	// keys len > 0, part exist
	redisService.Set("test3", []byte("value3"), 10)
	if err4 := redisService.Dels([]string{"", "test3"}); err4 != nil {
		t.Errorf("redis string Dels[keys len > 0, part exist] check fail")
	}
}

func Test_String_Keys(t *testing.T) {
	redisService := InitService()
	// keyFormat is ""
	if reply1, err1 := redisService.Keys(""); len(reply1) != 0 || err1 != nil {
		t.Errorf("redis string keys[keyFormat is ''] check fail")
	}
	// keyFormat not exist
	if reply2, err2 := redisService.Keys("test1"); len(reply2) != 0 || err2 != nil {
		t.Errorf("redis string keys[keyFormat not exist] check fail")
	}
	// keyFormat exist
	redisService.Set("test1", []byte("value1"), 1)
	if reply3, err3 := redisService.Keys("test1"); len(reply3) != 1 || err3 != nil {
		t.Errorf("redis string keys[keyFormat not exist] check fail")
	}
	redisService.Del("test1")
	// keyFormat rule
	redisService.Set("test2", []byte("value2"), 10)
	redisService.Set("test3", []byte("value3"), 10)
	reply4, err4 := redisService.Keys("test*")
	fmt.Println(reply4)
	if len(reply4) != 2 || err4 != nil {
		t.Errorf("redis string keys[keyFormat not exist] check fail")
	}
}

func Test_Set_SAdd(t *testing.T) {
	redisService := InitService()
	// members is nil
	if err1 := redisService.SAdd("test1", 3, nil); err1 != nil {
		t.Errorf("redis set sadd[members is nil] fail")
	}
	if reply1, _ := redisService.SCard("test1"); reply1 != 1 {
		t.Errorf("redis set sadd[members is nil] check fail")
	}
	redisService.Del("test1")
	// members one
	if err2 := redisService.SAdd("test2", 3, []byte("value1")); err2 != nil {
		t.Errorf("redis set sadd[members one] fail")
	}
	if reply2, _ := redisService.SCard("test2"); reply2 != 1 {
		t.Errorf("redis set sadd[members one] check fail")
	}
	redisService.Del("test2")
	// members repeate one
	if err3 := redisService.SAdd("test3", 3, []byte("value1"), []byte("value1")); err3 != nil {
		t.Errorf("redis set sadd[members repeate one] fail")
	}
	if reply3, _ := redisService.SCard("test3"); reply3 != 1 {
		t.Errorf("redis set sadd[members repeate one] check fail")
	}
	redisService.Del("test3")
	// members one two
	if err4 := redisService.SAdd("test4", 3, []byte("value1"), []byte("value2")); err4 != nil {
		t.Errorf("redis set sadd[members one two] fail")
	}
	if reply4, _ := redisService.SCard("test4"); reply4 != 2 {
		t.Errorf("redis set sadd[members one two] check fail")
	}
	redisService.Del("test4")
}

func Test_Set_SRem(t *testing.T) {
	// redisService := InitService()
}

func Test_Set_SCard(t *testing.T) {
	redisService := InitService()
	// key not exist
	if reply1, err1 := redisService.SCard("test1"); reply1 != 0 || err1 != nil {
		t.Errorf("redis set scard[key not exist] fail")
	}
	// key exist
	if err2 := redisService.SAdd("test2", 3, []byte("value1"), []byte("value2")); err2 != nil {
		t.Errorf("redis set scard[key exist] add fail")
	}
	if reply3, err3 := redisService.SCard("test2"); reply3 != 2 || err3 != nil {
		t.Errorf("redis set scard[key exist] check fail")
	}
	redisService.Del("test2")
}

func Test_Set_SIsMember(t *testing.T) {
	redisService := InitService()
	// key not exist
	if reply1, err1 := redisService.SIsMember("test1", []byte("value1")); reply1 || err1 != nil {
		t.Errorf("redis set sismember[key not exist] fail")
	}
	// key exist, member not exist
	if err2 := redisService.SAdd("test2", 3, []byte("value2")); err2 != nil {
		t.Errorf("redis set sismember[key exist, member not exist] add fail")
	}
	if reply21, err21 := redisService.SIsMember("test2", []byte("value1")); reply21 || err21 != nil {
		t.Errorf("redis set sismember[key exist, member not exist] fail")
	}
	redisService.Del("test2")
	// key exist, member exist
	if err3 := redisService.SAdd("test3", 3, []byte("value3")); err3 != nil {
		t.Errorf("redis set sismember[key exist, member not exist] add fail")
	}
	if reply31, err31 := redisService.SIsMember("test3", []byte("value3")); !reply31 || err31 != nil {
		t.Errorf("redis set sismember[key exist, member not exist] fail")
	}
}

func Test_Set_SMembers(t *testing.T) {
	redisService := InitService()
	// key is ""
	if reply1, err1 := redisService.SMembers(""); len(reply1) != 0 || err1 != nil {
		t.Errorf("redis set smembers[key is ''] fail")
	}
	// key not exist
	if reply2, err2 := redisService.SMembers("test1"); len(reply2) != 0 || err2 != nil {
		t.Errorf("redis set smembers[key not exist] fail")
	}
	// key exist
	if err3 := redisService.SAdd("test3", 3, []byte("value3")); err3 != nil {
		t.Errorf("redis set sismember[key exist, member not exist] add fail")
	}
	if reply31, err31 := redisService.SMembers("test3"); len(reply31) != 1 || err31 != nil {
		t.Errorf("redis set smembers[key exist] fail")
	}
	redisService.Del("test3")
}

func Test_ZSet_ZAdd(t *testing.T) {
	redisService := InitService()
	// args is nil
	if err1 := redisService.ZAdd("test1", 5, nil); err1 == nil {
		t.Errorf("redis set zadd[args is nil] fail")
	}
	redisService.Del("test1")
	// args len == 0
	if err2 := redisService.ZAdd("test1", 5); err2 == nil {
		t.Errorf("redis set zadd[args len == 0] fail")
	}
	redisService.Del("test1")
	// args len == 1
	if err3 := redisService.ZAdd("test1", 5, []byte("value1")); err3 == nil {
		t.Errorf("redis set zadd[args len == 1] fail")
	}
	redisService.Del("test1")
	// args len == 2
	if err4 := redisService.ZAdd("test1", 5, []byte("value1"), []byte("value1")); err4 == nil {
		t.Errorf("redis set zadd[args len == 2] fail")
	}
	redisService.Del("test1")
	if err51 := redisService.ZAdd("test1", 5, []byte("200"), []byte("value1")); err51 != nil {
		t.Errorf("redis set zadd[args len == 2] fail")
	}
	if err52 := redisService.ZAdd("test1", 5, []byte("100"), []byte("value2")); err52 != nil {
		t.Errorf("redis set zadd[args len == 2] fail")
	}
	if reply51, _ := redisService.ZCard("test1"); reply51 != 2 {
		t.Errorf("redis set zadd[args len == 2] check fail")
	}
	redisService.Del("test1")
}

func Test_ZSet_ZRem(t *testing.T) {
	redisService := InitService()
	// key is ''
	if reply1, err1 := redisService.ZRem("", []byte("value1")); reply1 != 0 || err1 != nil {
		t.Errorf("redis set zrem[key is ''] fail")
	}
	// key not '' && member is ''
	if err2 := redisService.ZAdd("test2", 5, []byte("100"), []byte("value2")); err2 != nil {
		t.Errorf("redis set zrem[key not '' && member is ''] add fail")
	}
	if reply21, err21 := redisService.ZRem("test2", []byte("")); reply21 != 0 || err21 != nil {
		t.Errorf("redis set zrem[key not '' && member is ''] fail")
	}
	redisService.Del("test2")
	// key not '' && member not exist
	if err3 := redisService.ZAdd("test3", 5, []byte("100"), []byte("value2")); err3 != nil {
		t.Errorf("redis set zrem[key not '' && member is ''] add fail")
	}
	if reply31, err31 := redisService.ZRem("test3", []byte("value3")); reply31 != 0 || err31 != nil {
		t.Errorf("redis set zrem[key not '' && member is ''] fail")
	}
	redisService.Del("test3")
	// key not '' && member exist del one
	if err4 := redisService.ZAdd("test4", 5, []byte("100"), []byte("value2")); err4 != nil {
		t.Errorf("redis set zrem[key not '' && member exist del one ] add fail")
	}
	if reply41, err41 := redisService.ZRem("test4", []byte("value2")); reply41 != 1 || err41 != nil {
		t.Errorf("redis set zrem[key not '' && member exist del one ] fail")
	}
	redisService.Del("test4")
	// key not '' && member exist del more
	if err5 := redisService.ZAdd("test5", 5, []byte("200"), []byte("value1"), []byte("100"), []byte("value2")); err5 != nil {
		t.Errorf("redis set zadd[key not '' && member exist del more] add fail")
	}
	if reply51, err51 := redisService.ZRem("test5", []byte("value1"), []byte("value2")); reply51 != 2 || err51 != nil {
		t.Errorf("redis set zrem[key not '' && member exist del more fail")
	}
	redisService.Del("test5")
}

func Test_ZSet_ZCard(t *testing.T) {
	redisService := InitService()
	// key is ''
	if reply1, err1 := redisService.ZCard(""); reply1 != 0 || err1 != nil {
		t.Errorf("redis set zcard[key is ''] fail")
	}
	// key is not '' && key not exist
	if reply2, err2 := redisService.ZCard("test1"); reply2 != 0 || err2 != nil {
		t.Errorf("redis set zcard[key is not '' && key not exist] fail")
	}
	// key is not '' && key exist one
	if err3 := redisService.ZAdd("test2", 5, []byte("100"), []byte("value2")); err3 != nil {
		t.Errorf("redis set zcard[key is not '' && key exist one] add fail")
	}
	if reply31, err31 := redisService.ZCard("test2"); reply31 != 1 || err31 != nil {
		t.Errorf("redis set zcard[key is not '' && key exist one] fail")
	}
	redisService.Del("test2")
	// key is not '' && key exist one more
	if err4 := redisService.ZAdd("test3", 5, []byte("200"), []byte("value1"), []byte("100"), []byte("value2")); err4 != nil {
		t.Errorf("redis set zcard[key is not '' && key exist one more] add fail")
	}
	if reply41, err41 := redisService.ZCard("test3"); reply41 != 2 || err41 != nil {
		t.Errorf("redis set zcard[key is not '' && key exist one more] fail")
	}
	redisService.Del("test3")
}

func Test_ZSet_ZRank(t *testing.T) {
	redisService := InitService()
	// key not exist
	if reply1, err1 := redisService.ZRank("test1", []byte("value1")); reply1 != -1 || err1 != nil {
		t.Errorf("redis set zrank[key not exist] fail")
	}
	// key exist && member is nil
	if err2 := redisService.ZAdd("test2", 5, []byte("200"), []byte("value1"), []byte("100"), []byte("value2")); err2 != nil {
		t.Errorf("redis set zrank[key exist && member is nil] add fail")
	}
	if reply21, err21 := redisService.ZRank("test2", nil); reply21 != -1 || err21 != nil {
		t.Errorf("redis set zrank[key exist && member is nil] fail")
	}
	redisService.Del("test2")
	// key exist && member not in Zset
	if err3 := redisService.ZAdd("test3", 5, []byte("200"), []byte("value1"), []byte("100"), []byte("value2")); err3 != nil {
		t.Errorf("redis set zrank[key exist && member is nil] add fail")
	}
	if reply31, err31 := redisService.ZRank("test3", []byte("value3")); reply31 != -1 || err31 != nil {
		t.Errorf("redis set zrank[key exist && member in Zset] fail")
	}
	redisService.Del("test3")
	// key exist && member in Zset
	if err4 := redisService.ZAdd("test4", 5, []byte("200"), []byte("value1"), []byte("100"), []byte("value2")); err4 != nil {
		t.Errorf("redis set zrank[key exist && member is nil] add fail")
	}
	if reply41, err41 := redisService.ZRank("test4", []byte("value1")); reply41 != 1 || err41 != nil {
		t.Errorf("redis set zrank[key exist && member in Zset] fail")
	}
	redisService.Del("test4")
}

func Test_ZSet_ZRevRank(t *testing.T) {
	redisService := InitService()
	// key not exist
	if reply1, err1 := redisService.ZRevRank("test1", []byte("value1")); reply1 != -1 || err1 != nil {
		t.Errorf("redis set zrevrank[key not exist] fail")
	}
	// key exist && member is nil
	if err2 := redisService.ZAdd("test2", 5, []byte("200"), []byte("value1"), []byte("100"), []byte("value2")); err2 != nil {
		t.Errorf("redis set zrevrank[key exist && member is nil] add fail")
	}
	if reply21, err21 := redisService.ZRevRank("test2", nil); reply21 != -1 || err21 != nil {
		t.Errorf("redis set zrevrank[key exist && member is nil] fail")
	}
	redisService.Del("test2")
	// key exist && member not in Zset
	if err3 := redisService.ZAdd("test3", 5, []byte("200"), []byte("value1"), []byte("100"), []byte("value2")); err3 != nil {
		t.Errorf("redis set zrevrank[key exist && member is nil] add fail")
	}
	if reply31, err31 := redisService.ZRevRank("test3", []byte("value3")); reply31 != -1 || err31 != nil {
		t.Errorf("redis set zrevrank[key exist && member in Zset] fail")
	}
	redisService.Del("test3")
	// key exist && member in Zset
	if err4 := redisService.ZAdd("test4", 5, []byte("200"), []byte("value1"), []byte("100"), []byte("value2")); err4 != nil {
		t.Errorf("redis set zrevrank[key exist && member is nil] add fail")
	}
	if reply41, err41 := redisService.ZRevRank("test4", []byte("value2")); reply41 != 1 || err41 != nil {
		t.Errorf("redis set zrevrank[key exist && member in Zset] fail")
	}
	redisService.Del("test4")
}

func Test_ZSet_ZRange(t *testing.T) {
	redisService := InitService()
	// member not exist
	if reply1, err1 := redisService.ZRange("test1", 0, 1, false); len(reply1) != 0 || err1 != nil {
		t.Errorf("redis set zrange[member not exist] fail")
	}
	// member exist && start > 0 && start < stop
	if err2 := redisService.ZAdd("test2", 5, []byte("200"), []byte("value1"), []byte("100"), []byte("value2")); err2 != nil {
		t.Errorf("redis set zrange[member exist && start > 0 && start < stop] add fail")
	}
	if reply21, err21 := redisService.ZRange("test2", 1, 4, false); len(reply21) != 1 || err21 != nil {
		t.Errorf("redis set zrange[member exist && start > 0 && start < stop] fail")
	}
	redisService.Del("test2")
	// member exist && start > 0 && start == stop
	if err3 := redisService.ZAdd("test3", 5, []byte("200"), []byte("value1"), []byte("100"), []byte("value2")); err3 != nil {
		t.Errorf("redis set zrange[member exist && start > 0 && start == stop] add fail")
	}
	if reply31, err31 := redisService.ZRange("test3", 1, 1, false); len(reply31) != 1 || err31 != nil {
		t.Errorf("redis set zrange[member exist && start > 0 && start == stop] fail")
	}
	redisService.Del("test3")
	// member exist && start > 0 && start > stop
	if err4 := redisService.ZAdd("test4", 5, []byte("200"), []byte("value1"), []byte("100"), []byte("value2")); err4 != nil {
		t.Errorf("redis set zrange[member exist && start > 0 && start > stop] add fail")
	}
	if reply41, err41 := redisService.ZRange("test4", 1, 0, false); len(reply41) != 0 || err41 != nil {
		t.Errorf("redis set zrange[member exist && start > 0 && start > stop] fail")
	}
	redisService.Del("test4")
	// member exist && start == -1 && stop > 0
	if err5 := redisService.ZAdd("test5", 5, []byte("200"), []byte("value1"), []byte("100"), []byte("value2")); err5 != nil {
		t.Errorf("redis set zrange[member exist && start == -1 && stop > 0] add fail")
	}
	if reply51, err51 := redisService.ZRange("test5", -1, 1, false); len(reply51) != 1 || err51 != nil {
		t.Errorf("redis set zrange[member exist && start == -1 && stop > 0] fail")
	}

	redisService.Del("test5")
	// member exist && start == -1 && stop < -1
	if err6 := redisService.ZAdd("test6", 5, []byte("200"), []byte("value1"), []byte("100"), []byte("value2")); err6 != nil {
		t.Errorf("redis set zrange[member exist && start == -1 && stop > 0] add fail")
	}
	if reply61, err61 := redisService.ZRange("test6", -1, -3, false); len(reply61) != 0 || err61 != nil {
		t.Errorf("redis set zrange[member exist && start == -1 && stop > 0] fail")
	}
	if reply61, err61 := redisService.ZRange("test6", -2, -1, false); len(reply61) != 2 || err61 != nil {
		t.Errorf("redis set zrange[member exist && start == -1 && stop > 0] fail")
	}
	redisService.Del("test6")
	// member exist && start == 0
	if err7 := redisService.ZAdd("test7", 5, []byte("200"), []byte("value1"), []byte("100"), []byte("value2")); err7 != nil {
		t.Errorf("redis set zrange[member exist && start == -1 && stop > 0] add fail")
	}
	if reply71, err71 := redisService.ZRange("test7", 0, 1, false); len(reply71) != 2 || err71 != nil {
		t.Errorf("redis set zrange[member exist && start == -1 && stop > 0] fail")
	}
	redisService.Del("test7")
	// member exist && start ==0 && stop == -1
	if err8 := redisService.ZAdd("test8", 5, []byte("200"), []byte("value1"), []byte("100"), []byte("value2")); err8 != nil {
		t.Errorf("redis set zrange[member exist && start ==0 && stop == -1] add fail")
	}
	if reply81, err81 := redisService.ZRange("test8", 0, -1, false); len(reply81) != 2 || err81 != nil {
		t.Errorf("redis set zrange[member exist && start ==0 && stop == -1] fail")
	}
	redisService.Del("test8")
}

func Test_ZSet_ZRangeByScore(t *testing.T) {
	redisService := InitService()

	// test1
	if err := redisService.ZAdd("test1", 5, []byte("200"), []byte("value1"), []byte("100"), []byte("value2")); err != nil {
		t.Errorf("redis set ZRangeByScore[0, 100] add fail")
	}
	if ret, err := redisService.ZRangeByScore("test1", 0, 100, false); err != nil {
		t.Errorf("redis set ZRangeByScore[0, 100] get fail")
	} else if len(ret) != 1 {
		t.Errorf("redis set ZRangeByScore[0, 100] get fail")
	}
	redisService.Del("test1")

	// test1
	redisService.ZAdd("test2", 5, []byte("100.1234"), []byte("v1"))
	redisService.ZAdd("test2", 5, Float64ToByte(100.1234), []byte("v2"))
	redisService.ZAdd("test2", 5, Float64ToByte(200), []byte("v3"))

	ret, err := redisService.ZRange("test2", 0, 0, true)
	if err != nil {
		t.Errorf("redis set ZRange[0, 1] get fail")
	}
	if len(ret) != 2 {
		t.Errorf("redis set ZRange[0, 1] get fail ")
	}
	if string(ret[0]) != "v1" {
		t.Errorf("redis set ZRange[0, 1] get fail ")
	}

	if string(ret[1]) != "100.1234" {
		t.Errorf("redis set ZRange[0, 1] get fail")
	}

	if ret, err := redisService.ZRangeByScore("test2", 100.1234, 100.1234, false); err != nil {
		t.Errorf("redis set ZRangeByScore[0, 100] get fail")
	} else if len(ret) != 2 {
		t.Errorf("redis set ZRangeByScore[0, 100] get fail")
	} else if string(ret[0]) != "v1" {
		t.Errorf("redis set ZRangeByScore[0, 1] get fail ")
	} else if string(ret[1]) != "v2" {
		t.Errorf("redis set ZRangeByScore[0, 1] get fail ")
	}
	redisService.Del("test2")
}

func Test_ZSet_ZRevRange(t *testing.T) {
	redisService := InitService()
	// member not exist
	if reply1, err1 := redisService.ZRevRange("test1", 0, 1, false); len(reply1) != 0 || err1 != nil {
		t.Errorf("redis set zrevrange[member not exist] fail")
	}
	// member exist && start > 0 && start < stop
	if err2 := redisService.ZAdd("test2", 5, []byte("200"), []byte("value1"), []byte("100"), []byte("value2")); err2 != nil {
		t.Errorf("redis set zrevrange[member exist && start > 0 && start < stop] add fail")
	}
	if reply21, err21 := redisService.ZRevRange("test2", 1, 4, false); len(reply21) != 1 || err21 != nil {
		t.Errorf("redis set zrevrange[member exist && start > 0 && start < stop] fail")
	}
	redisService.Del("test2")
	// member exist && start > 0 && start == stop
	if err3 := redisService.ZAdd("test3", 5, []byte("200"), []byte("value1"), []byte("100"), []byte("value2")); err3 != nil {
		t.Errorf("redis set zrevrange[member exist && start > 0 && start == stop] add fail")
	}
	if reply31, err31 := redisService.ZRevRange("test3", 1, 1, false); len(reply31) != 1 || err31 != nil {
		t.Errorf("redis set zrevrange[member exist && start > 0 && start == stop] fail")
	}
	redisService.Del("test3")
	// member exist && start > 0 && start > stop
	if err4 := redisService.ZAdd("test4", 5, []byte("200"), []byte("value1"), []byte("100"), []byte("value2")); err4 != nil {
		t.Errorf("redis set zrevrange[member exist && start > 0 && start > stop] add fail")
	}
	if reply41, err41 := redisService.ZRevRange("test4", 1, 0, false); len(reply41) != 0 || err41 != nil {
		t.Errorf("redis set zrevrange[member exist && start > 0 && start > stop] fail")
	}
	redisService.Del("test4")
	// member exist && start == -1 && stop > 0
	if err5 := redisService.ZAdd("test5", 5, []byte("200"), []byte("value1"), []byte("100"), []byte("value2")); err5 != nil {
		t.Errorf("redis set zrevrange[member exist && start == -1 && stop > 0] add fail")
	}
	if reply51, err51 := redisService.ZRevRange("test5", -1, 1, false); len(reply51) != 1 || err51 != nil {
		t.Errorf("redis set zrevrange[member exist && start == -1 && stop > 0] fail")
	}

	redisService.Del("test5")
	// member exist && start == -1 && stop < -1
	if err6 := redisService.ZAdd("test6", 5, []byte("200"), []byte("value1"), []byte("100"), []byte("value2")); err6 != nil {
		t.Errorf("redis set zrevrange[member exist && start == -1 && stop > 0] add fail")
	}
	if reply61, err61 := redisService.ZRevRange("test6", -1, -3, false); len(reply61) != 0 || err61 != nil {
		t.Errorf("redis set zrevrange[member exist && start == -1 && stop > 0] fail")
	}
	if reply61, err61 := redisService.ZRevRange("test6", -2, -1, false); len(reply61) != 2 || err61 != nil {
		t.Errorf("redis set zrevrange[member exist && start == -1 && stop > 0] fail")
	}
	redisService.Del("test6")
	// member exist && start == 0
	if err7 := redisService.ZAdd("test7", 5, []byte("200"), []byte("value1"), []byte("100"), []byte("value2")); err7 != nil {
		t.Errorf("redis set zrevrange[member exist && start == -1 && stop > 0] add fail")
	}
	if reply71, err71 := redisService.ZRevRange("test7", 0, 1, false); len(reply71) != 2 || err71 != nil {
		t.Errorf("redis set zrevrange[member exist && start == -1 && stop > 0] fail")
	}
	redisService.Del("test7")
	// member exist && start ==0 && stop == -1
	if err8 := redisService.ZAdd("test8", 5, []byte("200"), []byte("value1"), []byte("100"), []byte("value2")); err8 != nil {
		t.Errorf("redis set zrevrange[member exist && start ==0 && stop == -1] add fail")
	}
	if reply81, err81 := redisService.ZRevRange("test8", 0, -1, false); len(reply81) != 2 || err81 != nil {
		t.Errorf("redis set zrevrange[member exist && start ==0 && stop == -1] fail")
	}
	redisService.Del("test8")
}

func Test_Hash_HSet(t *testing.T) {
	// redisService := InitService()
	// TODO
}

func Test_Hash_HGet(t *testing.T) {
	// redisService := InitService()
	// TODO
}

func Test_Hash_HMSet(t *testing.T) {
	// redisService := InitService()
	// TODO
}

func Test_Hash_HMGet(t *testing.T) {
	// redisService := InitService()
	// TODO
}

func Test_Hash_HDel(t *testing.T) {
	// redisService := InitService()
	// TODO
}

func Test_Hash_HExists(t *testing.T) {
	// redisService := InitService()
	// TODO
}

func Test_Hash_HKeys(t *testing.T) {
	// redisService := InitService()
	// TODO
}

func Test_Hash_HVals(t *testing.T) {
	// redisService := InitService()
	// TODO
}

func Test_Hash_HGetAll(t *testing.T) {
	// redisService := InitService()
	// TODO
}

func Test_Hash_HLen(t *testing.T) {
	// redisService := InitService()
	// TODO
}

func Test_Hash_Rpush(t *testing.T) {
	// redisService := InitService()
	// TODO
}

func Test_Hash_LLpush(t *testing.T) {
	// redisService := InitService()
	// TODO
}

func Test_Hash_LRpop(t *testing.T) {
	// redisService := InitService()
	// TODO
}

func Test_Hash_LLpop(t *testing.T) {
	// redisService := InitService()
	// TODO
}

func Test_Hash_LIndex(t *testing.T) {
	// redisService := InitService()
	// TODO
}
