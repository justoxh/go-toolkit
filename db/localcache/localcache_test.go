package localcache

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/naoina/toml"
)

type RelayConfig struct {
	Name string
	Role string

	Cache CacheOptions
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

func Test_Initialize(t *testing.T) {
	var config_relay1 *RelayConfig = &RelayConfig{}
	var config_file1 string = "./localcache.toml"
	loadConfig(config_file1, config_relay1)
	if config_relay1.Name != "miner" {
		t.Errorf("config item Name[%s] not equal miner", config_relay1.Name)
	}
	var cacheService *LocalCacheService = &LocalCacheService{}
	cacheService.Initialize(config_relay1.Cache)
	// test cache
	err := cacheService.cache.Add("test1", "value1", 5)
	if err != nil {
		t.Errorf("cache string add fail")
	}
	reply, flag := cacheService.cache.Get("test1")
	fmt.Println(reply, flag)
	if !flag || reply.(string) != "value1" {
		t.Errorf("cache string get fail")
	}
	cacheService.cache.Delete("test1")
}
