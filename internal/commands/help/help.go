package help

import (
	"fmt"

	"github.com/richie-z/whatsapp-bot/internal/service"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

func Run(client *whatsmeow.Client, evt *events.Message, args string) {
	message := "!halo\n!halo Chassie\n!help\n!tolong\n!ingatkan makan 20 menit lagi\n!ingatkan makan 5 detik lagi\n!ingatkan makan jam 10:00\n!ingatkan makan besok\n!ingatkan makan besok jam 18:00"
	service.SendReply(client, evt.Info.Sender, message)

	err := service.LogMessage(evt.Info.Sender.String(), evt.Message.GetConversation(), message)
	if err != nil {
		fmt.Println("Failed to log message:", err)
	}
}
