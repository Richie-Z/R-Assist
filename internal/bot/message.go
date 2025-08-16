package bot

import (
	"context"
	"fmt"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

func SendReply(client *whatsmeow.Client, to types.JID, text string) {
	msg := &waE2E.Message{
		Conversation: proto.String(text),
	}
	_, err := client.SendMessage(context.Background(), to, msg)
	if err != nil {
		fmt.Println("Error sending message:", err)
	}
}
