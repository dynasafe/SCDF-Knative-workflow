package pkg

import (
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Master   string
	Host     string
	Password string
	Username string
	Database string
}

func GetRedisOptions(config RedisConfig) (*redis.FailoverOptions, error) {
	redisOptions := &redis.FailoverOptions{
		MasterName:    config.Master,
		SentinelAddrs: []string{config.Host},
	}

	// password
	if config.Password != "" {
		redisOptions.Password = config.Password
	}

	// database
	if config.Database != "" {
		db, err := strconv.Atoi(config.Database)
		if err != nil {
			return nil, fmt.Errorf("invalid database: %v", err)
		}
		redisOptions.DB = db
	}

	// username
	if config.Username != "" {
		redisOptions.Username = config.Username
	}

	return redisOptions, nil
}
