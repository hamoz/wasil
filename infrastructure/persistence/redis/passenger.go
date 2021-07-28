package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/hamoz/wasil/commons/types"
	entity "github.com/hamoz/wasil/domain/entity"
	"github.com/hamoz/wasil/domain/repository"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

const PASSENGER_KEY string = "passenger"
const PASSENGERS_COLLECTION string = "Passengers"

type PassengerRepo struct {
	rds *redis.Client
}

var _ repository.PassengerRepository = &PassengerRepo{}

func NewPassengerRepository(rds *redis.Client) *PassengerRepo {
	return &PassengerRepo{rds: rds}
}

func (a *PassengerRepo) Create(passenger *entity.Passenger) (*entity.Passenger, error) {
	ctx := context.Background()
	passenger.ID = PASSENGERS_COLLECTION + ":" + passenger.ChannelType + ":" + passenger.ChannelUserID
	log.Infoln("PassengerID :" + passenger.ID)
	if _, err := a.rds.HSet(ctx, passenger.ID, PASSENGER_KEY, passenger).Result(); err != nil {
		return nil, fmt.Errorf("create: redis error: %w", err)
	}
	if !passenger.Registered {
		a.rds.Expire(ctx, passenger.ID, time.Duration(5)*time.Minute).Result()
	}
	return passenger, nil
}

func (a *PassengerRepo) Find(id string) (*entity.Passenger, error) {
	ctx := context.Background()
	//id = PASSENGERS_COLLECTION + ":" + id
	result, err := a.rds.HGet(ctx, id, PASSENGER_KEY).Result()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("find: redis error: %w", err)
	}
	if result == "" {
		return nil, nil
	}

	passenger := new(entity.Passenger)
	if err := passenger.UnmarshalBinary([]byte(result)); err != nil {
		return nil, fmt.Errorf("find: unmarshal error: %w", err)
	}

	return passenger, nil
}

func (a *PassengerRepo) FindByChannelUser(channelType types.ChannelType, channelUserID string) (*entity.Passenger, error) {

	ctx := context.Background()
	id := PASSENGERS_COLLECTION + ":" + string(channelType) + ":" + channelUserID
	result, err := a.rds.HGet(ctx, id, PASSENGER_KEY).Result()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("find: redis error: %w", err)
	}
	if result == "" {
		return nil, nil
	}

	passenger := new(entity.Passenger)
	if err := passenger.UnmarshalBinary([]byte(result)); err != nil {
		return nil, fmt.Errorf("find: unmarshal error: %w", err)
	}

	return passenger, nil
}

func (a PassengerRepo) Update(passenger *entity.Passenger) error {
	// Find token:     a.rds.HGet()
	// Override token: a.rds.HSet()
	return nil
}

func (a PassengerRepo) Delete(key string) error {
	// Find token:   a.rds.HGet()
	// Delete token: a.rds.Del()
	return nil
}
