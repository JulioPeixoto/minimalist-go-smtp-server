package main

import (
	"fmt"
	"log"
	"minimalist-go-smtp-server/service/internal/smtp"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	server := smtp.NewServer("localhost", 2525)

	go func() {
		fmt.Println("Server started on port 2525")
		if err := server.Start(); err != nil {
			log.Fatalf("Erro ao iniciar o servidor: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	server.Stop()
	fmt.Println("Server stopped")
}
