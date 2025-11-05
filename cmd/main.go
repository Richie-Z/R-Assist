package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/richie-z/R-Assist/internal/bot"
)

func main() {
	client := bot.NewConnection()

	handler := bot.NewHandler(client)
	client.AddEventHandler(handler.Route)

	fmt.Println("WhatsApp bot is running...")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	client.Disconnect()
}
