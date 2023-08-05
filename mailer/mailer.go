package mailer

import (
	"fmt"
	"net/http"

	"github.com/Gigfinder-io/util.gigfinder.io/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	Host = "localhost"
	Port = 8010
)

const (
	NewMessageID      = 0
	RequestPasswordID = 1
	NewUserID         = 2
	VerifyUserID      = 3
	NearbyUserID      = 4
)

type mail struct {
	messageType int
	recipient   string
	payload     string
}

var (
	mailQueue chan mail
	done      chan struct{}
)

func Start() {
	mailQueue = make(chan mail, 5)
	done = make(chan struct{})
	log.Msgf(log.V, "starting mail service")

	go func() {
		for m := range mailQueue {
			log.Msg(log.VV, "sending mail to server")
			err := makeRequest(m)
			if err != nil {
				log.Msgf(log.V, "could not send mail to server: %v", err)
			}
		}
		close(done)
	}()
}

func Close() {
	log.Msgf(log.V, "closing mail service")
	close(mailQueue)
	<-done
}

func makeRequest(m mail) error {
	url := fmt.Sprintf("http://%v:%v/email/send?pl=%v&rp=%v&mt=%v", Host, Port, m.payload, m.recipient, m.messageType)
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return fmt.Errorf("could not make post request: %v", err)
	}

	if resp.StatusCode != 202 {
		return fmt.Errorf("did not receive 202 response, received [%v] instead", resp.StatusCode)
	}

	return nil
}

func Queue(emailType int, recipient string, payload string) {
	mailQueue <- mail{recipient: recipient, payload: payload, messageType: emailType}
}

func MakeNewsletterRequest(id primitive.ObjectID) error {
	url := fmt.Sprintf("http://%v:%v/email/newsletter?id=%v", Host, Port, id.Hex())
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return fmt.Errorf("could not make post request: %v", err)
	}

	if resp.StatusCode != 202 {
		return fmt.Errorf("did not receive 202 response for [POST %v], received [%v] instead", url, resp.StatusCode)
	}

	return nil
}
