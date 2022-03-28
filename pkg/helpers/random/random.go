package random

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Animal returns a random animal name.
func Animal() string {
	return animals[rand.Intn(len(animals))]
}

// Adjective returns a random adjective.
func Adjective() string {
	return adjectives[rand.Intn(len(adjectives))]
}
