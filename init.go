package libwekan

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Wekan struct {
	url           string
	databaseName  string
	client        *mongo.Client
	db            *mongo.Database
	adminUsername Username
	adminUser     *User
}

// Connect retourne un objet de type `Wekan`
func Connect(ctx context.Context, uri string, databaseName string, adminUsername Username) (Wekan, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return Wekan{}, err
	}
	w := Wekan{
		url:           uri,
		databaseName:  databaseName,
		client:        client,
		db:            client.Database(databaseName),
		adminUsername: adminUsername,
	}

	if err != nil {
		return Wekan{}, err
	}

	return w, nil
}

func (wekan *Wekan) AdminUser(ctx context.Context) (*User, error) {
	if wekan.adminUser == nil {
		admin, err := wekan.GetUserFromUsername(ctx, wekan.adminUsername)
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%s is not in database\nmongo => %s", wekan.adminUsername, err.Error())
		}
		if !admin.IsAdmin {
			return nil, errors.New("%s is not admin")
		}
		wekan.adminUser = &admin
	}
	return wekan.adminUser, nil
}

func (wekan *Wekan) Ping() error {
	return wekan.client.Ping(context.Background(), nil)
}
