package startup

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Gigfinder-io/util.gigfinder.io/db"
	"github.com/Gigfinder-io/util.gigfinder.io/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	DB struct {
		Address  string `yaml:"address"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"db"`
}

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
	fName := "./resource/application.yml"
	yamlFile, err := ioutil.ReadFile(fName)
	if err != nil {
		return fmt.Errorf("failed to read file [\"%v\"]: %v", fName, err)
	}
	config := &AppConfig{}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		return fmt.Errorf("error unmarshalling: %v", err)
	}
	db.Address = config.DB.Address
	db.User = config.DB.User
	db.Pass = config.DB.Password

	return nil
}
