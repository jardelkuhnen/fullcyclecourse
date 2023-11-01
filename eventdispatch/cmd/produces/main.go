package main

import "github.com/jardelkuhnen/eventdispatch/pkg/rabbitmq"

func main() {
	ch, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	err = rabbitmq.Publish(ch, "Hello Jardel!")
	if err != nil {
		panic(err)
	}
}
