package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/serkanerip/platform-service/config"
	"github.com/serkanerip/platform-service/internal/messagebroker"
	"github.com/serkanerip/platform-service/internal/mongodb"
	"github.com/serkanerip/platform-service/platform"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		message := map[string]string{
			"msg": "Pong!",
		}

		b, err := json.Marshal(&message)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(b)
	})

	//platformRepo := platform.PlatformInMemoryRepository{}
	conn, err := messagebroker.GetConn()
	if err != nil {
		log.Printf("cannot connect rabbit mq client err is: %v", err)
	}
	defer conn.Close()

	messageBusClient := platform.NewRabbitMQClient(conn)

	platformRepo := platform.NewMongoDBRepository(mongodb.CreateClient())
	platformRepo.Seed()
	commandDataClient := platform.NewHttpCommandDataClient()
	platformService := platform.NewPlatformService(platformRepo, commandDataClient, messageBusClient)
	platformRouter := platform.NewPlatformRouter(*platformService)

	r.Route("/api", func(r chi.Router) {
		platformRouter.SetupRoutes(r)
	})

	log.Printf("--> Command Service Endpoint: %s", config.ENV.CommandService)
	uri := fmt.Sprintf("0.0.0.0:%s", config.ENV.Port)
	if err := http.ListenAndServe(uri, r); err != nil {
		log.Fatalf("cannot start server err is :%v", err)
	}
}
