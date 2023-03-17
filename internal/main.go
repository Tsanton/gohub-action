package main

import (
	"fmt"
	"os"
)

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func Greet(greeting string, greeted string) {
	fmt.Printf("%s, %s!", greeting, greeted)
}

func main() {
	greeting := getenv("GREETING", "Hello")
	greeted := getenv("GREETED", "World")
	Greet(greeting, greeted)
}
