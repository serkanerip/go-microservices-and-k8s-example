package command

import "github.com/serkanerip/commands-service/events"

type Platform struct {
	Id         string
	ExternalId string `json:"external_id"`
	Name       string `json:"name"`
	Commands   []Command
}

type PlatformPublishedDTO struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	events.GenericEventDTO
}
