package utils

import (
	"fmt"
	"log"
)

// Returns the bytes read from the given name.
func MustReadMap(name string) []byte {
	bytes, er := Assets.Maps.ReadFile(fmt.Sprintf("maps/%s.json", name))
	if er != nil {
		log.Fatal(er)
	}

	return bytes
}
