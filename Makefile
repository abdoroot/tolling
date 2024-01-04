obu:
	go run OBU/main.go
receiver:
	go build -o data_receiver/bin/dreceiver ./data_receiver && ./data_receiver
.PHONY:obu

