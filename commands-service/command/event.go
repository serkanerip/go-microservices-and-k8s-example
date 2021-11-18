package command

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/serkanerip/commands-service/events"
)

type CommandEventProcessor struct {
	commandRepository CommandRepo
}

const (
	PLATFORM_PUBLISHED events.EventType = iota
	UNDETERMINED
)

func NewCommandEventProcessor(cr CommandRepo) *CommandEventProcessor {
	return &CommandEventProcessor{
		commandRepository: cr,
	}
}

func (c *CommandEventProcessor) ProcessEvent(message string) {
	eventType, err := c.determineEvent(message)
	if err != nil {
		log.Printf("cannot determine event type err is: %v", err)
	}
	switch eventType {
	case PLATFORM_PUBLISHED:
		c.addPlatform(message)
	}
}

func (c *CommandEventProcessor) determineEvent(message string) (events.EventType, error) {
	var dto events.GenericEventDTO
	if err := json.Unmarshal([]byte(message), &dto); err != nil {
		return UNDETERMINED, fmt.Errorf("cannot unmarshal message err is: %v", err)
	}
	switch dto.Event {
	case "Platform_Published":
		log.Println("Platform published event detected")
		return PLATFORM_PUBLISHED, nil
	}
	return UNDETERMINED, fmt.Errorf("cannot determine event type")
}

func (c *CommandEventProcessor) addPlatform(platformPublishedMessage string) error {
	var dto PlatformPublishedDTO
	if err := json.Unmarshal([]byte(platformPublishedMessage), &dto); err != nil {
		return fmt.Errorf("cannot unmarshal message err is: %v", err)
	}

	if exists, err := c.commandRepository.ExternalPlatformExists(dto.Id); err != nil || exists {
		return nil
	}

	c.commandRepository.CreatePlatform(Platform{
		ExternalId: dto.Id,
		Name:       dto.Name,
	})

	return nil
}
