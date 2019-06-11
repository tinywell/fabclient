package sdk

import (
	"fmt"
	"sync"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/status"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/events/deliverclient/seek"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/util/pathvar"

	"github.com/tinywell/fabclient/pkg/common"
)

// GoSDK implement for sdk base on fabric-go-sdk
type GoSDK struct {
	sdk       *fabsdk.FabricSDK
	chclients map[string]*channel.Client
	chMutex   sync.RWMutex
	eventRegs map[string]eventWrapper
	userCtx   UserContext
}

type eventWrapper struct {
	StatusChan chan bool
	EventReg   fab.Registration
}

// UserContext user context for the sdk
type UserContext struct {
	UserName string
	Org      string
}

// NewSDK return new GoSDK
func NewSDK(configPath string, userCtx UserContext) (*GoSDK, error) {
	cfg := config.FromFile(pathvar.Subst(configPath))
	// logprovider := &log.LogProvider{}
	// sdk, err := fabsdk.New(cfg, fabsdk.WithLoggerPkg(logprovider))
	sdk, err := fabsdk.New(cfg)
	if err != nil {
		return nil, err
	}
	service := &GoSDK{
		sdk:       sdk,
		chclients: make(map[string]*channel.Client),
		eventRegs: make(map[string]eventWrapper),
		userCtx:   userCtx,
	}
	return service, nil
}

// Excute chaincode execution
func (s *GoSDK) Excute(channelID, chaincode, fcn string, args [][]byte) RspMsg {
	client := s.getClient(channelID)
	request := channel.Request{
		ChaincodeID:  chaincode,
		Fcn:          fcn,
		Args:         args,
		TransientMap: nil,
	}
	response, err := client.Execute(request)
	rspmsg := RspMsg{
		TxID: string(response.Proposal.TxnID),
	}
	if err != nil {
		rspmsg.Code = common.RspServerError
		rspmsg.Data = []byte(err.Error())
		if sta, ok := err.(*status.Status); ok {
			if sta.Code == status.Timeout.ToInt32() {
				rspmsg.Code = common.RspTimeout
			}
		}
		return rspmsg
	}
	rspmsg.Code = common.RspSuccess
	rspmsg.Data = response.Payload
	return rspmsg
}

// Query chaincode query
func (s *GoSDK) Query(channelID, chaincode, fcn string, args [][]byte) RspMsg {
	client := s.getClient(channelID)
	request := channel.Request{
		ChaincodeID: chaincode,
		Fcn:         fcn,
		Args:        args,
	}
	response, err := client.Query(request)
	rspmsg := RspMsg{
		TxID: string(response.Proposal.TxnID),
	}
	if err != nil {
		rspmsg.Code = common.RspServerError
		rspmsg.Data = []byte(err.Error())
		return rspmsg
	}
	rspmsg.Code = common.RspSuccess
	rspmsg.Data = response.Payload
	return rspmsg
}

// RegisterEvent register chaincode event
func (s *GoSDK) RegisterEvent(channel, ccName, eventName string) (<-chan *Event, error) {
	eventKey := fmt.Sprintf("%s-%s-%s", channel, ccName, eventName)
	_, exist := s.eventRegs[eventKey]
	if exist {
		return nil, fmt.Errorf("event %s already registered", eventName)
	}
	context := s.sdk.ChannelContext(channel,
		fabsdk.WithUser(s.userCtx.UserName), fabsdk.WithOrg(s.userCtx.Org))
	client, err := event.New(context, event.WithBlockEvents(), event.WithSeekType(seek.Newest))
	reg, notifier, err := client.RegisterChaincodeEvent(ccName, eventName)
	if err != nil {
		return nil, err
	}
	statusChan := make(chan bool)
	s.eventRegs[eventKey] = eventWrapper{StatusChan: statusChan, EventReg: reg}
	return convertEvent(notifier, statusChan), nil
}

// UnRegisterEvent chaincode event unresgitered
func (s *GoSDK) UnRegisterEvent(channel, ccName, eventName string) error {
	eventKey := fmt.Sprintf("%s-%s-%s", channel, ccName, eventName)
	e, exist := s.eventRegs[eventKey]
	if !exist {
		return fmt.Errorf("This event had not registed")
	}
	context := s.sdk.ChannelContext(channel,
		fabsdk.WithUser(s.userCtx.UserName), fabsdk.WithOrg(s.userCtx.Org))
	client, err := event.New(context, event.WithBlockEvents())
	if err != nil {
		return err
	}

	client.Unregister(e.EventReg)
	e.StatusChan <- true
	delete(s.eventRegs, eventKey)
	return nil
}

func (s *GoSDK) getClient(channelID string) *channel.Client {
	s.chMutex.RLock()
	client, ok := s.chclients[channelID]
	s.chMutex.RUnlock()
	if ok {
		return client
	}
	s.chMutex.Lock()
	defer s.chMutex.Unlock()
	context := s.sdk.ChannelContext(channelID,
		fabsdk.WithUser(s.userCtx.UserName),
		fabsdk.WithOrg(s.userCtx.Org))
	client, err := channel.New(context)
	if err != nil {
		panic(err)
	}
	s.chclients[channelID] = client
	return client
}

func convertEvent(notifier <-chan *fab.CCEvent, quit chan bool) <-chan *Event {
	ccEvent := make(chan *Event)
	go func() {
		for {
			select {
			case event := <-notifier:
				ccEvent <- &Event{
					TxID:        event.TxID,
					CCName:      event.ChaincodeID,
					EventName:   event.EventName,
					Payload:     event.Payload,
					BlockNumber: event.BlockNumber,
				}
			case <-quit:
				return
			}
		}
	}()
	return ccEvent
}
