package libwekan

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Wekan struct {
	url           string
	databaseName  string
	client        *mongo.Client
	db            *mongo.Database
	adminUsername Username
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

func (wekan *Wekan) AdminUser(ctx context.Context) (User, error) {
	admin, err := wekan.GetUserFromUsername(ctx, wekan.adminUsername)
	if _, ok := err.(UnknownUserError); ok {
		return User{}, err
	}
	if !admin.IsAdmin {
		return User{}, UserIsNotAdminError{admin.ID}
	}
	return admin, nil
}

func (wekan *Wekan) Ping() error {
	return wekan.client.Ping(context.Background(), nil)
}
