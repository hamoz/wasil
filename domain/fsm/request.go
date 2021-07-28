package fsm

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hamoz/wasil/commons/types"
	"github.com/hamoz/wasil/domain/entity"
	"github.com/hamoz/wasil/domain/service"

	log "github.com/sirupsen/logrus"
)

const (
	LocationReceived EventType = "LocationReceived"
	ContactReceived  EventType = "ContactReceived"
	TextReceived     EventType = "TextReceived"
	Failed           EventType = "FailedEvent"
	Reset            EventType = "ResetEvent"
	Start            EventType = "StartEvent"
	InvalidInput     EventType = "InvalidInput"
)

// StartAction represents the action executed on the first time intering bot entering
type StartAction struct{}

// OrigLocationAction represents the action executed on entering the WaitingOrigLocation state.
type UpdateOrigLocationAction struct{}

// DestLocationAction represents the action executed on entering the WaitingDestLocation state.
//type DestLocationAction struct{}

// WaitingContactAction represents the action executed on entering the WaitingContact state.
type UpdateDestLocationAction struct{}

// WaitingConfirmAction represents the action executed on entering the WaitingConfirm state.
type UpdateContactAction struct{}

type ConfirmAction struct{}

type CancellAction struct{}

type RequestContext struct {
	requestMessage   *types.Message
	responseMessage  *types.Message
	requestService   service.RequestService
	passengerService service.PassengerService
	Fsm              *StateMachine
}

func (ctx *RequestContext) getFsm() *StateMachine {
	return ctx.Fsm
}

var fsmEntry log.Entry

func NewRequestContext(requestMessage *types.Message, responseMessage *types.Message,
	requestService service.RequestService,
	passengerService service.PassengerService) *RequestContext {
	return &RequestContext{
		requestMessage:   requestMessage,
		responseMessage:  responseMessage,
		requestService:   requestService,
		passengerService: passengerService,
	}

}

func (a *StartAction) Execute(eventCtx EventContext) EventType {
	ctx := eventCtx.(*RequestContext)
	logEntry := log.WithFields(
		log.Fields{
			"channel":       ctx.requestMessage.ChannelType,
			"from":          ctx.requestMessage.Sender.ID,
			"to":            ctx.requestMessage.Destination.Name,
			"currentState":  ctx.Fsm.Current,
			"previousState": ctx.Fsm.Previous,
			"event":         ctx.Fsm.CurrentEvent,
			"inputType":     ctx.requestMessage.Type,
			"input":         ctx.requestMessage.GetPayload(),
		},
	)
	logEntry.Infoln("StartAction")
	ctx.responseMessage.Text = "لطلب رحلتك رجاءً أرسل موقعك الحالي"
	return NoOp
}

func (a *UpdateOrigLocationAction) Execute(eventCtx EventContext) EventType {
	ctx := eventCtx.(*RequestContext)
	logEntry := log.WithFields(
		log.Fields{
			"channel":       ctx.requestMessage.ChannelType,
			"from":          ctx.requestMessage.Sender.ID,
			"to":            ctx.requestMessage.Destination.Name,
			"currentState":  ctx.Fsm.Current,
			"previousState": ctx.Fsm.Previous,
			"event":         ctx.Fsm.CurrentEvent,
			"inputType":     ctx.requestMessage.Type,
			"input":         ctx.requestMessage.GetPayload(),
		},
	)
	logEntry.Infoln("UpdateOrigLocationAction")
	//Input validation
	if ctx.Fsm.CurrentEvent != LocationReceived {
		ctx.responseMessage.Text = "أرسل وجهة الرحلة"
		return NoOp
	}

	request, err := ctx.requestService.FindByChannelUser(ctx.requestMessage.ChannelType,
		ctx.requestMessage.Sender.ID)
	if err != nil {
		logEntry.Errorf("error, cannot find request", err)
		ctx.requestMessage.Text = "حدث خطأ أثناء إجراء طلبك الرجاء المعاودة لاحقا أو الاتصال بخدمة العملاء"
		return Failed
	}
	if request == nil {
		logEntry.Debugln("Request not found, creating new one")
		request = new(entity.Request)
		//request.ID = requestID
		request.Time = time.Now()
		request.Status = types.Default
		request.ChannelType = ctx.requestMessage.ChannelType
		request.ChannelUserID = ctx.requestMessage.Sender.ID
	}
	request.FromLoc = ctx.requestMessage.Loc
	request.Status = types.OriginLocationUpdated
	if request.ID == "" {
		_, err = ctx.requestService.Create(request)
	} else {
		_, err = ctx.requestService.Update(request)
	}
	if err != nil {
		logEntry.Errorln("error, cannot update request", err)
		ctx.responseMessage.Text = "حدث خطأ أثناء إجراء طلبك الرجاء المعاودة لاحقا أو الاتصال بخدمة العملاء"
		return Failed
	}
	ctx.responseMessage.Text = "أرسل وجهة الرحلة"
	return NoOp
}

func (a *UpdateDestLocationAction) Execute(eventCtx EventContext) EventType {
	ctx := eventCtx.(*RequestContext)
	logEntry := log.WithFields(
		log.Fields{
			"channel":       ctx.requestMessage.ChannelType,
			"from":          ctx.requestMessage.Sender.ID,
			"to":            ctx.requestMessage.Destination.Name,
			"currentState":  ctx.Fsm.Current,
			"previousState": ctx.Fsm.Previous,
			"event":         ctx.Fsm.CurrentEvent,
			"inputType":     ctx.requestMessage.Type,
			"input":         ctx.requestMessage.GetPayload(),
		},
	)
	logEntry.Infoln("UpdateDestLocationAction")
	//Input validation
	if ctx.Fsm.CurrentEvent != LocationReceived {
		ctx.responseMessage.Text = "لإكمال الطلب الرجاء مشاركة هاتفك بالضغط على الزر أدناه (يطلب منك عند أول مرة فقط)"
		return NoOp
	}

	request, err := ctx.requestService.FindByChannelUser(ctx.requestMessage.ChannelType, ctx.requestMessage.Sender.ID)
	//request == nil, this  shouldn't happen
	if err != nil || request == nil {
		logEntry.Errorln("error, cannot find request", err)
		ctx.responseMessage.Text = "حدث خطأ أثناء إجراء طلبك الرجاء المعاودة لاحقا أو الاتصال بخدمة العملاء"
		return Failed
	}

	request.ToLocation = ctx.requestMessage.Loc
	request.Status = types.DestLocationUpdated
	if request.ID == "" {
		_, err = ctx.requestService.Create(request)
	} else {
		_, err = ctx.requestService.Update(request)
	}
	if err != nil {
		logEntry.Errorln("error, cannot update request", err)
		ctx.responseMessage.Text = "حدث خطأ أثناء إجراء طلبك الرجاء المعاودة لاحقا أو الاتصال بخدمة العملاء"
		return Failed
	}
	passenger, err := ctx.passengerService.FindByChannelUser(ctx.requestMessage.ChannelType, ctx.requestMessage.Sender.ID)
	if err != nil {
		logEntry.Errorln("error, cannot find passenger", err)
		return Failed
	}
	//
	if passenger == nil {
		ctx.responseMessage.Type = types.GetContactMessage
		ctx.responseMessage.Text = "لإكمال الطلب الرجاء مشاركة هاتفك بالضغط على الزر أدناه (يطلب منك عند أول مرة فقط)"
		return NoOp
	} else {
		ctx.requestMessage.Contact = &types.Contact{
			Phone: passenger.Phone,
			Name:  passenger.Name,
		}
		return ContactReceived
	}
}

func (a *UpdateContactAction) Execute(eventCtx EventContext) EventType {
	ctx := eventCtx.(*RequestContext)
	logEntry := log.WithFields(
		log.Fields{
			"channel":       ctx.requestMessage.ChannelType,
			"from":          ctx.requestMessage.Sender.ID,
			"to":            ctx.requestMessage.Destination.Name,
			"currentState":  ctx.Fsm.Current,
			"previousState": ctx.Fsm.Previous,
			"event":         ctx.Fsm.CurrentEvent,
			"inputType":     ctx.requestMessage.Type,
			"input":         ctx.requestMessage.GetPayload(),
		},
	)
	logEntry.Infoln("UpdateContactAction")
	//Input validation
	if ctx.Fsm.CurrentEvent != ContactReceived {
		ctx.responseMessage.Text = "لإكمال الطلب الرجاء مشاركة هاتفك بالضغط على الزر أدناه (يطلب منك عند أول مرة فقط)"
		return NoOp
	}

	passenger, err := ctx.passengerService.FindByChannelUser(ctx.requestMessage.ChannelType, ctx.requestMessage.Sender.ID)
	if err != nil {
		logEntry.Errorln("error, cannot find passenger", err)
		ctx.responseMessage.Text = "حدث خطأ أثناء إجراء طلبك الرجاء المعاودة لاحقا أو الاتصال بخدمة العملاء"
		return Failed
	}
	if passenger == nil {
		passenger = new(entity.Passenger)
		passenger.Name = ctx.requestMessage.Contact.Name
		passenger.Phone = ctx.requestMessage.Contact.Phone
		passenger.ChannelType = string(ctx.requestMessage.ChannelType)
		passenger.ChannelUserID = ctx.requestMessage.Sender.ID
		passenger.Registered = true
		passenger, err = ctx.passengerService.Register(passenger)
		if err != nil {
			logEntry.Errorln("error, cannot register passenger", err)
			ctx.responseMessage.Text = "حدث خطأ أثناء إجراء طلبك الرجاء المعاودة لاحقا أو الاتصال بخدمة العملاء"
			return Failed
		}
	}
	request, err := ctx.requestService.FindByChannelUser(ctx.requestMessage.ChannelType, ctx.requestMessage.Sender.ID)
	//request == nil this shouldn't happen
	if err != nil {
		logEntry.Errorln("error, cannot find request", err)
		ctx.responseMessage.Text = "حدث خطأ أثناء إجراء طلبك الرجاء المعاودة لاحقا أو الاتصال بخدمة العملاء"
		return Failed
	}
	request.Status = types.ContactUpdated
	request.PassengerID = passenger.ID
	cost := -10
	cost, err = ctx.requestService.Price(request)
	if err != nil {
		logEntry.Errorln("error, cannot calculate request price", err)
		ctx.responseMessage.Text = "حدث خطأ أثناء إجراء طلبك الرجاء المعاودة لاحقا أو الاتصال بخدمة العملاء"
		return Failed
	}
	ctx.responseMessage.Text = fmt.Sprintf("تكلفة الرحلة %d جنيه. أرسل الرقم : \n1.للتأكيد\n2.للإلغاء", cost)
	_, err = ctx.requestService.Update(request)
	if err != nil {
		logEntry.Errorln("error, cannot update request", err)
		ctx.responseMessage.Text = "حدث خطأ أثناء إجراء طلبك الرجاء المعاودة لاحقا أو الاتصال بخدمة العملاء"
		return Failed
	}
	return NoOp
}

func (a *ConfirmAction) Execute(eventCtx EventContext) EventType {
	ctx := eventCtx.(*RequestContext)
	logEntry := log.WithFields(
		log.Fields{
			"channel":       ctx.requestMessage.ChannelType,
			"from":          ctx.requestMessage.Sender.ID,
			"to":            ctx.requestMessage.Destination.Name,
			"currentState":  ctx.Fsm.Current,
			"previousState": ctx.Fsm.Previous,
			"event":         ctx.Fsm.CurrentEvent,
			"inputType":     ctx.requestMessage.Type,
			"input":         ctx.requestMessage.GetPayload(),
		},
	)
	logEntry.Infoln("ConfirmAction")
	//Input validation
	ctx.responseMessage.Text = strings.TrimSpace(ctx.requestMessage.Text)
	matched, err := regexp.MatchString(`^[12١٢]$`, ctx.responseMessage.Text)
	if ctx.Fsm.CurrentEvent != TextReceived || !matched {
		ctx.responseMessage.Text = "أرسل الرقم : \n1. للتأكيد\n2. للإلغاء"
		return NoOp
	}

	request, err := ctx.requestService.FindByChannelUser(ctx.requestMessage.ChannelType, ctx.requestMessage.Sender.ID)
	//request == nil this shouldn't happen
	if err != nil {
		logEntry.Errorln("error, cannot find request", err)
		ctx.responseMessage.Text = "حدث خطأ أثناء إجراء طلبك الرجاء المعاودة لاحقا أو الاتصال بخدمة العملاء"
		return Failed
	}
	if ctx.requestMessage.Text == "1" || ctx.requestMessage.Text == "١" {
		request.Status = types.Confirmed
		ctx.responseMessage.Text = "شكراً، تم تأكيد طلبك، الرجاء الانتطار قليلا يجري البحث عن أقرب سائق لموقعك"
	} else if ctx.requestMessage.Text == "2" || ctx.requestMessage.Text == "٢" {
		request.Status = types.Cancelled
		ctx.responseMessage.Text = "تم إلغاء الطلب"
	} else {
		ctx.responseMessage.Text = "أرسل الرقم : \n1. للتأكيد\n2. للإلغاء"
		ctx.responseMessage.Type = types.Error
		return Failed
	}
	_, err = ctx.requestService.Update(request)
	if err != nil {
		logEntry.Errorln("error, cannot find request", err)
		ctx.responseMessage.Text = "حدث خطأ أثناء إجراء طلبك الرجاء المعاودة لاحقا أو الاتصال بخدمة العملاء"
		return Failed
	}
	return NoOp
}

func (a *CancellAction) Execute(eventCtx EventContext) EventType {
	ctx := eventCtx.(*RequestContext)
	logEntry := log.WithFields(
		log.Fields{
			"channel":       ctx.requestMessage.ChannelType,
			"from":          ctx.requestMessage.Sender.ID,
			"to":            ctx.requestMessage.Destination.Name,
			"currentState":  ctx.Fsm.Current,
			"previousState": ctx.Fsm.Previous,
			"event":         ctx.Fsm.CurrentEvent,
			"inputType":     ctx.requestMessage.Type,
			"input":         ctx.requestMessage.GetPayload(),
		},
	)
	logEntry.Infoln("CancellAction")
	//Input validation
	ctx.responseMessage.Text = strings.TrimSpace(ctx.requestMessage.Text)
	matched, err := regexp.MatchString(`^[1١]$`, ctx.responseMessage.Text)
	if ctx.Fsm.CurrentEvent != TextReceived || !matched {
		ctx.responseMessage.Text = "لديك طلب قيد الانتطار، يمكنك المحاولة بعد إكتمال الرحلة أو إلغاء الطلب الحالي.\nللإلغاء أرسل الرقم 1"
		return NoOp
	}

	request, err := ctx.requestService.FindByChannelUser(ctx.requestMessage.ChannelType, ctx.requestMessage.Sender.ID)
	//request == nil this shouldn't happen
	if err != nil {
		logEntry.Errorln("error, cannot find request", err)
		ctx.responseMessage.Text = "حدث خطأ أثناء إجراء طلبك الرجاء المعاودة لاحقا أو الاتصال بخدمة العملاء"
		return Failed
	}
	request.Status = types.Cancelled
	_, err = ctx.requestService.Update(request)
	if err != nil {
		logEntry.Errorln("error, cannot find request", err)
		ctx.responseMessage.Text = "حدث خطأ أثناء إجراء طلبك الرجاء المعاودة لاحقا أو الاتصال بخدمة العملاء"
		return Failed
	}
	ctx.responseMessage.Text = "تم إلغاء الطلب"
	return NoOp

}

func NewRequestFSM() *StateMachine {
	return &StateMachine{
		States: States{
			types.Default: State{
				Events: Events{
					LocationReceived: types.OriginLocationUpdated,
					Start:            types.Started,
					TextReceived:     types.Started,
					ContactReceived:  types.Started,
					InvalidInput:     types.Started,
					Failed:           types.Started,
				},
			},
			types.Started: State{
				Action: &StartAction{},
				Events: Events{
					Start:            types.Started,
					LocationReceived: types.OriginLocationUpdated,
					TextReceived:     types.Started,
					ContactReceived:  types.Started,
					InvalidInput:     types.Started,
					Failed:           types.Started,
				},
			},
			types.OriginLocationUpdated: State{
				Action: &UpdateOrigLocationAction{},
				Events: Events{
					LocationReceived: types.DestLocationUpdated,
					TextReceived:     types.OriginLocationUpdated,
					ContactReceived:  types.OriginLocationUpdated,
					InvalidInput:     types.OriginLocationUpdated,
					Failed:           types.Started,
				},
			},
			types.DestLocationUpdated: State{
				Action: &UpdateDestLocationAction{},
				Events: Events{
					ContactReceived:  types.ContactUpdated,
					TextReceived:     types.DestLocationUpdated,
					LocationReceived: types.DestLocationUpdated,
					InvalidInput:     types.OriginLocationUpdated,
					Failed:           types.Started,
				},
			},
			types.ContactUpdated: State{
				Action: &UpdateContactAction{},
				Events: Events{
					TextReceived: types.Confirmed,
					//wrong events cause transition to the same state
					ContactReceived:  types.ContactUpdated,
					LocationReceived: types.ContactUpdated,
					InvalidInput:     types.DestLocationUpdated,
					Failed:           types.Started,
				},
			},
			types.Confirmed: State{
				Action: &ConfirmAction{},
				Events: Events{
					//wrong events cause transition to the same state
					TextReceived: types.Cancelled,
					Failed:       types.Default,
				},
			},
			types.Cancelled: State{
				Action: &CancellAction{},
			},
		},
	}
}
