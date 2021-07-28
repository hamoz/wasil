package application

import (
	"github.com/hamoz/wasil/commons/types"
	"github.com/hamoz/wasil/domain/entity"
	"github.com/hamoz/wasil/domain/fsm"
	"github.com/hamoz/wasil/domain/service"
	redis "github.com/hamoz/wasil/infrastructure/persistence/redis"
	log "github.com/sirupsen/logrus"
)

type PassengerApp struct {
	passengerService service.PassengerService
	requestService   service.RequestService
	broker           *redis.MessageBroker
	InQueue          string
	OutQueue         string
}

type PassengerAppInterface interface {
	RegisterPassenger(*entity.Passenger) (*entity.Passenger, error)
	GetPassenger(id string) (*entity.Passenger, error)
	GetPassengers() ([]entity.Passenger, error)
	DeletePassenger(id string) error
	Start()
}

var ctx *fsm.RequestContext //= fsm.NewRequestContext(nil, nil, nil, nil)

func NewPassengerApp(passengerService service.PassengerService,
	requestService service.RequestService, broker *redis.MessageBroker) *PassengerApp {
	//msg := new(common.Message)
	return &PassengerApp{
		passengerService: passengerService,
		requestService:   requestService,
		broker:           broker,
		InQueue:          "Passenger_in",
		OutQueue:         "Passenger_out",
	}
}

func (app *PassengerApp) RegisterPassenger(passenger *entity.Passenger) (*entity.Passenger, error) {
	return app.passengerService.Register(passenger)
}

func (app *PassengerApp) GetPassenger(id string) (*entity.Passenger, error) {
	return app.passengerService.Find(id)
}

func (app *PassengerApp) GetPassengers() ([]entity.Passenger, error) {
	return nil, nil
}

func (app *PassengerApp) DeletePassenger(id string) error {
	return app.passengerService.Delete(id)
}

func (app *PassengerApp) Start() {
	app.broker.Subscribe(app.handle, "Telegram_PassengerBot_in")
}

func (app *PassengerApp) handle(message *types.Message) {
	responseMessage := composeResponseMessage(message)
	requestFsm := fsm.NewRequestFSM()
	ctx := fsm.NewRequestContext(message, responseMessage, app.requestService, app.passengerService)
	ctx.Fsm = requestFsm
	request, err := app.requestService.FindByChannelUser(message.ChannelType, message.Sender.ID)
	if err != nil {
		log.Errorf("error, cannot find request", err)
		responseMessage.Text = "حدث خطأ أثناء إجراء طلبك الرجاء المعاودة لاحقا أو الاتصال بخدمة العملاء"
	} else {
		if request == nil {
			requestFsm.Current = types.Default
		} else {
			requestFsm.Current = request.Status
		}

		switch message.Type {
		case types.StartMessage:
			err = requestFsm.SendEvent(fsm.Start, ctx)
		case types.LocationMessage:
			err = requestFsm.SendEvent(fsm.LocationReceived, ctx)
		case types.TextMessage:
			err = requestFsm.SendEvent(fsm.TextReceived, ctx)
		case types.ContactMessage:
			err = requestFsm.SendEvent(fsm.ContactReceived, ctx)
		default:
			err = requestFsm.SendEvent(fsm.Start, ctx)
		}
		if err != nil {
			log.Errorf("error, cannot update request", err)
			responseMessage.Text = "حدث خطأ أثناء إجراء طلبك الرجاء المعاودة لاحقا أو الاتصال بخدمة العملاء"
		}
	}
	//Sending response
	brokerChannel := string(responseMessage.ChannelType) + "_" + responseMessage.ChannelID + "_out"
	//glog.Info("publishing response")
	err = app.broker.Publish(brokerChannel, responseMessage)
	if err != nil {
		log.Error(err)
	}

}

func composeResponseMessage(message *types.Message) *types.Message {
	responseMessage := new(types.Message)
	responseMessage.Sender = message.Destination
	responseMessage.Destination = message.Sender
	responseMessage.ChannelID = message.ChannelID
	responseMessage.ChannelType = message.ChannelType
	responseMessage.Type = types.TextMessage
	return responseMessage
}
