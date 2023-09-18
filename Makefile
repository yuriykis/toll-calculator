OBU_BINARY_NAME=obu
RECEIVER_BINARY_NAME=receiver
CALCULATOR_BINARY_NAME=calculator
AGGREAGTOR_BINARY_NAME=aggregator

obu:
	@go build -o bin/$(OBU_BINARY_NAME) obu/main.go
	@./bin/$(OBU_BINARY_NAME)
	
receiver:
	@go build -o bin/$(RECEIVER_BINARY_NAME) ./data_receiver
	@./bin/$(RECEIVER_BINARY_NAME)

calculator:
	@go build -o bin/$(CALCULATOR_BINARY_NAME) ./distance_calculator
	@./bin/$(CALCULATOR_BINARY_NAME)

aggregator:
	@go build -o bin/$(AGGREAGTOR_BINARY_NAME) ./aggregator
	@./bin/$(AGGREAGTOR_BINARY_NAME)

.PHONY: obu, aggregator