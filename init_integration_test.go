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
	"github.com/stretchr/testify/assert"
)

// var dbClient *mongo.Client
var wekan Wekan
var cwd, _ = os.Getwd()

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		panic("ne peut démarrer mongodb")
	}
	// pulls an image, creates a container based on it and runs it
	mongodbContainerName := "wekandb-ti-" + strconv.Itoa(time.Now().Nanosecond())

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
	if err := pool.Retry(func() error {
		fmt.Println("Mongo n'est pas encore prêt")
		var err error
		mongoUrl := fmt.Sprintf("mongodb://root:password@localhost:%s", mongodb.GetPort("27017/tcp"))
		wekan, err = Connect(context.Background(), mongoUrl, "wekan", "signaux.faibles")
		if err != nil {
			return err
		}
		return wekan.client.Ping(context.TODO(), nil)
	}); err != nil {
		fmt.Printf("N'arrive pas à démarrer/restaurer Mongo: %s", err)
	}
	fmt.Println("Mongo est prêt, on lance le restore dump")
	err = restoreDump(mongodb)
	if err != nil {
		panic("Foirage du restore dump : " + err.Error())
	}

	fmt.Println("On peut lancer les tests")
	code := m.Run()
	kill(mongodb)
	// You can't defer this because os.Exit doesn't care for defer

	os.Exit(code)
}

func Test_OnlyOneAdminInDB(t *testing.T) {
	ass := assert.New(t)
	users, err := wekan.GetUsers(context.TODO())
	var admin User
	var admins int
	for _, u := range users {
		if u.IsAdmin {
			admin = u
			admins += 1
		}
	}
	ass.Nil(err)
	ass.Equal(1, admins)

	adminUser, err := wekan.AdminUser(context.Background())
	ass.Nil(err)
	ass.Equal(admin.ID, (adminUser).ID)
}

func TestGetUser_when_user_not_exist(t *testing.T) {
	ass := assert.New(t)
	user, err := wekan.GetUserFromUsername(context.TODO(), Username("unexistant user"))
	ass.NotNil(err)
	ass.ErrorIs(err, err.(UnknownUserError))
	ass.Empty(user.ID)

}

func kill(mongodb *dockertest.Resource) {
	if mongodb == nil {
		panic("mongodb n'a pas démarré")
	}
	if err := mongodb.Close(); err != nil {
		panic(fmt.Sprintf("Could not purge mongodb: %s", err))
	}
}

func restoreDump(mongodb *dockertest.Resource) error {
	var b bytes.Buffer
	output := bufio.NewWriter(&b)

	options := dockertest.ExecOptions{
		StdOut: output,
		StdErr: output,
	}

	if _, err := mongodb.Exec([]string{"/bin/bash", "-c", "mongorestore  --uri mongodb://root:password@localhost/ /dump"}, options); err != nil {
		return nil
	}
	return output.Flush()
}
