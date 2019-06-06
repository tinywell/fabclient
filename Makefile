
mock: 
	rm -rf test/mocks

	mkdir -p test/mocks/sdk
	mockgen -package=msdk -source=pkg/sdk/sdk.go > test/mocks/sdk/sdk.go

	mkdir -p test/mocks/server
	mockgen -package=mserver -source=pkg/server/server.go > test/mocks/server/server.go

	mkdir -p test/mocks/handler
	mockgen -package=mhandler -source=pkg/handler/handler.go > test/mocks/handler/handler.go


.PHONY: mock