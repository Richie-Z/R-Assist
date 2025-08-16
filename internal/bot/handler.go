package bot

import (
	"fmt"
	"strings"

	"github.com/richie-z/whatsapp-bot/internal/commands/halo"
	"github.com/richie-z/whatsapp-bot/internal/commands/help"
	"github.com/richie-z/whatsapp-bot/internal/commands/ingatkan"
	"github.com/richie-z/whatsapp-bot/internal/service"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

type Handler struct {
	client *whatsmeow.Client
}

func NewHandler(client *whatsmeow.Client) *Handler {
	return &Handler{client: client}
}

func (h *Handler) Route(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		if v.Message.GetConversation() != "" {
			h.handleTextMessage(v)
		}
	}
}

func (h *Handler) handleTextMessage(evt *events.Message) {
	text := strings.TrimSpace(evt.Message.GetConversation())
	if !strings.HasPrefix(text, "!") {
		return
	}

	parts := strings.SplitN(text, " ", 2)
	cmd := strings.ToLower(parts[0])
	args := ""
	if len(parts) > 1 {
		args = parts[1]
	}

	fmt.Printf("%s command received: %s \n", cmd, args)

	switch cmd {
	case "!halo":
		halo.Run(h.client, evt, args)
	case "!ingatkan", "!ingatin":
		ingatkan.Run(h.client, evt, args)
	case "!help", "!tolong":
		help.Run(h.client, evt, args)
	default:
		service.SendReply(h.client, evt.Info.Sender, "Perintah tidak dikenal")
		help.Run(h.client, evt, args)
	}
}
