package main

import (
	"log"

	"github.com/ByteNode1/GolangCalculator/internal/server"
)

func main() {
	if err := server.New().Start(); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
