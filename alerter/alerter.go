package alerter

import (
	"fmt"
	"net/http"

	"github.com/Gigfinder-io/util.gigfinder.io/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	Address = "localhost:9080"
)

type alert struct {
	userID string
}

var (
	alertQueue chan alert
	done       chan struct{}
)

func Start() {
	alertQueue = make(chan alert, 5)
	done = make(chan struct{})
	log.Msgf(log.V, "starting alert service")

	go func() {
		for m := range alertQueue {
			log.Msg(log.VV, "sending alert to server")
			err := makeRequest(m)
			if err != nil {
				log.Msgf(log.V, "could not send alert to server: %v", err)
			}
		}
		close(done)
	}()
}

func Close() {
	log.Msgf(log.V, "closing alert service")
	close(alertQueue)
	<-done
}

func makeRequest(m alert) error {
	url := fmt.Sprintf("http://%v/alerts/hooks/usersignup?uid=%v", Address, m.userID)
	log.Msgf(0, "making request [%v] to alert.gigfinder.io", url)
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return fmt.Errorf("could not make post request")
	}

	if resp.StatusCode != 202 {
		return fmt.Errorf("did not receive 202 response, instead received [%v]", resp.StatusCode)
	}

	return nil
}

func Queue(userID primitive.ObjectID) {
	alertQueue <- alert{userID.Hex()}
}
