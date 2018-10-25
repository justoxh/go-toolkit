# config

## support

support format: toml, ymal, json, ini, properties and so on

## usage

`govendor fetch github.com/justoxh/go-toolkit/config`

## example

```go
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
		Redis redisOptions
	}

	var conf testConf

	testfile := "./testdata/app_test.toml"
	err := LoadConfig(testfile, &conf)
	if err != nil {
		// some deal
    }
    
    fmt.Println(conf.Host)
```
