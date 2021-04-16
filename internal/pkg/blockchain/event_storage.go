package blockchain

type EventStorage interface {
	Exists(eventId EventId) bool
}
