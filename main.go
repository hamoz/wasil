package main

import (
	"github.com/hamoz/wasil/commons/logs"
	"github.com/hamoz/wasil/commons/types"
	"github.com/hamoz/wasil/infrastructure/persistence/redis"
	"github.com/hamoz/wasil/interfaces/channels/telegram"
	log "github.com/sirupsen/logrus"
)

func main() {
	logs.Init()
	log.Info("init")
	repos, err := redis.NewRepositories("localhost:6379", "")
	if err != nil {
		log.Error(err)
		return
	}
	params := make(map[string]string)
	params["token"] = "1811929227:AAGzPvGlagBa8VQBG_A__5kK9VDFU6Q5XBw"
	setting := types.Settings{
		Name:             "PassengerBot",
		ChannelType:      types.Telegram,
		ChannelID:        "PassengerBot",
		OutQueue:         string(types.Telegram) + "_" + "PassengerBot" + "_out",
		InQueue:          string(types.Telegram) + "_" + "PassengerBot" + "_in",
		AdditionalParams: params,
	}
	log.Info("init")
	broker := redis.NewMessageBroker(repos.Db)
	passengerTelegramBot, err := telegram.NewTelegramBot(setting, broker)
	if err != nil {
		log.Error(err)
		return
	}
	passengerTelegramBot.Start()
	log.Info("bot started")
}
