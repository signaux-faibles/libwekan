//go:build integration

// nolint:errcheck
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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()
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
		wekan, err = Init(ctx, mongoUrl, "wekan", "signaux.faibles", "^tableau-crp.*")
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

	//adminUser, err := wekan.AdminUser(ctx)
	ass.Nil(wekan.AssertPrivileged(ctx))
	ass.Equal(admin.ID, wekan.adminUserID)
}

func TestGetUser_when_user_not_exist(t *testing.T) {
	ass := assert.New(t)
	user, err := wekan.GetUserFromUsername(context.TODO(), "unexistant user")
	ass.NotNil(err)
	ass.ErrorIs(err, err.(UserNotFoundError))
	ass.Empty(user.ID)

}

func TestWekan_Init_WithBadParams(t *testing.T) {
	// WHEN
	badWekan, err := Init(context.Background(), "badUri", "badDb", "badAdmin", "badDomain")

	// THEN
	assert.Empty(t, badWekan)
	assert.IsType(t, InvalidMongoConfigurationError{}, err)

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

	dockerOptions := dockertest.ExecOptions{
		StdOut: output,
		StdErr: output,
	}

	if _, err := mongodb.Exec([]string{"/bin/bash", "-c", "mongorestore  --uri mongodb://root:password@localhost/ /dump"}, dockerOptions); err != nil {
		return nil
	}
	return output.Flush()
}

func newTestBadWekan(dbname string) Wekan {
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:80").SetConnectTimeout(time.Millisecond).SetTimeout(time.Millisecond)
	client, _ := mongo.Connect(context.Background(), clientOptions)
	return Wekan{
		client: client,
		db:     client.Database(dbname),
	}
}

func TestWekan_CheckDocuments_WithBadDocuments(t *testing.T) {
	// GIVEN
	badUserID := UserID("badUser")

	// WHEN
	err := wekan.CheckDocuments(ctx, badUserID)

	// THEN
	assert.IsType(t, UserNotFoundError{}, err)
}

func TestWekan_AssertPrivileged(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	badAdminWekan := wekan
	False := false
	badAdminWekan.privileged = &False

	// On vérifie les implémentations dans les différentes fonctions qui utilisent cette fonction
	errs := []error{
		badAdminWekan.AssertPrivileged(ctx),
		badAdminWekan.AddMemberToBoard(ctx, "", BoardMember{}),
		badAdminWekan.AddMemberToCard(ctx, "", ""),
		badAdminWekan.AddLabelToCard(ctx, "", ""),
		badAdminWekan.DisableBoardMember(ctx, "", ""),
		badAdminWekan.DisableUser(ctx, User{}),
		badAdminWekan.DisableUsers(ctx, Users{User{}}),
		badAdminWekan.DisableBoardMember(ctx, "", ""),
		badAdminWekan.EnableBoardMember(ctx, "", ""),
		badAdminWekan.EnableUser(ctx, User{}),
		badAdminWekan.EnableUsers(ctx, Users{User{}}),
		badAdminWekan.InsertBoard(ctx, Board{}),
		badAdminWekan.InsertBoardLabel(ctx, Board{}, BoardLabel{}),
		badAdminWekan.InsertAction(ctx, Action{}),
		badAdminWekan.InsertCard(ctx, Card{}),
		badAdminWekan.InsertUser(ctx, User{}),
		badAdminWekan.InsertUsers(ctx, Users{User{}}),
		badAdminWekan.InsertRule(ctx, Rule{}),
		badAdminWekan.InsertList(ctx, List{}),
		badAdminWekan.InsertTrigger(ctx, Trigger{}),
		badAdminWekan.InsertTemplates(ctx, UserTemplates{}),
		badAdminWekan.RemoveMemberFromCard(ctx, "", ""),
		badAdminWekan.RemoveRuleWithID(ctx, ""),
	}
	_, err := badAdminWekan.EnsureMemberInCard(ctx, "", "")
	errs = append(errs, err)
	_, err = badAdminWekan.EnsureMemberOutOfCard(ctx, "", "")
	errs = append(errs, err)
	_, err = badAdminWekan.EnsureRuleAddTaskforceMemberExists(ctx, User{}, Board{}, BoardLabel{})
	errs = append(errs, err)
	_, err = badAdminWekan.EnsureRuleRemoveTaskforceMemberExists(ctx, User{}, Board{}, BoardLabel{})
	errs = append(errs, err)
	_, err = badAdminWekan.EnsureUserIsBoardAdmin(ctx, "", "")
	errs = append(errs, err)
	_, err = badAdminWekan.EnsureUserIsActiveBoardMember(ctx, "", "")
	errs = append(errs, err)
	_, err = badAdminWekan.EnsureUserIsInactiveBoardMember(ctx, "", "")
	errs = append(errs, err)

	for i, err := range errs {
		ass.IsType(NotPrivilegedError{}, err, "echec pour la fonction %d", i)
	}
}
