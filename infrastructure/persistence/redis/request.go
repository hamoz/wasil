package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/hamoz/wasil/commons/types"
	entity "github.com/hamoz/wasil/domain/entity"
	"github.com/hamoz/wasil/domain/repository"

	redis "github.com/go-redis/redis/v8"
)

const REQUEST_KEY string = "request"
const REQUESTS_COLLECTION = "Requests"

type RequestRepo struct {
	rds *redis.Client
}

var _ repository.RequestRepository = &RequestRepo{}

func NewRequestRepository(rds *redis.Client) *RequestRepo {
	return &RequestRepo{rds: rds}
}

func (a *RequestRepo) Create(request *entity.Request) (*entity.Request, error) {
	//ctx := context.Background()
	request.ID = REQUESTS_COLLECTION + ":" + string(request.ChannelType) + ":" + request.ChannelUserID
	return a.Update(request)
}

func (a *RequestRepo) Find(id string) (*entity.Request, error) {
	ctx := context.Background()
	result, err := a.rds.HGet(ctx, id, REQUEST_KEY).Result()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("find: redis error: %w", err)
	}
	if result == "" {
		return nil, nil
	}

	request := new(entity.Request)
	if err := request.UnmarshalBinary([]byte(result)); err != nil {
		return nil, fmt.Errorf("find: unmarshal error: %w", err)
	}
	return request, nil
}

func (a *RequestRepo) FindByChannelUser(channelType types.ChannelType, channelUserID string) (*entity.Request, error) {
	id := REQUESTS_COLLECTION + ":" + string(channelType) + ":" + channelUserID
	return a.Find(id)
}

func (a *RequestRepo) Update(request *entity.Request) (*entity.Request, error) {
	// Find token:     a.rds.HGet()
	// Override token: a.rds.HSet()
	ctx := context.Background()
	if _, err := a.rds.HSet(ctx, request.ID, REQUEST_KEY, request).Result(); err != nil {
		return nil, fmt.Errorf("create: redis error: %w", err)
	}
	if request.Status == types.Placed {
		a.rds.Persist(ctx, request.ID)
	} else {
		a.rds.Expire(ctx, request.ID, time.Duration(5)*time.Minute)
	}
	//a.rds.Expire(ctx, key, time.Minute)

	return request, nil
}

func (a *RequestRepo) Delete(id string) error {
	// Find token:   a.rds.HGet()
	// Delete token: a.rds.Del()
	return nil
}
