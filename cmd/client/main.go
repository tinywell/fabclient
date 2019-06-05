package main

import (
	"context"
	"flag"

	exphandler "github.com/tinywell/fabclient/examples/handler"
	exprest "github.com/tinywell/fabclient/examples/server/restful"
	"github.com/tinywell/fabclient/pkg/common/config"
	"github.com/tinywell/fabclient/pkg/handler"
	"github.com/tinywell/fabclient/pkg/sdk"
)

var (
	configFile string
)

func main() {
	flag.StringVar(&configFile, "config", "./config,yaml", "config path")
	flag.Parse()

	config.InitConfig(configFile)

	// sdkConfig := config.GetSDKConfig()
	// userCtx := sdk.UserContext{
	// 	UserName: "User1",
	// 	MSPID:    "msporg1",
	// }
	// sdkcore := sdk.NewSDK(sdkConfig, userCtx)

	core := handler.Core{
		TxHandler:  &sdk.MockTxHandler{},
		SrcManager: nil,
	}
	params := handler.Params{}
	handler := handler.NewHandler(core, params)
	expHandler := exphandler.CreateHandler("mychannel", "expcc", "User1")
	err := handler.RegisterTxHandleFunc("EXP100", expHandler.SaveData)
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
