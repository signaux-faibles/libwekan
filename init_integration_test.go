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
	"go.mongodb.org/mongo-driver/mongo"
)

var dbClient *mongo.Client
var wekan Wekan
var cwd, _ = os.Getwd()

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
			Mounts: []string{cwd + "/test/resources/:/dump/"},
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
		mongoUrl := fmt.Sprintf("mongodb://root:password@localhost:%s", mongodb.GetPort("27017/tcp"))
		wekan, err = Connect(context.Background(), mongoUrl, "wekan", "signaux.faibles")
		if err != nil {
			fmt.Println(err)
			return err
		}

		if err := wekan.client.Ping(context.TODO(), nil); err != nil {
			return err
		}

		err = restoreDump(mongodb)
		if err != nil {
			return err
		}

		users, err := wekan.GetUsers(context.TODO())
		if err != nil {
			fmt.Println(err)
			return err
		}

		if len(users) == 0 {
			panic("pas d'utilisateurs dans la base de test ?!")
		}

		_, err = wekan.AdminUser(context.TODO())
		if err != nil {
			panic(fmt.Sprintf("pas d'admin dans la base de test: %s", err.Error()))
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
	var b bytes.Buffer
	output := bufio.NewWriter(&b)

	options := dockertest.ExecOptions{
		StdOut: output,
		StdErr: output,
	}

	_, err := mongodb.Exec([]string{"/bin/bash", "-c", "mongorestore  --uri mongodb://root:password@localhost/ /dump"}, options)
	// _, err = mongodb.Exec([]string{"/bin/bash", "-c", "mongo mongodb://root:password@localhost/wekan --authenticationDatabase admin --eval 'printjson(db.users.find({}).toArray())'"}, options)
	output.Flush()
	// fmt.Println(b.String())

	return err
}

// func Test_getUser(t *testing.T) {
// 	asserte := assert.New()
// 	users :=
// }
