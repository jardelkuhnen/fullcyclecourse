

First of all, you should have docker running locally. 

1 - run `docker-compose up` to start the containers of rabbitmq
2 - run curl the project `go run cmd/producer/main.go` to start the producer of messages
3 - run curl the project `go run cmd/consumer/main.go` to start the consumer of messages
