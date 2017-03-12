package msg

import (
	"log"

	mb "github.com/messagebird/go-rest-api"
)

// Birdman or simply aviculturist. Knows how to deal with MessageBird.com API.
// Wraps original MessageBird client in order to avoid coupling with vendor specific structs in packages that will consume this functionality.
type Birdman struct {
	mbClient *mb.Client
	msgChan  chan Request
}

// NewMsgBirdClient creates instance of Birdman
func NewMsgBirdClient(token string, mc chan Request) *Birdman {
	c := &Birdman{
		mbClient: mb.New(token),
		msgChan:  mc,
	}

	balance, err := c.mbClient.Balance()
	if err != nil {
		log.Fatal("Unable to get balance")
	}

	if balance.Amount < 1 {
		log.Fatal("Balance is too low for proper voting campaign")
	}

	log.Printf("Current MessageBird.com balance type: %q, amount: %f\n", balance.Type, balance.Amount)

	return c
}

// SendTest sends SMS from sender to a recepient with provided text.
func (c *Birdman) SendText(sender, recepient, text string) error {
	m, err := c.mbClient.NewMessage(sender, []string{recepient}, text, &mb.MessageParams{})
	if err != nil {
		if err == mb.ErrResponse {
			for _, mbError := range m.Errors {
				log.Printf("Error: %#v\n", mbError)
			}
		}
		return err
	}

	return nil
}

// Lookup is used to get detailes about MSISDN. We need only country code.
func (c *Birdman) Lookup(msisdn string) (string, error) {
	lr, err := c.mbClient.Lookup(msisdn, &mb.LookupParams{})
	if err != nil {
		return "", err
	}

	return lr.CountryCode, nil
}

// RequestSMS adds SMS request to the channel, it will be send sometime in the future.
func (c *Birdman) RequestSMS(sender, recepient, text string) {
	c.msgChan <- Request{Sender: sender, Recepient: recepient, Text: text}
}
