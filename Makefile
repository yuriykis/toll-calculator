OBU_BINARY_NAME=obu
RECEIVER_BINARY_NAME=receiver
CALCULATOR_BINARY_NAME=calculator

obu:
	@go build -o bin/$(OBU_BINARY_NAME) obu/main.go
	@./bin/$(OBU_BINARY_NAME)
	
receiver:
	@go build -o bin/$(RECEIVER_BINARY_NAME) ./data_receiver
	@./bin/$(RECEIVER_BINARY_NAME)

calculator:
	@go build -o bin/$(CALCULATOR_BINARY_NAME) ./distance_calculator
	@./bin/$(CALCULATOR_BINARY_NAME)

.PHONY: obu