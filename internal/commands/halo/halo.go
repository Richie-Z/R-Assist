package halo

import (
	"fmt"

	"github.com/richie-z/whatsapp-bot/internal/service"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

func Run(client *whatsmeow.Client, evt *events.Message, args string) {
	fmt.Println("Halo command received:", args)
	message := fmt.Sprintf("Hi %s,this is Richie's bot answering", args)
	service.SendReply(client, evt.Info.Sender, message)
}
