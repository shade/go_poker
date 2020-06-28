package cache

type CacheKey []byte
type CacheValue []byte

type ICache interface {
	// This is blocking.
	Poll(key CacheKey, values chan string)
	// This is non-blocking.
	Keys() []CacheKey
	Get(key CacheKey) CacheValue
	Push(key CacheKey, value CacheValue)
}
