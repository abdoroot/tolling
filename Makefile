obu:
	@go run OBU/main.go
receiver:
	@go run ./data_receiver
calc:
	@go run ./distance_calculator
agg:
	@go run ./aggregator
proto:
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative types/ptypes.proto
	
.PHONY:obu

