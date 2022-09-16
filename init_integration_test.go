package libwekan

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbClient *mongo.Client

func TestMain(m *testing.M) {

	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		panic("ne peut démarrer mongodb")
	}
	// pulls an image, creates a container based on it and runs it
	mongodbContainerName := "keycloakUpdater-ti-" + strconv.Itoa(time.Now().Nanosecond())

	mongodb, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Name:       mongodbContainerName,
			Repository: "mongo",
			Tag:        "5.0",
			Env: []string{
				// username and password for mongodb superuser
				"MONGO_INITDB_ROOT_USERNAME=root",
				"MONGO_INITDB_ROOT_PASSWORD=password",
			},
			// Cmd: []string{"mongod", "--logpath /dev/null", "--oplogSize 128", "--quiet", "--noauth"},
		},
		func(config *docker.HostConfig) {
			// set AutoRemove to true so that stopped container goes away by itself
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
		},
	)
	if err != nil {
		fmt.Println(err.Error())
		kill(mongodb)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	_ = pool.Retry(func() error {
		var err error
		dbClient, err = mongo.Connect(
			context.TODO(),
			options.Client().ApplyURI(
				fmt.Sprintf("mongodb://root:password@localhost:%s", mongodb.GetPort("27017/tcp")),
			),
		)
		if err != nil {
			return err
		}
		if err := dbClient.Ping(context.TODO(), nil); err != nil {
			return err
		}

		err = restoreDump(mongodb)
		if err != nil {
			return err
		}

		nbUsers, err := dbClient.Database("wekan").Collection("users").CountDocuments(context.Background(), bson.M{}, nil)
		if nbUsers == 0 {
			panic("no dump")
		}
		return err
	})

	code := m.Run()
	kill(mongodb)
	// You can't defer this because os.Exit doesn't care for defer

	os.Exit(code)
}

func kill(resource *dockertest.Resource) {
	if resource == nil {
		panic("mongodb n'a pas démarré")
	}
	if err := resource.Close(); err != nil {
		panic(fmt.Sprintf("Could not purge resource: %s", err))
	}
}

func restoreDump(mongodb *dockertest.Resource) error {
	dump, err := os.Open("test/resources/wekan-users.dump")
	if err != nil {
		panic("le dump de test n'est pas accessible")
	}

	var b bytes.Buffer
	output := bufio.NewWriter(&b)

	options := dockertest.ExecOptions{
		StdIn:  dump,
		StdOut: output,
		StdErr: output,
	}

	_, err = mongodb.Exec([]string{"/bin/bash", "-c", "mongorestore mongodb://root:password@127.0.0.1/ --noIndexRestore --archive"}, options)

	return err
}
