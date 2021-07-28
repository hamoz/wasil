package redis

import (
	"context"

	"github.com/hamoz/wasil/domain/repository"
	log "github.com/sirupsen/logrus"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	*redis.Client
}

type Repositories struct {
	Passenger repository.PassengerRepository
	Driver    repository.DriverRepository
	Request   repository.RequestRepository
	Db        *redis.Client
}

func GetRedisClient() *redis.Client {
	//once.Do(func() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	//redisClient = &RedisClient{client}
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Could not connect to redis %v", err)
	}
	//})

	return client
}

func NewRepositories(addr, password string) (*Repositories, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       0,        // use default DB
	})

	//redisClient = &RedisClient{client}
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Could not connect to redis %v", err)
	}
	//})

	return &Repositories{
		Passenger: NewPassengerRepository(client),
		Driver:    NewDriverRepository(client),
		Request:   NewRequestRepository(client),
		Db:        client,
	}, nil
}

//closes the  database connection
func (r *Repositories) Close() error {
	return r.Db.Close()
}
