package halo

import (
	"fmt"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

func Run(client *whatsmeow.Client, evt *events.Message, args string) {
	fmt.Print("hello")
}
