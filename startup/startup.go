package startup

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Gigfinder-io/util.gigfinder.io/db"
	"github.com/Gigfinder-io/util.gigfinder.io/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
	ID          primitive.ObjectID `bson:"_id"`
	ServiceName string             `bson:"serviceName"`
	Port        int                `bson:"port"`
	Identifier  string             `bson:"identifier"`
	Version     int                `bson:"version"`
	RunnerArgs  []string           `bson:"runArgs"`
	BuildArgs   []string           `bson:"buildArgs"`
	ReleasePort int                `bson:"releasePort"`
}

// Go updates os.Args according to arguments fetched from a database. run this pre-flag.Parse() (or run it again after)
func Go(serviceName string) error {
	log.Msgf(log.V, "fetching startup information...")
	err := findDBAccess()
	if err != nil {
		return fmt.Errorf("startup failed: %v", err)
	}

	// Query DB
	err = db.Connect()
	if err != nil {
		return fmt.Errorf("failed to start db connection: %v", err)
	}
	defer db.Disconnect()

	svc := service{}
	err = db.FindOne("services", db.Query{"serviceName": serviceName}, &svc)
	if err != nil {
		return fmt.Errorf("Could not find service with serviceName[%v]: %v", serviceName, err)
	}
	log.Msgf(log.V, "found startup information for [%v-%v:v%v]", svc.ID.Hex(), svc.ServiceName, svc.Version)
	log.Msgf(2, "Startup information: [%v]", svc)

	newArgs := []string{os.Args[0]}
	newArgs = append(newArgs, svc.RunnerArgs...)
	newArgs = append(newArgs, fmt.Sprintf("-port=%v", svc.Port))
	os.Args = newArgs

	return nil
}

func findDBAccess() error {
	fName := "../args/db.txt"
	file, err := os.Open(fName)
	if err != nil {
		return fmt.Errorf("could not open file: [%v] error: %v", fName, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		txt := scanner.Text()
		if val := strings.TrimPrefix(txt, "addr="); val != txt {
			db.Address = val
			log.Msgf(3, "Found db address: [%v]", val)
		}
		if val := strings.TrimPrefix(txt, "usr="); val != txt {
			db.User = val
			log.Msgf(3, "Found db user: [%v]", val)
		}
		if val := strings.TrimPrefix(txt, "pass="); val != txt {
			db.Pass = val
			log.Msgf(3, "Found db password: [%v]", val)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("could not read db descriptor file, error: %v", err)
	}
	return nil
}
