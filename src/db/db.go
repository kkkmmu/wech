package db

type Database interface {
	Save(key string, value []byte) error
	Get(key string) ([]byte, error)
	Destroy(string) error
}
