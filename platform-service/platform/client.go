package platform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/serkanerip/platform-service/config"
)

type CommandDataClient interface {
	SendPlatformToCommand(PlatformDTO)
}

type HttpCommandDataClient struct {
	c http.Client
}

func NewHttpCommandDataClient() *HttpCommandDataClient {
	return &HttpCommandDataClient{
		c: http.Client{Timeout: time.Duration(1) * time.Second},
	}
}

func (h *HttpCommandDataClient) SendPlatformToCommand(dto PlatformDTO) {
	b, err := json.Marshal(dto)
	if err != nil {
		log.Printf("cannot marshal platform dto err is: ", err)
		return
	}

	uri := fmt.Sprintf("%s/api/c/platforms", config.ENV.CommandService)
	resp, err := h.c.Post(uri, "application/json", bytes.NewReader(b))

	if err != nil || resp.StatusCode != 201 {
		log.Printf("--> Sync POST to CommandService was Not OK!  err is: %v", err)
		return
	}

	log.Printf("--> Sync POST to CommandService was OK! ")
}
