package cache

type CacheKey []byte
type CacheValue []byte

type ICache interface {
	// This is blocking.
	Poll() CacheValue
	// This is non-blocking.
	Keys() []CacheKey
	Get(key CacheKey) CacheValue
	Push(key CacheKey, value CacheValue)
	Update(key CacheKey, value CacheValue)

	// Attempts to reconnect returns true if succssful
	Reconnect() bool
}
