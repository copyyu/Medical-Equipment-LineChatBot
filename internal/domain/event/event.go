package event

import "time"

// Event types for the medical equipment system
const (
	EquipmentCreated = "equipment.created"
	EquipmentUpdated = "equipment.updated"
	EquipmentDeleted = "equipment.deleted"
	TicketCreated    = "ticket.created"
	TicketUpdated    = "ticket.updated"
)

// Event represents a domain event that is published through the event bus
type Event struct {
	Type      string      `json:"type"`
	Payload   interface{} `json:"payload"`
	Timestamp time.Time   `json:"timestamp"`
}

// NewEvent creates a new event with the current timestamp
func NewEvent(eventType string, payload interface{}) Event {
	return Event{
		Type:      eventType,
		Payload:   payload,
		Timestamp: time.Now(),
	}
}
