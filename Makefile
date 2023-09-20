OBU_BINARY_NAME=obu
RECEIVER_BINARY_NAME=receiver
CALCULATOR_BINARY_NAME=calculator
AGGREAGTOR_BINARY_NAME=agg
GATE_BINARY_NAME=gate
obu:
	@go build -o bin/$(OBU_BINARY_NAME) obu/main.go
	@./bin/$(OBU_BINARY_NAME)
	
receiver:
	@go build -o bin/$(RECEIVER_BINARY_NAME) ./data_receiver
	@./bin/$(RECEIVER_BINARY_NAME)

calculator:
	@go build -o bin/$(CALCULATOR_BINARY_NAME) ./distance_calculator
	@./bin/$(CALCULATOR_BINARY_NAME)

agg:
	@go build -o bin/$(AGGREAGTOR_BINARY_NAME) ./aggregator
	@./bin/$(AGGREAGTOR_BINARY_NAME)

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	types/ptypes.proto

gate:
	@go build -o bin/$(GATE_BINARY_NAME) ./gateway
	@./bin/$(GATE_BINARY_NAME)

.PHONY: obu agg