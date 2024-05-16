package simulator

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"sync"
)

func generateToken() (string, error) {
	var s string
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Println("generateToken err: ", err)
		return s, err
	}
	s = hex.EncodeToString(b)
	return s, err
}

func writeMap(auth []AuthRequest) (sync.Map, error) {
	var m sync.Map
	var err error
	for i := 0; i < len(auth); i++ {
		//Store key-valeu
		m.Store(auth[i].Auth.Username, auth[i].Auth.Password)
		//Load and print the value
		if value, ok := m.Load(auth[i].Auth.Username); ok {
			log.Printf("Goroutine %d: Key %s - Value %d\n", i, auth[i].Auth.Username, value)
		} else {
			log.Println("Error load from map", err)
			return m, err
		}
	}

	return m, err
}
