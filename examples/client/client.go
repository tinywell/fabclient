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
	params := handler.Params{}
	handler := handler.NewHandler(core, params)
	expHandler := exphandler.CreateHandler("mychannel", "public", "User1")
	err = handler.RegisterTxHandleFunc("EXP100", expHandler.SaveData)
	if err != nil {
		panic("register handler error")
	}
	err = handler.RegisterTxHandleFunc("EXP200", expHandler.ReadData)
	if err != nil {
		panic("register handler error")
	}
	server := exprest.NewServer(":8000")
	msgs := server.ReceiveMessage()
	ctx := context.Background()
	for msg := range msgs {
		handler.HandleMessage(ctx, msg)
	}
}
