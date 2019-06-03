package sdk

import (
	"sync"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/status"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"

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
	MSPID    string
}

// NewSDK return new GoSDK
func NewSDK(configPath string, userCtx UserContext) *GoSDK {
	cfg := config.FromFile(configPath)
	// logprovider := &log.LogProvider{}
	// sdk, err := fabsdk.New(cfg, fabsdk.WithLoggerPkg(logprovider))
	sdk, err := fabsdk.New(cfg)
	if err != nil {
		return nil
	}
	service := &GoSDK{
		sdk:       sdk,
		chclients: make(map[string]*channel.Client),
		eventRegs: make(map[string]eventWrapper),
		userCtx:   userCtx,
	}
	return service
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
func (s *GoSDK) Query(channel, chaincode, fcn string, args [][]byte) RspMsg {
	return RspMsg{}
}

// RegisterEvent register chaincode event
func (s *GoSDK) RegisterEvent(channel, ccName, event string) <-chan *Event {
	return nil
}

func (s *GoSDK) getClient(channelID string) *channel.Client {
	s.chMutex.RLock()
	client, ok := s.chclients[channelID]
	s.chMutex.RUnlock()
	if ok {
		return client
	}
	context := s.sdk.ChannelContext(channelID,
		fabsdk.WithUser(s.userCtx.UserName), fabsdk.WithOrg(s.userCtx.MSPID))
	client, err := channel.New(context)
	if err != nil {
		panic(err)
	}
	s.chMutex.Lock()
	s.chclients[channelID] = client
	s.chMutex.Unlock()
	return client
}
