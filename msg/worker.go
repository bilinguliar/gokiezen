package msg

import (
	"context"
	"log"
)

type Request struct {
	Sender    string
	Recepient string
	Text      string
}

// Messenger is used to send text messages.
type Messenger interface {
	SendText(sender, msisdn, text string) error
}

// StartSendingMessages starts background worker that sends short messages.
func StartSendingMessages(ctx context.Context, mc chan Request, m Messenger) {
	var req Request

	for {
		select {
		case <-ctx.Done():
			return
		case req = <-mc:
			log.Printf("SMS requested, details: %#v\n", req)
			err := m.SendText(req.Sender, req.Recepient, req.Text)
			if err != nil {
				log.Println("Failed to send SMS, error:", err)
			}
		}
	}
}
