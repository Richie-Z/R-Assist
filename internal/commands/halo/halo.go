package halo

import (
	"fmt"

	"github.com/richie-z/whatsapp-bot/internal/service"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

func Run(client *whatsmeow.Client, evt *events.Message, args string) {
	message := fmt.Sprintf("Hi %s! this is Richie's bot answering", args)
	service.SendReply(client, evt.Info.Sender, message)

	err := service.LogMessage(evt.Info.Sender.String(), evt.Message.GetConversation(), message)
	if err != nil {
		fmt.Println("Failed to log message:", err)
	}
}
