package config

type DBType int

const (
	REDIS DBType = iota
	FLAT_FILE
)

type Config struct {
	identity_db      string
	identity_address string
	socket_address   string
}

func ParseFromPath(path string) {

}

func Parse(path string) *Config {

}
