package redis

import (
	"bufio"
	"errors"
	"github.com/naoina/toml"
	"os"
	"testing"

	_ "github.com/justoxh/go-toolkit/log"
	"github.com/justoxh/go-toolkit/log/logruslogger"


)

type RelayConfigNew struct {
	Name string
	Role string

	Redis RedisOptions
}

func loadConfigNew(file string, cfg *RelayConfig) {
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

func InitServiceNew() *RedisCacheService {
	var config_relay1 *RelayConfig = &RelayConfig{}
	var config_file1 string = "./pool_test_nopass.toml"
	loadConfigNew(config_file1, config_relay1)
	var redisService *RedisCacheService = &RedisCacheService{}
	var conf = logruslogger.Options{}
	log := logruslogger.GetLoggerWithOptions("test", &conf)
	redisService.Initialize(config_relay1.Redis, log)
	return redisService
}

func Test_List_LLlen(t *testing.T) {
	var redisService *RedisCacheService = InitServiceNew()
	// test key nil, value not nil
	res1,err1 := redisService.LLlen("ethereum_accounts")
	if err1 != nil {
		t.Errorf("redis llen error 01")
	}

	if res1 != 0 {
		t.Errorf("redis llen error 02:%v",res1)
	}
	redisService.Del("ethereum_accounts")

	res2,err2 := redisService.LLlen("ethereum_accounts")
	if err2 != nil {
		t.Errorf("redis llen error 03")
	}

	if res2 != 0 {
		t.Errorf("redis llen error 04")
	}

	redisService.LLpush("ethereum_accounts",[]byte("123"))

	res3,err3 := redisService.LLlen("ethereum_accounts")
	if err3 != nil {
		t.Errorf("redis llen error 05")
	}

	if res3 != 1 {
		t.Errorf("redis llen error 06")
	}
}
