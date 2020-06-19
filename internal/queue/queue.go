package queue

type QueueKey []byte
type QueueValue []byte

type IQueue interface {
	// This is blocking.
	Poll() QueueValue
	// This is non-blocking.
	Get(key QueueKey) QueueValue
	Push(key QueueKey, value QueueValue)

	// Attempts to reconnect returns true if succssful
	Reconnect() bool
}
