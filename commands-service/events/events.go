package events

type EventType int64

type GenericEventDTO struct {
	Event string `json:"event"`
}

type EventProcessor interface {
	ProcessEvent(message string)
}
