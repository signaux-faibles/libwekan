package libwekan

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Wekan struct {
	url          string
	databaseName string
	db           *mongo.Database
	admin        User
}

// Connect returns a Wekan context
func Connect(ctx context.Context, server string, databaseName string, username string) (Wekan, error) {
	uri := fmt.Sprintf("mongodb://%s/", server)
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return Wekan{}, err
	}
	w := Wekan{
		url:          uri,
		databaseName: databaseName,
		db:           client.Database(databaseName),
	}
	w.admin, err = w.GetUser(ctx, username)
	if err == mongo.ErrNoDocuments {
		return Wekan{}, fmt.Errorf("%s is not in database\nmongo => %s", username, err.Error())
	}
	if err != nil {
		return Wekan{}, err
	}
	if !w.admin.IsAdmin {
		return Wekan{}, errors.New("%s is not admin")
	}
	return w, nil
}
