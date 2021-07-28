package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/hamoz/wasil/commons/types"
	log "github.com/sirupsen/logrus"
)

type MessageBroker struct {
	rClient *redis.Client
}

func NewMessageBroker(rClient *redis.Client) *MessageBroker {
	return &MessageBroker{
		rClient: rClient,
	}
}

func (p *MessageBroker) Subscribe(handler func(message *types.Message), channels ...string) {
	log.Infoln("subscribing")
	ctx := context.Background()
	rPubSub := p.rClient.Subscribe(ctx, channels...)
	_, err := rPubSub.Receive(ctx)
	if err != nil {
		log.Fatalln(err)
		panic(err)
	}
	log.Infoln("Subscribed successfully to channels")
	ch := rPubSub.Channel()
	for msg := range ch {
		message := new(types.Message)
		message.UnmarshalBinary([]byte(msg.Payload))
		handler(message)
	}

}

func (p *MessageBroker) Publish(channel string, message *types.Message) error {
	ctx := context.Background()
	_, err := p.rClient.Publish(ctx, channel, message).Result()
	return err
}
