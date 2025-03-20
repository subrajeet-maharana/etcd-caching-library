package cache

// Cache defines the interface for different cache implementations
type Cache interface {
	Get(key string) (string, bool)
	Set(key, value string)
	Delete(key string)
}
