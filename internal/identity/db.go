package identity

import (
	"log"
	"os"
)

type DB struct {
	file *os.File
	primary
}

func NewDB(path string) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	return DB{
		file,
	}
}

func Get(path) * {

}
