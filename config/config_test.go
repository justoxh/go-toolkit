package config

import (
	"testing"
)

func Test_LoadConfig(t *testing.T) {
	type mysqlOptions struct {
		Hostname           string
		Port               string
		User               string
		Password           string
		DBName             string
		TablePrefix        string
		MaxOpenConnections int
		MaxIdleConnections int
		ConnMaxLifetime    int
		Debug              bool
	}

	type redisOptions struct {
		Host        string
		Port        string
		Password    string
		IdleTimeout int
		MaxIdle     int
		MaxActive   int
	}

	type testConf struct {
		Name  string
		Mysql mysqlOptions
		Redis redisOptions
	}

	var conf testConf
	testfile := ""
	err := LoadConfig(testfile, &conf)
	if err != nil {
		t.Log(err)
	}

	testfile = "./testdata/app_test.toml"
	err = LoadConfig(testfile, &conf)
	if err != nil {
		t.Error(err)
	}
	t.Logf("conf content: %#v", conf)

	testfile = "./testdata/app_test_no_exist.toml"
	err = LoadConfig(testfile, &conf)
	if err != nil {
		t.Log(err)
	}

	testfile = "./testdata/app_test_bad.toml"
	var conf1 testConf
	err = LoadConfig(testfile, &conf1)
	if err != nil {
		t.Log(err)
	}

	t.Logf("conf content: %#v", conf)

}
