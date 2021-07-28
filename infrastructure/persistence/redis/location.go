package redis

import (
	"context"

	"github.com/hamoz/wasil/commons/types"
	"github.com/hamoz/wasil/domain/repository"

	"github.com/go-redis/redis/v8"
)

type LocationRepo struct {
	rds *redis.Client
}

var _ repository.LocationRepository = &LocationRepo{}

func NewLocationRepository(rds *redis.Client) *LocationRepo {
	return &LocationRepo{rds}
}

func (repo *LocationRepo) Add(location *types.Location) error {
	ctx := context.Background()
	return repo.rds.GeoAdd(ctx, DRIVERS_COLLECTION,
		&redis.GeoLocation{
			Name:      location.Name,
			Longitude: float64(location.Lng),
			Latitude:  float64(location.Lat),
			Dist:      0,
			GeoHash:   0,
		},
	).Err()
}

func (repo *LocationRepo) Find(id string) (*types.Location, error) {
	//TO DO
	//implement Find
	return nil, nil
}
func (repo *LocationRepo) Update(driver *types.Location) error {
	//TO DO
	//implement update
	return nil
}

func (repo *LocationRepo) Delete(id string) error {
	ctx := context.Background()
	return repo.rds.ZRem(ctx, DRIVERS_COLLECTION, id).Err()

}
func (repo *LocationRepo) FindByDistance(center *types.Location, count int, distance float64) []types.Location {
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
	res, _ := repo.rds.GeoRadius(ctx, DRIVERS_COLLECTION, float64(center.Lng), float64(center.Lat), &redis.GeoRadiusQuery{
		Radius:      distance,
		Unit:        "km",
		WithGeoHash: true,
		WithCoord:   true,
		WithDist:    true,
		Count:       count,
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
