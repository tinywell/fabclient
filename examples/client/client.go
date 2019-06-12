package client

import (
	"context"

	exphandler "github.com/tinywell/fabclient/examples/handler"
	exprest "github.com/tinywell/fabclient/examples/server/restful"
	"github.com/tinywell/fabclient/pkg/common/config"
	"github.com/tinywell/fabclient/pkg/handler"
	"github.com/tinywell/fabclient/pkg/sdk"
)

// Cfg client config
type Cfg struct {
	ConfigFile string
}

// Run run client
func Run(cfg Cfg) {
	config.InitConfig(cfg.ConfigFile)

	usrCtx := sdk.UserContext{UserName: "User1", Org: config.GetOrgName()}
	gosdk, err := sdk.NewSDK(config.GetSDKConfig(), usrCtx)
	if err != nil {
		panic("create sdk error:" + err.Error())
	}
	core := handler.Core{
		TxHandler:  gosdk,
		SrcManager: nil,
	}
	h := handler.NewHandler(core, handler.Params{})
	expBox := exphandler.CreateHandler("mychannel", "public", "User1")
	err = h.FillHandlerFunc(expBox)
	if err != nil {
		panic("fill handler function error" + err.Error())
	}

	server := exprest.NewServer(":8000")
	msgs := server.ReceiveMessage()
	ctx := context.Background()
	for msg := range msgs {
		h.HandleMessage(ctx, msg)
	}
}
