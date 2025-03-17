.PHONY: start-server start-agent metricstest

# Запускать из metrigo/

# make start-server a=<HOST:PORT>
start-server:
	@go build -o cmd/server/server cmd/server/*.go
	@./cmd/server/server $(if $(a),-a "$(a)",)

# make start-server a=<HOST:PORT>
start-agent:
	@go build -o cmd/agent/agent cmd/agent/*.go
	@./cmd/agent/agent $(if $(a),-a "$(a)",) $(if $(p),-p "$(p)",) $(if $(r),-r "$(r)",) $(if $(t),-t "$(t)",)

# make metricstest iter=<НОМЕР ИТЕРА>
metricstest:
	@go build -o cmd/server/server cmd/server/*.go
	@go build -o cmd/agent/agent cmd/agent/*.go
	@metricstest -test.v -test.run="^TestIteration$(iter)" -source-path=. -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -server-port=8080
