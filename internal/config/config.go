package config

import (
	"gopoker/internal/identity/db"
	"gopoker/internal/mediator/cache"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type DBType int

type yamlConfig struct {
	TableKeys *struct {
		Queue string `yaml: "queue"`
		Store string `yaml: "store"`
	} `yaml:"table_keys"`
	RestAPI *struct {
		Port uint32 `yaml:"port"`
		File *struct {
			Path string `yaml:"path"`
		} `yaml:"file"`
		Redis *struct {
			Username string `yaml:"username"`
			Password string `yaml:"password"`
			Db       int    `yaml:"db"`
		} `yaml:"redis"`
	} `yaml:"rest_api"`
	SocketAPI *struct {
		Port  uint32 `yaml:"port"`
		Redis *struct {
			Username string `yaml:"username"`
			Password string `yaml:"password"`
			Db       int    `yaml:"db"`
		} `yaml:"redis"`
	} `yaml:"socket_api"`
}

type Config struct {
	tableStore string
	tableQueue string

	socketPort  int
	restAPIPort int

	db    db.IIDB
	cache cache.ICache
}

func (c Config) GetTableKeys() (tableStore string, tableQueue string) {
	return c.tableStore, c.tableQueue
}

func (c Config) GetCache() cache.ICache {
	return c.cache
}

func (c Config) GetDB() db.IIDB {
	return c.db
}

func ParseFromPath(path string) *Config {
	content, err := ioutil.ReadFile(path)

	if err != nil {
		panic("Could not parse config file from path!")
	}

	return Parse(string(content))
}

func Parse(content string) *Config {
	var input yamlConfig
	config := &Config{}

	if err := yaml.Unmarshal([]byte(content), &input); err != nil {
		panic(err)
	}

	if input.TableKeys == nil {
		panic("table_keys is not set in Config..")
	}

	if input.TableKeys.Store == input.TableKeys.Queue {
		panic("Queue and Store cannot share the same key.")
	}

	if input.RestAPI == nil {
		panic("REST API parameters not set..")
	}

	if input.RestAPI.Port > 0xFFFF {
		panic("Invalid REST API port number")
	}

	if input.SocketAPI == nil {
		panic("Socket API parameters not set")
	}

	if input.SocketAPI.Port > 0xFFFF {
		panic("Invalid REST API port number")
	}

	if input.SocketAPI.Port == input.RestAPI.Port {
		panic("Socket API and REST API cannot share a port")
	}
	if input.RestAPI.Redis != nil {
		panic("Redis Identity DB not set up in go yet.")
	} else if input.RestAPI.File != nil {
		config.db = db.NewFileDB(input.RestAPI.File.Path)
	} else {
		panic("Invalid Identity DB config set in for REST API")
	}

	if input.SocketAPI.Redis != nil {
		config.cache = cache.NewRedisCache(input.SocketAPI.Redis.Username, input.SocketAPI.Redis.Password, input.SocketAPI.Redis.Db, config.tableQueue, config.tableStore)
	} else {
		panic("Invalid Cache config set for Socket API")
	}
	config.tableStore = input.TableKeys.Store
	config.tableQueue = input.TableKeys.Queue

	config.socketPort = int(input.SocketAPI.Port)
	config.restAPIPort = int(input.RestAPI.Port)

	return config
}
