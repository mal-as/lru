package lru

// Cache - интерфейс для получения и установки значений в кэше
type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
}
