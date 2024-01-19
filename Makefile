obu:
	@go run OBU/main.go
receiver:
	@go run ./data_receiver
calc:
	@go run ./distance_calculator
agg:
	@go run ./aggregator	
.PHONY:obu

