OBU_BINARY_NAME=obu
RECEIVER_BINARY_NAME=receiver

obu:
	@go build -o bin/$(OBU_BINARY_NAME) obu/main.go
	@./bin/$(OBU_BINARY_NAME)
	
receiver:
	@go build -o bin/$(RECEIVER_BINARY_NAME) data_receiver/main.go
	@./bin/$(RECEIVER_BINARY_NAME)

.PHONY: obu 