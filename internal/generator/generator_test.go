package generator

import (
	"log"
	"math/rand"
	"testing"
	"time"

	"practice/storage"
)

func TestGenerator(t *testing.T) {
	telematicCh := make(chan *storage.Car)
	log.Println("Starting generator...")
	go GenerateTelematic(rand.Intn(3), telematicCh)

	emptyStruct := &storage.Car{}

	go func() {
		for {
			val, ok := <-telematicCh
			if ok == false {
				break // exit break loop
			} else {
				if val == emptyStruct {
					log.Fatal("Generator error: returned empty struct")
				}
			}
		}
	}()

	time.Sleep(5 * time.Second)
	close(telematicCh)
}
