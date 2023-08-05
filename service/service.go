package service

import (
	"fmt"
	"net/http"

	"github.com/Gigfinder-io/util.gigfinder.io/db"
	"github.com/Gigfinder-io/util.gigfinder.io/db/models"
	"github.com/Gigfinder-io/util.gigfinder.io/jwt"
	"github.com/Gigfinder-io/util.gigfinder.io/log"
	"github.com/Gigfinder-io/util.gigfinder.io/stat"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// HandlerFunc is the function signature for the handlers this package can install.
type HandlerFunc func(http.ResponseWriter, *http.Request, primitive.ObjectID) (int, []byte, error)

// CreateHandler installs a handler function with a method route and permissions specification
func CreateHandler(path string, handler HandlerFunc, method string, requireLogin bool) {
	http.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		defer func() { // Panic Protection
			r := recover()
			if r != nil {
				log.Msgf(log.V, "recovered from panic: %v", r)
				stat.Atomic(stat.ServerErrors)
				stat.Error(fmt.Errorf("recovered from panic: %v", r), req.Method, path)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		recordMS := stat.StartTimer(path)

		if req.Method != method {
			log.Msgf(log.VVV, "received request with invalid method")
			stat.Atomic(stat.UserErrors)
			stat.Error(fmt.Errorf("received request with invalid method"), req.Method, path)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		log.Msgf(log.V, "received request [%v %v]", method, path)

		uid := primitive.NilObjectID
		if requireLogin {
			u, err := extractUserID(req)
			if err != nil {
				log.Msgf(0, "received request with invalid user, err: %v", err)
				w.WriteHeader(http.StatusUnauthorized)
				stat.Atomic(stat.UnauthorizedErrors)
				stat.Error(fmt.Errorf("received request with invalid user, err: %v", err), req.Method, path)
				w.Write([]byte(`{"success":false,"error":{"login":{"absent":true}}}`))
				return
			}
			uid = u
		}

		status, body, err := handler(w, req, uid)
		if err != nil {
			stat.Error(err, req.Method, path)
			level := 2
			if status == http.StatusInternalServerError {
				level = 0
				stat.Atomic(stat.ServerErrors)
			} else if status != http.StatusOK && status != http.StatusAccepted {
				stat.Atomic(stat.UserErrors)
			}
			log.Msgf(level, "encountered request error: [%v]", err)
		}

		w.WriteHeader(status)
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)

		recordMS()
	})
}

// CreateAdminHandler installs an admin handler function with a method route.
// Differs from CreateHandler by hardcoding admin requirements
func CreateAdminHandler(path string, handler HandlerFunc, method string) {
	http.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		defer func() { // Panic Protection
			r := recover()
			if r != nil {
				log.Msgf(log.V, "recovered from panic: %v", r)
				stat.Atomic(stat.ServerErrors)
				stat.Error(fmt.Errorf("recovered from panic: %v", r), req.Method, path)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		recordMS := stat.StartTimer(path)

		if req.Method != method {
			log.Msgf(log.VVV, "received request with invalid method")
			stat.Atomic(stat.UserErrors)
			w.WriteHeader(http.StatusMethodNotAllowed)
			stat.Error(fmt.Errorf("received request with invalid method"), req.Method, path)
			return
		}

		log.Msgf(log.V, "received request [%v %v]", method, path)

		// Extract user ID
		uid := primitive.NilObjectID
		u, err := extractUserID(req)
		if err != nil {
			log.Msgf(0, "received request with invalid user, err: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			stat.Atomic(stat.UnauthorizedErrors)
			stat.Error(fmt.Errorf("received request with invalid user, err: %v", err), req.Method, path)
			w.Write([]byte(`{"success":false,"error":{"login":{"absent":true}}}`))
			return
		}
		uid = u

		// Check that user is admin, required for every route on this service
		var user *models.User
		err = db.FindOne("users", db.Query{"_id": uid}, &user)
		if err != nil || user == nil {
			log.Msgf(log.V, "could not find user: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			stat.Atomic(stat.ServerErrors)
			w.Write([]byte(`{"success":false,"error":{"login":{"absent":true}}}`))
			return
		}
		if !user.Admin {
			log.Msgf(log.VVV, "received request from non admin user")
			w.WriteHeader(http.StatusUnauthorized)
			stat.Atomic(stat.UnauthorizedErrors)
			w.Write([]byte(`{"success":false,"error":{"login":{"unauthorized":true}}}`))
			return
		}

		status, body, err := handler(w, req, primitive.ObjectID{})
		if err != nil {
			stat.Error(err, req.Method, path)
			level := 2
			if status == http.StatusInternalServerError {
				level = 0
				stat.Atomic(stat.ServerErrors)
			} else if status != http.StatusOK && status != http.StatusAccepted {
				stat.Atomic(stat.UserErrors)
			}
			log.Msgf(level, "encountered request error: [%v]", err)
		}

		w.WriteHeader(status)
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)

		recordMS()
	})
}

func extractUserID(req *http.Request) (primitive.ObjectID, error) {
	header := req.Header.Get("Authorization")
	if len(header) < 8 {
		return primitive.NilObjectID, fmt.Errorf("no, or invalid bearer token provided")
	}

	id, err := jwt.Parse(header[7:])
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("could not parse jwt: %v", err)
	}

	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("could not cast to objectID: %v", err)
	}

	return oID, nil
}
