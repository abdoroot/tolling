obu:
	@go run OBU/main.go
receiver:
	@go run ./data_receiver
calc:
	@go run ./distance_calculator
.PHONY:obu

