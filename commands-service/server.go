package main

import (
	"github.com/serkanerip/commands-service/command"
	"github.com/serkanerip/commands-service/internal/infra"
	"github.com/serkanerip/commands-service/internal/mongodb"
)

func main() {
	client := mongodb.CreateClient()
	commandRepo := command.NewCommandMongoRepository(client)
	commandService := command.NewCommandService(commandRepo)
	platformRouter := command.NewPlatformRouter(*commandService)

	forever := make(chan bool)
	// rabbitmq
	rabbitMQClient := infra.NewRabbitMQClient(commandRepo)
	go rabbitMQClient.Run()

	httpHandler := infra.NewHttpHandler()
	httpHandler.RegisterRoutes(platformRouter)
	go httpHandler.Run()
	//end rabbitmq

	<-forever
}
