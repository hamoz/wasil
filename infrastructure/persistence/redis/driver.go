package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/hamoz/wasil/commons/types"
	entity "github.com/hamoz/wasil/domain/entity"
	"github.com/hamoz/wasil/domain/repository"

	"github.com/go-redis/redis/v8"
)

const DRIVER_FILED string = "driver"
const DRIVERS_COLLECTION = "Drivers"

type DriverRepo struct {
	rds *redis.Client
}

var _ repository.DriverRepository = &DriverRepo{}

func NewDriverRepository(rds *redis.Client) *DriverRepo {
	return &DriverRepo{rds: rds}
}

func (repo *DriverRepo) Create(driver *entity.Driver) (*entity.Driver, error) {
	ctx := context.Background()
	key := driver.Channel + ":" + driver.ChannelId
	if _, err := repo.rds.HSet(ctx, key, DRIVER_FILED, driver).Result(); err != nil {
		return nil, fmt.Errorf("create: redis error: %w", err)
	}
	if driver.Registered {
		repo.rds.Persist(ctx, key)
	} else {
		repo.rds.Expire(ctx, key, time.Duration(5)*time.Minute)
	}
	//a.rds.Expire(ctx, key, time.Minute)

	return driver, nil
}

func (repo *DriverRepo) Find(key string) (*entity.Driver, error) {
	ctx := context.Background()
	result, err := repo.rds.HGet(ctx, key, DRIVER_FILED).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("find: redis error: %w", err)
	}

	driver := new(entity.Driver)
	if err := driver.UnmarshalBinary([]byte(result)); err != nil {
		return nil, fmt.Errorf("find: unmarshal error: %w", err)
	}
	return driver, nil
}

func (repo *DriverRepo) Update(driver *entity.Driver) error {
	// Find token:     a.rds.HGet()
	// Override token: a.rds.HSet()
	return nil
}

func (repo *DriverRepo) Delete(tokenID string) error {
	// Find token:   a.rds.HGet()
	// Delete token: a.rds.Del()
	return nil
}

func (a *DriverRepo) UpdateLocation(lng, lat float64, id string) {
	ctx := context.Background()
	a.rds.GeoAdd(ctx, DRIVERS_COLLECTION,
		&redis.GeoLocation{Longitude: lng, Latitude: lat, Name: id},
	)
}

func (a *DriverRepo) RemoveLocation(id string) {
	ctx := context.Background()
	a.rds.ZRem(ctx, DRIVERS_COLLECTION, id)
}

func (a *DriverRepo) GetNearby(limit int, lat, lng, r float64) []types.Location {
	/*
		WITHDIST: Also return the distance of the returned items from the
		specified center. The distance is returned in the same unit as the unit
		specified as the radius argument of the command.
		WITHCOORD: Also return the longitude,latitude coordinates of the matching items.
		WITHHASH: Also return the raw geohash-encoded sorted set score of the item,
		in the form of a 52 bit unsigned integer. This is only useful for low level
		hacks or debugging and is otherwise of little interest for the general user.
	*/
	ctx := context.Background()
	res, _ := a.rds.GeoRadius(ctx, DRIVERS_COLLECTION, lng, lat, &redis.GeoRadiusQuery{
		Radius:      r,
		Unit:        "km",
		WithGeoHash: true,
		WithCoord:   true,
		WithDist:    true,
		Count:       limit,
		Sort:        "ASC",
	}).Result()
	slice := make([]types.Location, len(res))
	for i, v := range res {
		slice[i] = types.Location{
			Name: v.Name,
			Lat:  float32(v.Latitude),
			Lng:  float32(v.Longitude),
			Dist: float32(v.Dist),
		}
	}
	return slice
}
