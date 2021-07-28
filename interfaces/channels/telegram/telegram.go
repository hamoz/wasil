package telegram

import (
	"strconv"
	"time"

	"github.com/hamoz/wasil/commons/types"
	"github.com/hamoz/wasil/infrastructure/persistence/redis"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
)

type telegramBot struct {
	settings *types.Settings
	tb       *tb.Bot
	broker   *redis.MessageBroker
}

func NewTelegramBot(settings types.Settings, broker *redis.MessageBroker) (*telegramBot, error) {
	bot, err := tb.NewBot(tb.Settings{
		// You can also set custom API URL.
		// If field is empty it equals to "https://api.telegram.org".
		//URL: "http://195.129.111.17:8012",
		Token:  settings.AdditionalParams["token"],
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		return nil, err
	}
	b := &telegramBot{
		settings: &settings,
		tb:       bot,
		broker:   broker,
	}

	bot.Handle("/start", func(m *tb.Message) {
		b.incomingHandler(types.StartMessage, m)
	})

	bot.Handle(tb.OnText, func(m *tb.Message) {
		b.incomingHandler(types.TextMessage, m)
	})

	bot.Handle(tb.OnLocation, func(m *tb.Message) {
		b.incomingHandler(types.LocationMessage, m)
	})

	bot.Handle(tb.OnContact, func(m *tb.Message) {
		b.incomingHandler(types.ContactMessage, m)
	})

	return b, nil
}

func (b *telegramBot) Start() {
	//listent to message out going through OutQueue
	b.broker.Subscribe(b.outgoingHandler, b.settings.OutQueue)
	//b.tb.Start()

}
func (b *telegramBot) Stop() {
	b.tb.Stop()
}

func (b *telegramBot) Init() {
}

func (b *telegramBot) incomingHandler(msgType types.MessageType, m *tb.Message) {

	message := new(types.Message)
	message.ID = strconv.Itoa(m.ID)
	message.ChannelType = b.settings.ChannelType
	message.ChannelID = b.settings.ChannelID
	message.Type = msgType
	message.Sender = new(types.User)
	message.Sender.ID = strconv.FormatInt(m.Chat.ID, 10)
	message.Sender.ChannelType = message.ChannelType
	message.Sender.Name = m.Sender.FirstName + " " + m.Sender.LastName
	message.Destination = new(types.User)
	message.Destination.Name = b.settings.Name
	message.Destination.ChannelType = message.ChannelType
	if msgType == types.TextMessage {
		message.Text = m.Text
	}
	if msgType == types.ContactMessage {
		message.Contact = new(types.Contact)
		message.Contact.Phone = m.Contact.PhoneNumber
		message.Contact.Name = m.Contact.FirstName + " " + m.Contact.LastName
	}
	if msgType == types.LocationMessage {
		message.Loc = new(types.Location)
		message.Loc.Lng = m.Location.Lng
		message.Loc.Lat = m.Location.Lat
	}
	//glog.Info("Publising message to App")
	err := b.broker.Publish(b.settings.InQueue, message)
	if err != nil {
		log.Error(err)
	}
}

func (b *telegramBot) outgoingHandler(message *types.Message) {
	/*bArr, error := message.MarshalBinary()
	if error == nil {
		glog.Info(string(bArr))
	}*/
	userId, err := strconv.Atoi(message.Destination.ID)
	if err != nil {
		log.Error(err)
		return
	}
	switch message.Type {
	case types.TextMessage:
		b.tb.Send(&tb.User{ID: userId}, message.Text)
	case types.LocationMessage:
		loc := &tb.Location{
			Lat:        message.Loc.Lat,
			Lng:        message.Loc.Lng,
			LivePeriod: 0,
		}
		b.tb.Send(&tb.User{ID: userId}, message.Text, loc)
	case types.GetContactMessage:
		r := &tb.ReplyMarkup{ResizeReplyKeyboard: true}
		r.Reply(r.Row(r.Contact("أرسل")))
		b.tb.Send(&tb.User{ID: userId}, message.Text, r)
	}

}
